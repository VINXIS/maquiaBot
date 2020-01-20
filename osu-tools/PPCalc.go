package osutools

import (
	"math"
	"strconv"
	"strings"

	osuapi "../osu-api"
)

// PPCalc calculates the pp given by the beatmap with specified acc and mods based off of https://github.com/ppy/osu/blob/master/osu.Game.Rulesets.Osu/Difficulty/OsuPerformanceCalculator.cs
func PPCalc(beatmap osuapi.Beatmap, score osuapi.Score, store chan<- string) {
	unrankable := osuapi.Mods(536870912 + 2048 + 4194304 + 8192 + 128)

	if score.Mods&unrankable != 0 {
		store <- "0.00"
		return
	}

	totalPP := 0.00
	switch beatmap.Mode {
	case osuapi.ModeOsu:
		multiplier := 1.12
		if score.Mods&osuapi.ModNoFail != 0 {
			multiplier *= 0.9
		}
		if score.Mods&osuapi.ModSpunOut != 0 {
			multiplier *= 0.95
		}

		aimPP := aimPP(beatmap, score)
		speedPP := speedPP(beatmap, score)
		accPP := accSTDPP(beatmap, score)
		totalPP = multiplier * math.Pow(math.Pow(aimPP, 1.1)+math.Pow(speedPP, 1.1)+math.Pow(accPP, 1.1), 1.0/1.1)
	case osuapi.ModeTaiko:
		multiplier := 1.1
		if score.Mods&osuapi.ModNoFail != 0 {
			multiplier *= 0.9
		}
		if score.Mods&osuapi.ModHidden != 0 {
			multiplier *= 1.1
		}

		strainPP := strainTKOPP(beatmap, score)
		accPP := accTKOPP(beatmap, score)
		totalPP = multiplier * math.Pow(math.Pow(strainPP, 1.1)+math.Pow(accPP, 1.1), 1.0/1.1)
	case osuapi.ModeOsuMania:
		scoreMultiplier := maniaModCheck(score)
		score.Score *= int64(1.0 / scoreMultiplier)

		multiplier := 0.8
		if score.Mods&osuapi.ModNoFail != 0 {
			multiplier *= 0.9
		}
		if score.Mods&osuapi.ModEasy != 0 {
			multiplier *= 0.5
		}

		strainPP := strainMANPP(beatmap, score)
		accPP := accMANPP(beatmap, score, strainPP)
		totalPP = multiplier * math.Pow(math.Pow(strainPP, 1.1)+math.Pow(accPP, 1.1), 1.0/1.1)
	case osuapi.ModeCatchTheBeat:
		totalPP = catchPP(beatmap, score)
	}

	store <- strconv.FormatFloat(totalPP, 'f', 2, 64)
}

// osu!standard functions
func convertSTDSR(SR float64) float64 {
	return math.Pow(5.0*math.Max(1.0, SR/0.0675)-4.0, 3.0) / 100000.0
}

func aimPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	rawAim := beatmap.DifficultyAim
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)
	accuracy := float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))

	if score.Mods&osuapi.ModTouchDevice != 0 {
		rawAim = math.Pow(rawAim, 0.8)
	}

	aimValue := convertSTDSR(rawAim)

	lengthBonus := 0.95 + 0.4*math.Min(1.0, totalHits/2000.0)
	if totalHits > 2000 {
		lengthBonus += math.Log10(totalHits/2000.0) * 0.5
	}
	aimValue *= lengthBonus

	aimValue *= math.Pow(0.97, float64(score.CountMiss))

	if beatmap.MaxCombo > 0 {
		aimValue *= math.Min(math.Pow(float64(score.MaxCombo), 0.8)/math.Pow(float64(beatmap.MaxCombo), 0.8), 1.0)
	}

	ARBonus := 1.0
	if beatmap.ApproachRate > 10.33 {
		ARBonus += 0.3 * (beatmap.ApproachRate - 10.33)
	} else if beatmap.ApproachRate < 8.0 {
		ARBonus += 0.01 * (8.0 - beatmap.ApproachRate)
	}
	aimValue *= ARBonus

	if score.Mods&osuapi.ModHidden != 0 {
		aimValue *= 1.0 + 0.04*(12.0-beatmap.ApproachRate)
	}

	if score.Mods&osuapi.ModFlashlight != 0 {
		FLBonus := 1.0 + 0.35*math.Min(1.0, totalHits/200.0)
		if totalHits > 200 {
			FLBonus += 0.3 * math.Min(1.0, (totalHits-200)/300.0)
			if totalHits > 500 {
				FLBonus += (totalHits - 500) / 1200.0
			}
		}
		aimValue *= FLBonus
	}

	aimValue *= 0.5 + accuracy/2.0
	aimValue *= 0.98 + math.Pow(beatmap.OverallDifficulty, 2.0)/2500

	return aimValue
}

func speedPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	speedValue := convertSTDSR(beatmap.DifficultySpeed)
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)
	accuracy := float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))

	lengthBonus := 0.95 + 0.4*math.Min(1.0, totalHits/2000.0)
	if totalHits > 2000 {
		lengthBonus += math.Log10(totalHits/2000.0) * 0.5
	}
	speedValue *= lengthBonus

	speedValue *= math.Pow(0.97, float64(score.CountMiss))

	if beatmap.MaxCombo > 0 {
		speedValue *= math.Min(math.Pow(float64(score.MaxCombo), 0.8)/math.Pow(float64(beatmap.MaxCombo), 0.8), 1.0)
	}

	ARBonus := 1.0
	if beatmap.ApproachRate > 10.33 {
		ARBonus += 0.3 * (beatmap.ApproachRate - 10.33)
	}
	speedValue *= ARBonus

	if score.Mods&osuapi.ModHidden != 0 {
		speedValue *= 1.0 + 0.04*(12.0-beatmap.ApproachRate)
	}

	speedValue *= 0.02 + accuracy
	speedValue *= 0.96 + math.Pow(beatmap.OverallDifficulty, 2.0)/1600

	return speedValue
}

func accSTDPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	trueAcc := 0.0
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)

	if beatmap.Circles > 0 {
		trueAcc = ((float64(score.Count300)-(totalHits-float64(beatmap.Circles)))*6.0 + float64(score.Count100)*2.0 + float64(score.Count50)) / (float64(beatmap.Circles) * 6.0)
	}

	if trueAcc < 0 {
		trueAcc = 0
	}

	accValue := math.Pow(1.52163, beatmap.OverallDifficulty) * math.Pow(trueAcc, 24.0) * 2.83
	accValue *= math.Min(1.15, math.Pow(float64(beatmap.Circles)/1000.0, 0.3))

	if score.Mods&osuapi.ModHidden != 0 {
		accValue *= 1.08
	}
	if score.Mods&osuapi.ModFlashlight != 0 {
		accValue *= 1.02
	}

	return accValue
}

// osu!taiko functions
func strainTKOPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)
	accuracy := float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))

	strainValue := math.Pow(5.0*math.Max(1.0, beatmap.DifficultyRating/0.0075)-4.0, 2.0) / 100000.0

	lengthBonus := 1.0 + 0.1*math.Min(1.0, totalHits/1500.0)
	strainValue *= lengthBonus

	strainValue *= math.Pow(0.985, float64(score.CountMiss))

	if beatmap.MaxCombo > 0 {
		strainValue *= math.Min(math.Pow(float64(score.MaxCombo), 0.5)/math.Pow(float64(beatmap.MaxCombo), 0.5), 1.0)
	}

	if score.Mods&osuapi.ModHidden != 0 {
		strainValue *= 1.025
	}

	if score.Mods&osuapi.ModFlashlight != 0 {
		strainValue *= 1.05 * lengthBonus
	}

	return strainValue * accuracy
}

func accTKOPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	accuracy := float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)

	clockRate := 1.0
	OD := beatmap.OverallDifficulty
	if score.Mods&osuapi.ModDoubleTime != 0 || score.Mods&osuapi.ModNightcore != 0 {
		clockRate = 1.5
	} else if score.Mods&osuapi.ModHalfTime != 0 {
		clockRate = 0.75
	}
	ODScale := (80.0 - 6.0*OD) * clockRate
	OD = (80.0 - ODScale) / 6.0

	greatResult := 70 // greatResult calculation numbers obtained from here https://github.com/ppy/osu/blob/master/osu.Game/Rulesets/Scoring/HitWindows.cs
	if OD > 5.0 {
		greatResult = int(70 + (40-70)*(OD-5)/5)
	} else if OD < 5.0 {
		greatResult = int(70 - (70-100)*(5-OD)/5)
	}
	greatResult /= 2

	accValue := math.Pow(150.0/(float64(greatResult)/clockRate), 1.1) * math.Pow(accuracy, 15.0) * 22.0
	return accValue * math.Min(1.15, math.Pow(totalHits/1500.0, 0.3))
}

// osu!mania functions
func strainMANPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)

	strainValue := math.Pow(5.0*math.Max(1.0, beatmap.DifficultyRating/0.2)-4.0, 2.2) / 135.0

	strainValue *= 1.0 + 0.1*math.Min(1.0, totalHits/1500.0)

	if score.Score <= 500000 {
		strainValue = 0
	} else if score.Score <= 600000 {
		strainValue *= (float64(score.Score) - 500000) / 100000 * 0.3
	} else if score.Score <= 700000 {
		strainValue *= 0.3 + (float64(score.Score)-600000)/100000*0.25
	} else if score.Score <= 800000 {
		strainValue *= 0.55 + (float64(score.Score)-700000)/100000*0.2
	} else if score.Score <= 900000 {
		strainValue *= 0.75 + (float64(score.Score)-800000)/100000*0.15
	} else {
		strainValue *= 0.9 + (float64(score.Score)-900000)/100000*0.1
	}

	return strainValue
}

func accMANPP(beatmap osuapi.Beatmap, score osuapi.Score, strainValue float64) float64 {
	clockRate := 1.0
	OD := beatmap.OverallDifficulty
	if score.Mods&osuapi.ModDoubleTime != 0 || score.Mods&osuapi.ModNightcore != 0 {
		clockRate = 1.5
	} else if score.Mods&osuapi.ModHalfTime != 0 {
		clockRate = 0.75
	}
	ODScale := (80.0 - 6.0*OD) * clockRate
	OD = (80.0 - ODScale) / 6.0

	greatResult := 98 // greatResult calculation numbers obtained from here https://github.com/ppy/osu/blob/master/osu.Game/Rulesets/Scoring/HitWindows.cs
	if OD > 5.0 {
		greatResult = int(98 + (68-98)*(OD-5)/5)
	} else if OD < 5.0 {
		greatResult = int(98 - (98-128)*(5-OD)/5)
	}
	greatResult /= 2

	return math.Max(0, 0.2-(float64(greatResult)/clockRate-34)*0.006667) * strainValue * math.Pow(math.Max(0, float64(score.Score)-960000)/40000, 1.1)
}

func maniaModCheck(score osuapi.Score) float64 {
	val := 1.0
	ignoreMods := []string{"NM", "HR", "SD", "PF", "DT", "NC", "FI", "HD", "FL"}

	mods := strings.Split(score.Mods.String(), "")
	for i := 0; i < len(mods); i += 2 {
		exists := false
		for _, mod := range ignoreMods {
			if mod == mods[i]+mods[i+1] {
				exists = true
				break
			}
		}

		if !exists {
			val *= 0.5
		}
	}
	return val
}

// osu!fruits functions
func catchPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	accuracy := float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)

	value := math.Pow(5.0*math.Max(1.0, beatmap.DifficultyRating/0.0049)-4.0, 2.0) / 100000.0

	lengthBonus := 0.95 + 0.4*math.Min(1.0, totalHits/3000.0)
	if totalHits > 3000 {
		lengthBonus += math.Log10(totalHits/3000.0) * 0.5
	}
	value *= lengthBonus

	value *= math.Pow(0.97, float64(score.CountMiss))

	if beatmap.MaxCombo > 0 {
		value *= math.Min(math.Pow(float64(score.MaxCombo), 0.8)/math.Pow(float64(beatmap.MaxCombo), 0.8), 1.0)
	}

	ARBonus := 1.0
	if beatmap.ApproachRate > 9 {
		ARBonus += 0.1 * (beatmap.ApproachRate - 9.0)
	} else if beatmap.ApproachRate < 8.0 {
		ARBonus += 0.025 * (8.0 - beatmap.ApproachRate)
	}
	value *= ARBonus

	if score.Mods&osuapi.ModHidden != 0 {
		value *= 1.35
	}

	value *= math.Pow(accuracy, 5.5)

	if score.Mods*osuapi.ModNoFail != 0 {
		value *= 0.9
	}

	return value
}
