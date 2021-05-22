package gencommands

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	structs "maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// TaskCache is the list of all tasks running
var TaskCache []structs.TaskCache

// Task tasks the person after an x amount of specified time
func Task(s *discordgo.Session, m *discordgo.MessageCreate) {
	taskRegex, _ := regexp.Compile(`(?i)task\s+(.+)`)
	timeRegex, _ := regexp.Compile(`(?i)\s(\d+) (month|week|day|h(ou)?r|min(ute)?|sec(ond)?)s?`)
	dateRegex, _ := regexp.Compile(`(?i)starting\s+(at\s+)?(.+)`)
	taskDuration := time.Duration(0)
	startTime := time.Now().UTC()

	if !taskRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please provide information for the task you wish to create!")
		return
	}

	// Parse info
	text := taskRegex.FindStringSubmatch(m.Content)[1]
	if !timeRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Could not parse any form of a time interval to use for the task!")
		return
	}

	times := timeRegex.FindAllStringSubmatch(m.Content, -1)
	months := 0
	weeks := 0
	days := 0
	hours := 0
	minutes := 0
	seconds := 0
	for _, timeString := range times {
		timeVal, err := strconv.Atoi(timeString[1])
		if err != nil {
			break
		}
		timeUnit := timeString[2]
		switch timeUnit {
		case "month":
			months += timeVal
		case "week":
			weeks += timeVal
		case "day":
			days += timeVal
		case "hour", "hr":
			hours += timeVal
		case "minute", "min":
			minutes += timeVal
		case "second", "sec":
			seconds += timeVal
		}
		text = strings.Replace(text, strings.TrimSpace(timeString[0]), "", 1)
		text = strings.TrimSpace(text)
		text = strings.TrimSuffix(text, "and")
		text = strings.TrimSpace(text)
		text = strings.TrimSuffix(text, ",")
	}
	taskDuration += time.Second * time.Duration(months) * 2629744
	taskDuration += time.Second * time.Duration(weeks) * 604800
	taskDuration += time.Second * time.Duration(days) * 86400
	taskDuration += time.Second * time.Duration(hours) * 3600
	taskDuration += time.Second * time.Duration(minutes) * 60
	taskDuration += time.Second * time.Duration(seconds)

	if dateRegex.MatchString(m.Content) {
		// Parse date
		date := dateRegex.FindStringSubmatch(m.Content)[2]
		t, err := tools.TimeParse(date)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid datetime format! Error: "+err.Error())
			return
		}

		if t.Year() == 0 {
			t = t.AddDate(time.Now().Year(), 0, 0)
		} else if t.Year() == 1 {
			t = t.AddDate(time.Now().Year()+1, 0, 0)
		}

		for t.Before(time.Now().UTC()) {
			t = t.Add(taskDuration)
		}

		startTime = t
		text = dateRegex.ReplaceAllString(text, "")
	} else {
		startTime = startTime.Add(taskDuration)
	}

	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, "every")
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, "in")
	text = strings.TrimSpace(text)

	if taskDuration == 0 { // Default to 1 day
		taskDuration = time.Hour * 24
	} else if taskDuration.Hours() > 17520 { // Check if time duration is dumb as hell
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}

	// Obtain date
	startTimeString := startTime.Format(time.UnixDate)
	text = strings.TrimSpace(strings.ReplaceAll(text, "`", ""))

	// People can add huge time durations where the time may go backward
	if startTime.Before(time.Now()) {
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}

	// Create task and add to list of tasks
	task := structs.NewTask(m.Author.ID, text, taskDuration, startTime)
	tasks := []structs.Task{}
	_, err := os.Stat("./data/tasks.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/tasks.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &tasks)
	} else {
		s.ChannelMessageSend(m.ChannelID, "An error occurred obtaining task data! Please try later.")
		return
	}
	tasks = append(tasks, task)

	TaskCache = append(TaskCache, structs.TaskCache{
		Task:   task,
		Active: true,
	})

	// Save tasks
	jsonCache, err := json.Marshal(tasks)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/tasks.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	if text != "" {
		s.ChannelMessageSend(m.ChannelID, "Ok I'll task you for `"+task.Info+"` starting on "+startTimeString+"\nPlease make sure your DMs are open or else you will not receive the task!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Ok I'll task you starting on "+startTimeString+"\nPlease make sure your DMs are open or else you will not receive the task!")
	}

	// Run task
	go RunTask(s, task)
}

// RunTask runs the task
func RunTask(s *discordgo.Session, task structs.Task) {
	for task.LastRun.Before(time.Now().UTC()) {
		task.LastRun = task.LastRun.Add(time.Second * time.Duration(task.Seconds))
	}

	// Reach the point of the next iteration of sending the message
	<-time.After(task.LastRun.Sub(time.Now().UTC()))
	for i, activeTask := range TaskCache {
		if task.ID == activeTask.Task.ID {
			if !activeTask.Active {
				TaskCache[i] = TaskCache[len(TaskCache)-1]
				TaskCache = TaskCache[:len(TaskCache)-1]
				return
			}

			go TaskMessage(s, task)
			break
		}
	}

	// Start ticker
	taskTicker := time.NewTicker(time.Second * time.Duration(task.Seconds))
	for range taskTicker.C {
		for i, activeTask := range TaskCache {
			if task.ID == activeTask.Task.ID {
				if !activeTask.Active {
					TaskCache[i] = TaskCache[len(TaskCache)-1]
					TaskCache = TaskCache[:len(TaskCache)-1]
					taskTicker.Stop()
					return
				}

				go TaskMessage(s, task)
				break
			}
		}
	}
}

// Tasks lists the person's tasks
func Tasks(s *discordgo.Session, m *discordgo.MessageCreate) {
	userTimers := []structs.Task{}
	for _, task := range TaskCache {
		if task.Task.User == m.Author.ID && task.Active {
			userTimers = append(userTimers, task.Task)
		}
	}

	if len(userTimers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no pending tasks!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.String(),
			IconURL: m.Author.AvatarURL("2048"),
		},
		Description: "Please use `tremove <ID>` or `taskremove <ID>` to remove a task",
	}
	for _, task := range userTimers {
		info := task.Info
		if info == "" {
			info = "N/A"
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   strconv.FormatInt(task.ID, 10),
			Value:  "Task: " + info + "\nTask interval (seconds): " + strconv.FormatFloat(task.Seconds, 'f', 0, 64) + "\nLast Run: " + task.LastRun.Format(time.RFC822),
			Inline: true,
		})
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// RemoveTask removes a task (kind of)
func RemoveTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	taskRegex, _ := regexp.Compile(`(?i)t(ask)?remove\s+(\d+|all)`)
	if !taskRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give a task's snowflake ID to remove! You can see all of your tasks with `tasks`. If you want to remove all tasks, please state `taskremove all`")
		return
	}

	// Get tasks
	tasks := []structs.Task{}
	_, err := os.Stat("./data/tasks.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/tasks.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &tasks)
	} else {
		tools.ErrRead(s, err)
	}

	// Mark Active as false for the task in both slices and remove from json
	taskID := taskRegex.FindStringSubmatch(m.Content)[2]
	text := "tasks"
	if taskID == "all" {
		for i, task := range TaskCache {
			if task.Task.User == m.Author.ID {
				TaskCache[i].Active = false
			}
		}
		for i, task := range tasks {
			if task.User == m.Author.ID {
				tasks[i] = tasks[len(tasks)-1]
				tasks = tasks[:len(tasks)-1]
			}
		}
	} else {
		taskIDint, err := strconv.ParseInt(taskID, 10, 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error parsing ID.")
			return
		}
		for i, task := range TaskCache {
			if task.Task.ID == taskIDint {
				TaskCache[i].Active = false
				break
			}
		}
		for i, task := range tasks {
			if task.ID == taskIDint {
				tasks[i] = tasks[len(tasks)-1]
				tasks = tasks[:len(tasks)-1]
				break
			}
		}
		text = "task"
	}

	// Save tasks
	jsonCache, err := json.Marshal(tasks)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/tasks.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, "Removed "+text+"!")
}

// TaskMessage will send the user their task
func TaskMessage(s *discordgo.Session, task structs.Task) {
	// Update lastRun, get tasks
	tasks := []structs.Task{}
	_, err := os.Stat("./data/tasks.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/tasks.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &tasks)
	} else {
		tools.ErrRead(s, err)
	}

	for i, activeTask := range tasks {
		if task.ID == activeTask.ID {
			tasks[i].LastRun = time.Now().UTC()
			break
		}
	}

	// Save tasks
	jsonCache, err := json.Marshal(tasks)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/tasks.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S+`)
	dm, _ := s.UserChannelCreate(task.User)
	if task.Info == "" {
		s.ChannelMessageSend(dm.ID, "Task interval hit!")
	} else if linkRegex.MatchString(task.Info) {
		response, err := http.Get(linkRegex.FindStringSubmatch(task.Info)[0])
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Task: `"+task.Info+"`")
			return
		}
		img, _, err := image.Decode(response.Body)
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Task: `"+task.Info+"`")
			return
		}
		imgBytes := new(bytes.Buffer)
		err = png.Encode(imgBytes, img)
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Task: `"+task.Info+"`")
		}
		_, err = s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: "Task: ",
			Files: []*discordgo.File{
				{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Task: `"+task.Info+"`!")
		}
		response.Body.Close()
		return
	} else {
		s.ChannelMessageSend(dm.ID, "Task: `"+task.Info+"`!")
	}
}
