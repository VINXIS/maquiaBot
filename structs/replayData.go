package structs

import (
	"bytes"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	osuapi "../osu-api"
	parser "github.com/natsukagami/go-osu-parser"
	"github.com/ulikunitz/xz/lzma"
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
	Early        float64
	Late         float64
	Seed         float64
	HitErrors    []float64
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

// Object extends parser.HitObject
type Object struct {
	parser.HitObject
	stackHeight int
}

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
	r.PlayData = r.GetPlayData(false)
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

	// Score Rank
	percent300 := float64(score.Count300) / float64(score.CountMiss+score.Count50+score.Count100+score.Count300)
	percent50 := float64(score.Count50) / float64(score.CountMiss+score.Count50+score.Count100+score.Count300)
	switch {
	case percent300 == 1:
		score.Rank = "SS"
	case percent300 > 0.9 && percent50 < 0.01 && score.CountMiss == 0:
		score.Rank = "S"
	case percent300 > 0.8 && score.CountMiss == 0, percent300 > 0.9:
		score.Rank = "A"
	case percent300 > 0.7 && score.CountMiss == 0, percent300 > 0.8:
		score.Rank = "B"
	case percent300 > 0.6:
		score.Rank = "C"
	default:
		score.Rank = "D"
	}
	if (score.Mods&osuapi.ModFlashlight != 0 || score.Mods&osuapi.ModHidden != 0) && (score.Rank == "S" || score.Rank == "SS") {
		score.Rank += "H"
	}

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

// GetPlayData obtains the play data from the replay
func (r *ReplayData) GetPlayData(isAPI bool) []PlayData {
	if r.Mode != 0 {
		r.Data = []byte{}
		return []PlayData{}
	}

	// Get length, and decompress LZMA stream
	var playBytes []byte
	if !isAPI {
		l := int(uint32(r.Data[3])<<24 | uint32(r.Data[2])<<16 | uint32(r.Data[1])<<8 | uint32(r.Data[0]))
		r.Data = r.Data[4:]
		playBytes = r.Data[:l]
	} else {
		playBytes = r.Data
	}
	buffer := bytes.NewBuffer(playBytes)
	reader, err := lzma.NewReader(buffer)
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

		if parts[0] == "-12345" && parts[1] == "0" && parts[2] == "0" {
			r.Seed, _ = strconv.ParseFloat(parts[3], 64)
			break
		}

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

// GetUnstableRate gets the unstable rate of the play MUST BE CALLED AFTER osutools.BeatmapParse
func (r *ReplayData) GetUnstableRate() float64 {
	if len(r.PlayData) == 0 {
		return 0
	}

	// Get info for helping to determine hit error // TODO: stacks
	radius := 64.0 * (1.0 - 0.7*(r.Beatmap.CircleSize-5.0)/5.0) / 2.0
	window50 := 199.5 - r.Beatmap.OverallDifficulty*10.0

	// Get map
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)
	beatmap, err := parser.ParseFile("./data/osuFiles/" + strconv.Itoa(r.Beatmap.BeatmapID) + " " + replacer.ReplaceAllString(r.Beatmap.Artist, "") + " - " + replacer.ReplaceAllString(r.Beatmap.Title, "") + ".osu")
	if err != nil {
		return 0
	}

	if r.Score.Mods&osuapi.ModHardRock != 0 {
		for i := 0; i < len(beatmap.HitObjects); i++ {
			obj := &beatmap.HitObjects[i]

			if obj.ObjectName == "spinner" {
				continue
			}

			obj.Position.Y = 384 - obj.Position.Y
			if obj.ObjectName == "slider" {
				obj.EndPosition.Y = 384 - obj.Position.Y
				for j := 0; j < len(obj.Points); j++ {
					obj.Points[j].Y = 384 - obj.Points[j].Y
				}
			}
		}
	}

	version, _ := strconv.Atoi(beatmap.FileFormat[1:])
	if version >= 6 {
		applyStacking(&beatmap)
	} else {
		applyStackingOld(&beatmap)
	}

	usedPlays := []PlayData{}
	for _, obj := range beatmap.HitObjects {
		if obj.ObjectName == "spinner" {
			continue
		}

		for j, play := range r.PlayData {
			// Check if in 50 window
			if play.TimeStamp < int64(obj.StartTime)-int64(window50) {
				continue
			} else if play.TimeStamp > int64(obj.StartTime)+int64(window50) {
				break
			}

			// Check if used already
			used := false
			for _, usedPlay := range usedPlays {
				if usedPlay == play {
					used = true
					break
				}
			}
			if !used {
				// Check if play is a press and in circle
				inCircle := math.Pow(play.X-obj.Position.X, 2)+math.Pow(play.Y-obj.Position.Y, 2) < math.Pow(radius, 2)
				m1 := play.PressType&1 != 0 && r.PlayData[j-1].PressType&1 == 0
				m2 := play.PressType&2 != 0 && r.PlayData[j-1].PressType&2 == 0
				k1 := play.PressType&4 != 0 && r.PlayData[j-1].PressType&4 == 0
				k2 := play.PressType&8 != 0 && r.PlayData[j-1].PressType&8 == 0
				press := m1 || m2 || k1 || k2
				if inCircle && press {
					r.HitErrors = append(r.HitErrors, float64(play.TimeStamp-int64(obj.StartTime)))
					usedPlays = append(usedPlays, play)
					break
				}
			}
		}
	}

	// Get Std Deviation
	avgHitError := 0.0
	earlyCount := 0
	earlyTotal := float64(0)
	lateCount := 0
	lateTotal := float64(0)
	for _, hitError := range r.HitErrors {
		avgHitError += hitError
		if hitError >= 0 {
			lateTotal += hitError
			lateCount++
		} else {
			earlyTotal += hitError
			earlyCount++
		}
	}
	if earlyCount > 0 {
		r.Early = earlyTotal / float64(earlyCount)
	}
	if lateCount > 0 {
		r.Late = lateTotal / float64(lateCount)
	}
	avgHitError /= float64(len(r.HitErrors) - 1)
	stdDevHitError := 0.0
	for _, hitError := range r.HitErrors {
		stdDevHitError += math.Pow(hitError-avgHitError, 2)
	}
	stdDevHitError /= float64(len(r.HitErrors))
	stdDevHitError = math.Sqrt(stdDevHitError)
	unstableRate := stdDevHitError * 10
	if r.Score.Mods&osuapi.ModDoubleTime != 0 {
		unstableRate /= 1.5
	} else if r.Score.Mods&osuapi.ModHalfTime != 0 {
		unstableRate /= 0.75
	}
	return unstableRate
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

func applyStacking(beatmap *parser.Beatmap) {
	scale := (1.0 - 0.7*(beatmap.CircleSize-5.0)/5.0) / 2.0
	ARMS := diffRange(beatmap.ApproachRate)
	stackThresh := int(beatmap.StackLeniency * ARMS)

	// Add stack height feature
	rawObjs := beatmap.HitObjects
	objs := []Object{}
	for _, object := range rawObjs {
		objs = append(objs, Object{object, 0})
	}

	// Obtain stack heights
	for i := len(objs) - 1; i > 0; i-- {
		n := i - 1

		obji := &objs[i]
		if obji.stackHeight != 0 || obji.ObjectName == "spinner" {
			continue
		}

		if obji.ObjectName == "circle" {
			for n-1 >= 0 {
				objn := &objs[n]
				n--
				if objn.ObjectName == "spinner" {
					continue
				}

				if obji.StartTime-objn.EndTime > stackThresh {
					break
				}

				if objn.ObjectName == "slider" && distance(objn.EndPosition, obji.Position) < 3 {
					offset := obji.stackHeight - objn.stackHeight + 1

					for j := n + 1; j <= i; j++ {
						objj := &objs[j]
						if distance(objn.EndPosition, objj.Position) < 3 {
							objj.stackHeight -= offset
						}
					}

					break
				}

				if distance(objn.Position, obji.Position) < 3 {
					objn.stackHeight = obji.stackHeight + 1
					obji = objn
				}
			}
		} else if obji.ObjectName == "slider" {
			for n-1 >= 0 {
				objn := &objs[n]
				n--
				if objn.ObjectName == "spinner" {
					continue
				}

				if obji.StartTime-objn.StartTime > stackThresh {
					break
				}

				if distance(objn.Position, obji.Position) < 3 {
					objn.stackHeight = obji.stackHeight + 1
					obji = objn
				}
			}
		}

	}

	for i := 0; i < len(beatmap.HitObjects)-1; i++ {
		offset := float64(objs[i].stackHeight) * scale * -6.4
		beatmap.HitObjects[i].Position.X += offset
		beatmap.HitObjects[i].Position.Y += offset
	}
}

func applyStackingOld(beatmap *parser.Beatmap) {
	scale := (1.0 - 0.7*(beatmap.CircleSize-5.0)/5.0) / 2.0
	ARMS := diffRange(beatmap.ApproachRate)
	stackThresh := int(beatmap.StackLeniency * ARMS)

	// Add stack height feature
	rawObjs := beatmap.HitObjects
	objs := []Object{}
	for _, object := range rawObjs {
		objs = append(objs, Object{object, 0})
	}

	for i := 0; i < len(objs); i++ {
		currObj := &objs[i]
		if currObj.stackHeight != 0 && currObj.ObjectName != "slider" {
			continue
		}

		sliderStack := 0
		startTime := currObj.StartTime
		pos2 := currObj.Position
		if currObj.ObjectName == "slider" {
			startTime = currObj.EndTime
			pos2 = currObj.EndPosition
		}

		for j := i + 1; j < len(objs); j++ {
			nextObj := &objs[j]
			if nextObj.StartTime-stackThresh > startTime {
				break
			}

			if distance(nextObj.Position, currObj.Position) < 3 {
				currObj.stackHeight++
				startTime = nextObj.StartTime
				if nextObj.ObjectName == "slider" {
					startTime = nextObj.EndTime
				}
			}
			if distance(nextObj.Position, pos2) < 3 {
				sliderStack++
				currObj.stackHeight -= sliderStack
				startTime = nextObj.StartTime
				if nextObj.ObjectName == "slider" {
					startTime = nextObj.EndTime
				}
			}
		}
	}

	for i := 0; i < len(beatmap.HitObjects)-1; i++ {
		offset := float64(objs[i].stackHeight) * scale * -6.4
		beatmap.HitObjects[i].Position.X += offset
		beatmap.HitObjects[i].Position.Y += offset
	}
}

func diffRange(value float64) float64 {
	if value > 5.0 {
		return 1200 + (450-1200)*(value-5)/5
	} else if value < 5.0 {
		return 1200 - (1200-1800)*(5-value)/5
	}
	return 1200
}

func distance(v1, v2 parser.Point) float64 {
	x := v1.X - v2.X
	y := v1.Y - v2.Y
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}
