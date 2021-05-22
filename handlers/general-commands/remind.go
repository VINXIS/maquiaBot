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

// ReminderTimers is the list of all reminders running
var ReminderTimers []structs.ReminderTimer

// Remind reminds the person after an x amount of specified time
func Remind(s *discordgo.Session, m *discordgo.MessageCreate) {
	remindRegex, _ := regexp.Compile(`(?i)remind(er)?\s+(.+)`)
	timeRegex, _ := regexp.Compile(`(?i)\s(\d+) (month|week|day|h(ou)?r|min(ute)?|sec(ond)?)s?`)
	dateRegex, _ := regexp.Compile(`(?i)at\s+(.+)`)
	reminderTime := time.Duration(0)
	text := ""
	timeResultString := ""
	// Parse info
	if remindRegex.MatchString(m.Content) {
		text = remindRegex.FindStringSubmatch(m.Content)[2]
		if timeRegex.MatchString(m.Content) {
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
				text = strings.TrimSuffix(text, ",")
			}
			text = strings.TrimSpace(text)
			text = strings.TrimSuffix(text, "in")
			text = strings.TrimSpace(text)
			reminderTime += time.Second * time.Duration(months) * 2629744
			reminderTime += time.Second * time.Duration(weeks) * 604800
			reminderTime += time.Second * time.Duration(days) * 86400
			reminderTime += time.Second * time.Duration(hours) * 3600
			reminderTime += time.Second * time.Duration(minutes) * 60
			reminderTime += time.Second * time.Duration(seconds)
		} else if dateRegex.MatchString(m.Content) {
			// Parse date
			date := dateRegex.FindStringSubmatch(m.Content)[1]
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

			reminderTime = t.Sub(time.Now())
			text = dateRegex.ReplaceAllString(text, "")
		}
	}
	if reminderTime == 0 { // Default to 5 minutes
		reminderTime = time.Second * time.Duration(300)
	}
	// Check if time duration is dumb as hell
	if reminderTime.Hours() > 17520 {
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}

	// Obtain date
	timeResult := time.Now().UTC().Add(reminderTime)
	timeResultString = timeResult.Format(time.UnixDate)
	text = strings.ReplaceAll(text, "`", "")

	// People can add huge time durations where the time may go backward
	if timeResult.Before(time.Now()) {
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}

	// Create reminder and add to list of reminders
	reminder := structs.NewReminder(timeResult, m.Author.ID, text)
	reminders := []structs.Reminder{}
	_, err := os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		s.ChannelMessageSend(m.ChannelID, "An error occurred obtaining reminder data! Please try later.")
		return
	}
	reminders = append(reminders, reminder)
	reminderTimer := structs.ReminderTimer{
		Reminder: reminder,
		Timer:    *time.NewTimer(timeResult.Sub(time.Now().UTC())),
	}
	ReminderTimers = append(ReminderTimers, reminderTimer)

	// Save reminders
	jsonCache, err := json.Marshal(reminders)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	if text != "" {
		s.ChannelMessageSend(m.ChannelID, "Ok I'll remind you about `"+reminder.Info+"` on "+timeResultString+"\nPlease make sure your DMs are open or else you will not receive the reminder!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Ok I'll remind you on "+timeResultString+"\nPlease make sure your DMs are open or else you will not receive the reminder!")
	}
	// Run reminder
	go RunReminder(s, reminderTimer)
}

// RunReminder runs the reminder
func RunReminder(s *discordgo.Session, reminderTimer structs.ReminderTimer) {
	if time.Now().Before(reminderTimer.Reminder.Target) {
		<-reminderTimer.Timer.C
	}
	for i, rmndr := range ReminderTimers {
		if rmndr.Reminder.ID == reminderTimer.Reminder.ID {
			if rmndr.Reminder.Active {
				go ReminderMessage(s, reminderTimer)
			}
			ReminderTimers[i] = ReminderTimers[len(ReminderTimers)-1]
			ReminderTimers = ReminderTimers[:len(ReminderTimers)-1]
			break
		}
	}

	// Remove reminder
	reminders := []structs.Reminder{}
	_, err := os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		tools.ErrRead(s, err)
	}
	for i, reminder := range reminders {
		if reminder.ID == reminderTimer.Reminder.ID {
			reminders[i] = reminders[len(reminders)-1]
			reminders = reminders[:len(reminders)-1]
			break
		}
	}

	// Save reminders
	jsonCache, err := json.Marshal(reminders)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(s, err)
}

// Reminders lists the person's reminders
func Reminders(s *discordgo.Session, m *discordgo.MessageCreate) {
	userTimers := []structs.Reminder{}
	for _, reminder := range ReminderTimers {
		if reminder.Reminder.User == m.Author.ID && reminder.Reminder.Active {
			userTimers = append(userTimers, reminder.Reminder)
		}
	}

	if len(userTimers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no pending reminders!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.String(),
			IconURL: m.Author.AvatarURL("2048"),
		},
		Description: "Please use `rremove <ID>` or `remindremove <ID>` to remove a reminder",
	}
	for _, reminder := range userTimers {
		info := reminder.Info
		if info == "" {
			info = "N/A"
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   strconv.FormatInt(reminder.ID, 10),
			Value:  "Reminder: " + info + "\nRemind time: " + reminder.Target.Format(time.RFC822),
			Inline: true,
		})
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// RemoveReminder removes a reminder (kind of)
func RemoveReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	remindRegex, _ := regexp.Compile(`(?i)r(emind)?remove\s+(\d+|all)`)
	if !remindRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give a reminder's snowflake ID to remove! You can see all of your reminds with `reminders`. If you want to remove all reminders, please state `remindremove all`")
		return
	}

	// Get reminders
	reminders := []structs.Reminder{}
	_, err := os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(s, err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		tools.ErrRead(s, err)
	}

	// Mark Active as false for the reminder in both slices
	reminderID := remindRegex.FindStringSubmatch(m.Content)[2]
	if reminderID == "all" {
		for i, reminder := range ReminderTimers {
			if reminder.Reminder.User == m.Author.ID {
				ReminderTimers[i].Reminder.Active = false
			}
		}
		for i, reminder := range reminders {
			if reminder.User == m.Author.ID {
				reminders[i].Active = false
			}
		}
		s.ChannelMessageSend(m.ChannelID, "Removed reminders!")
	} else {
		reminderIDint, err := strconv.ParseInt(reminderID, 10, 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error parsing ID.")
			return
		}
		for i, reminder := range ReminderTimers {
			if reminder.Reminder.ID == reminderIDint {
				ReminderTimers[i].Reminder.Active = false
				break
			}
		}
		for i, reminder := range reminders {
			if reminder.ID == reminderIDint {
				reminders[i].Active = false
				break
			}
		}
		s.ChannelMessageSend(m.ChannelID, "Removed reminder!")
	}

	// Save reminders
	jsonCache, err := json.Marshal(reminders)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(s, err)
}

// ReminderMessage will send the user their reminder
func ReminderMessage(s *discordgo.Session, reminderTimer structs.ReminderTimer) {
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S+`)
	dm, _ := s.UserChannelCreate(reminderTimer.Reminder.User)
	if reminderTimer.Reminder.Info == "" {
		s.ChannelMessageSend(dm.ID, "Reminder!")
	} else if linkRegex.MatchString(reminderTimer.Reminder.Info) {
		response, err := http.Get(linkRegex.FindStringSubmatch(reminderTimer.Reminder.Info)[0])
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Reminder about `"+reminderTimer.Reminder.Info+"`!")
			return
		}
		img, _, err := image.Decode(response.Body)
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Reminder about `"+reminderTimer.Reminder.Info+"`!")
			return
		}
		imgBytes := new(bytes.Buffer)
		err = png.Encode(imgBytes, img)
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Reminder about `"+reminderTimer.Reminder.Info+"`!")
		}
		_, err = s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: "Reminder about this",
			Files: []*discordgo.File{
				{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
		if err != nil {
			s.ChannelMessageSend(dm.ID, "Reminder about `"+reminderTimer.Reminder.Info+"`!")
		}
		response.Body.Close()
		return
	} else {
		s.ChannelMessageSend(dm.ID, "Reminder about `"+reminderTimer.Reminder.Info+"`!")
	}
}
