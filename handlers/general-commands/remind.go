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

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// ReminderTimers is the list of all reminders running
var ReminderTimers []structs.ReminderTimer

// Remind reminds the person after an x amount of specified time
func Remind(s *discordgo.Session, m *discordgo.MessageCreate) {
	remindRegex, _ := regexp.Compile(`remind\s+(.+)`)
	timeRegex, _ := regexp.Compile(`\s(\d+) (month|week|day|hour|minute|second)s?`)
	reminderTime := time.Duration(0)
	text := ""
	timeResultString := ""
	// Parse info
	if remindRegex.MatchString(m.Content) {
		text = remindRegex.FindStringSubmatch(m.Content)[1]
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
				case "hour":
					hours += timeVal
				case "minute":
					minutes += timeVal
				case "second":
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
		}
	}
	if reminderTime == 0 { // Default to 5 minutes
		reminderTime = time.Second * time.Duration(300)
	}
	// Check if time duration is dumb as hell
	if reminderTime.Hours() > 8760 {
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
	reminder := structs.NewReminder(timeResult, *m.Author, text)
	reminders := []structs.Reminder{}
	_, err := os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(err)
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
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(err)

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
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		tools.ErrRead(err)
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
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(err)
}

// Reminders lists the person's reminders
func Reminders(s *discordgo.Session, m *discordgo.MessageCreate) {
	userTimers := []structs.Reminder{}
	all := false
	if strings.Contains(m.Content, "all") {
		if m.Author.ID != "92502458588205056" {
			s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
			return
		}
		all = true
		for _, reminder := range ReminderTimers {
			if reminder.Reminder.Active {
				userTimers = append(userTimers, reminder.Reminder)
			}
		}
	} else {
		for _, reminder := range ReminderTimers {
			if reminder.Reminder.User.ID == m.Author.ID && reminder.Reminder.Active {
				userTimers = append(userTimers, reminder.Reminder)
			}
		}
	}

	if len(userTimers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no pending reminders!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username + "#" + m.Author.Discriminator,
			IconURL: m.Author.AvatarURL(""),
		},
		Description: "Please use `rremove <ID>` or `remindremove <ID>` to remove a reminder",
	}
	if all {
		for _, reminder := range userTimers {
			info := reminder.Info
			if info == "" {
				info = "N/A"
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   strconv.FormatInt(reminder.ID, 10),
				Value:  "Reminder: " + info + "\nRemind time: " + reminder.Target.Format(time.RFC822) + "\nUser: " + reminder.User.Username + "#" + reminder.User.Discriminator,
				Inline: true,
			})
		}
	} else {
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
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// RemoveReminder removes a reminder (kind of)
func RemoveReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	remindRegex, _ := regexp.Compile(`r(emind)?remove\s+(\d+|all)`)
	if !remindRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give a reminder's snowflake ID to remove! You can see all of your reminds with `reminders`. If you want to remove all reminders, please state `remindremove all`")
		return
	}

	// Get reminders
	reminders := []structs.Reminder{}
	_, err := os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		tools.ErrRead(err)
	}

	// Mark Active as false for the reminder in both slices
	reminderID := remindRegex.FindStringSubmatch(m.Content)[2]
	if reminderID == "all" {
		for i, reminder := range ReminderTimers {
			if reminder.Reminder.User.ID == m.Author.ID {
				ReminderTimers[i].Reminder.Active = false
			}
		}
		for i, reminder := range reminders {
			if reminder.User.ID == m.Author.ID {
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
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/reminders.json", jsonCache, 0644)
	tools.ErrRead(err)
}

// ReminderMessage will send the user their reminder
func ReminderMessage(s *discordgo.Session, reminderTimer structs.ReminderTimer) {
	linkRegex, _ := regexp.Compile(`https?:\/\/\S+`)
	dm, _ := s.UserChannelCreate(reminderTimer.Reminder.User.ID)
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
				&discordgo.File{
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
