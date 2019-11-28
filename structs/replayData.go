package structs

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	osuapi "../osu-api"
	"github.com/ulikunitz/xz"
)

// ReplayData stores the replay information of a play
type ReplayData struct {
	Time         time.Time
	Data         []byte
	Mode         osuapi.Mode
	Beatmap      osuapi.Beatmap
	Player       osuapi.User
	Score        osuapi.Score
	LifeBar      []HealthData
	PlayData     []PlayData
	UnstableRate float64
}

// HealthData holds the health data of the player in the replay
type HealthData struct {
	TimeStamp int64
	Health    float64
}

// PlayData stores the actual play data of a replay
type PlayData struct {
	TimeStamp int64
	TimeSince int64
	X         float64
	Y         float64
	PressType Press
}

// Press is the type of press that occurred
type Press int

// Press types
const (
	M1 Press = 1 << iota
	M2
	K1
	K2
	Smoke
)

// ParseReplay parses the replay and fills in the ReplayData's values
func (r *ReplayData) ParseReplay(osuAPI *osuapi.Client) {
	r.Mode = r.getMode()
	r.Beatmap = r.getBeatmap(osuAPI)
	r.Player = r.getUser(osuAPI)

	// Skip replay hash cuz its useless Lol
	if r.Data[0] == 0 {
		r.Data = r.Data[1:]
	} else {
		r.Data = r.Data[1:]
		hashLength, offset := uleb(r.Data)
		r.Data = r.Data[hashLength+offset-1:]
	}

	r.Score = r.getScore(osuAPI)
	r.LifeBar = r.getLife()
	r.Time = r.getTime()
	r.PlayData = r.getPlayData()
}

func (r *ReplayData) getMode() osuapi.Mode {
	mode := r.Data[0]
	r.Data = r.Data[5:]
	if mode > 4 {
		return 0
	}
	return osuapi.Mode(mode)
}

func (r *ReplayData) getBeatmap(osuAPI *osuapi.Client) osuapi.Beatmap {
	hash := ""
	if r.Data[0] == 0 {
		r.Data = r.Data[1:]
		return osuapi.Beatmap{}
	}
	r.Data = r.Data[1:]
	hashLength, offset := uleb(r.Data)
	hash = string(r.Data[offset:hashLength])
	r.Data = r.Data[hashLength+offset-1:]

	beatmap, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
		BeatmapHash: hash,
	})
	if err != nil || len(beatmap) == 0 {
		return osuapi.Beatmap{FileMD5: hash}
	}
	return beatmap[0]
}

func (r *ReplayData) getUser(osuAPI *osuapi.Client) osuapi.User {
	username := ""
	if r.Data[0] == 0 {
		r.Data = r.Data[1:]
		return osuapi.User{}
	}
	r.Data = r.Data[1:]
	usernameLength, offset := uleb(r.Data)
	username = string(r.Data[offset:usernameLength])
	r.Data = r.Data[usernameLength+offset-1:]

	user, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		return osuapi.User{Username: username}
	}
	return *user
}

func (r *ReplayData) getScore(osuAPI *osuapi.Client) osuapi.Score {
	score := osuapi.Score{}
	score.Count300 = int(uint32(r.Data[1])<<8 | uint32(r.Data[0]))
	score.Count100 = int(uint32(r.Data[3])<<8 | uint32(r.Data[2]))
	score.Count50 = int(uint32(r.Data[5])<<8 | uint32(r.Data[4]))
	score.CountGeki = int(uint32(r.Data[7])<<8 | uint32(r.Data[6]))
	score.CountKatu = int(uint32(r.Data[9])<<8 | uint32(r.Data[8]))
	score.CountMiss = int(uint32(r.Data[11])<<8 | uint32(r.Data[10]))
	score.Score = int64(uint64(r.Data[15])<<24 | uint64(r.Data[14])<<16 | uint64(r.Data[13])<<8 | uint64(r.Data[12]))
	score.MaxCombo = int(uint32(r.Data[17])<<8 | uint32(r.Data[16]))
	score.FullCombo = r.Data[18] == 1
	score.Mods = osuapi.Mods(uint32(r.Data[22])<<24 | uint32(r.Data[21])<<16 | uint32(r.Data[20])<<8 | uint32(r.Data[19]))
	r.Data = r.Data[23:]
	return score
}

func (r *ReplayData) getLife() []HealthData {
	if r.Data[0] == 0 {
		r.Data = r.Data[1:]
		return []HealthData{}
	}
	r.Data = r.Data[1:]
	lifeLength, offset := uleb(r.Data)
	life := string(r.Data[offset:lifeLength])
	lifeData := strings.Split(life, ",")
	healthData := []HealthData{}
	for _, interval := range lifeData {
		parts := strings.Split(interval, "|")
		if len(parts) < 2 {
			continue
		}
		timeStamp, _ := strconv.ParseInt(parts[0], 10, 64)
		health, _ := strconv.ParseFloat(parts[1], 64)
		healthData = append(healthData, HealthData{
			TimeStamp: timeStamp,
			Health:    health,
		})
	}

	r.Data = r.Data[lifeLength+offset-1:]
	return healthData
}

func (r *ReplayData) getTime() time.Time {
	ticks := int64(
		uint64(r.Data[7])<<56 |
			uint64(r.Data[6])<<48 |
			uint64(r.Data[5])<<40 |
			uint64(r.Data[4])<<32 |
			uint64(r.Data[3])<<24 |
			uint64(r.Data[2])<<16 |
			uint64(r.Data[1])<<8 |
			uint64(r.Data[0]))
	divider := 100
	duration := time.Duration(ticks)
	scoreTime := time.Time{}
	for i := 0; i < divider; i++ {
		scoreTime = scoreTime.Add(duration)
	}
	r.Data = r.Data[8:]
	return scoreTime
}

func (r *ReplayData) getPlayData() []PlayData {
	if r.Mode != 0 {
		r.Data = []byte{}
		return []PlayData{}
	}

	// Get length, and decompress LZMA stream
	l := int(uint32(r.Data[3])<<24 | uint32(r.Data[2])<<16 | uint32(r.Data[1])<<8 | uint32(r.Data[0]))
	r.Data = r.Data[4:]
	playBytes := r.Data[:l]
	buffer := bytes.NewBuffer(playBytes)
	reader, err := xz.NewReader(buffer)
	if err != nil {
		r.Data = []byte{}
		return []PlayData{}
	}
	playBytes, _ = ioutil.ReadAll(reader)
	playDataString := string(playBytes)
	hits := strings.Split(playDataString, ",")

	// Get play data
	playData := []PlayData{}
	timeElapsed := int64(0)
	for _, hit := range hits {
		parts := strings.Split(hit, "|")

		// Obtain data
		timeSince, _ := strconv.ParseInt(parts[0], 10, 64)
		timeElapsed += timeSince
		x, _ := strconv.ParseFloat(parts[1], 64)
		y, _ := strconv.ParseFloat(parts[2], 64)
		press, _ := strconv.Atoi(parts[3])

		// Append
		playData = append(playData, PlayData{
			TimeStamp: timeElapsed,
			TimeSince: timeSince,
			X:         x,
			Y:         y,
			PressType: Press(press),
		})
	}

	r.Data = []byte{}
	return playData
}

func uleb(byteArray []byte) (int, int) {
	result := 0
	shift := 0
	i := 0
	for {
		b := byteArray[i]
		i++
		result = result | ((int(b) & 0b01111111) << shift)
		if (int(b) & 0b10000000) == 0 {
			break
		}
		shift += 7
	}
	return result + 1, i
}
