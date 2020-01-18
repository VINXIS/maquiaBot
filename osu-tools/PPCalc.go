package osutools

import (
	"math"
	"strconv"

	osuapi "../osu-api"
)

// PPCalc calculates the pp given by the beatmap with specified acc and mods based off of https://github.com/ppy/osu/blob/master/osu.Game.Rulesets.Osu/Difficulty/OsuPerformanceCalculator.cs
func PPCalc(beatmap osuapi.Beatmap, score osuapi.Score, store chan<- string) {
	unrankable := 536870912 + 2048 + 4194304 + 8192 + 128

	if score.Mods&osuapi.Mods(unrankable) != 0 {
		store <- "0.00"
		return
	}

	multiplier := 1.12
	if score.Mods&osuapi.ModNoFail != 0 {
		multiplier *= 0.9
	}
	if score.Mods&osuapi.ModSpunOut != 0 {
		multiplier *= 0.95
	}

	aimPP := aimPP(beatmap, score)
	speedPP := speedPP(beatmap, score)
	accPP := accPP(beatmap, score)
	totalPP := multiplier * math.Pow(math.Pow(aimPP, 1.1)+math.Pow(speedPP, 1.1)+math.Pow(accPP, 1.1), 1.0/1.1)

	store <- strconv.FormatFloat(totalPP, 'f', 2, 64)
}

func convertSR(SR float64) float64 { return math.Pow(5.0*math.Max(1.0, SR/0.0675)-4.0, 3.0) / 100000.0 }

func aimPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	rawAim := beatmap.DifficultyAim
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)
	accuracy := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300))

	if score.Mods&osuapi.ModTouchDevice != 0 {
		rawAim = math.Pow(rawAim, 0.8)
	}

	aimValue := convertSR(rawAim)

	lengthBonus := 0.95 + 0.4*math.Min(1.0, totalHits/2000.0)
	if totalHits > 2000 {
		lengthBonus += math.Log10(totalHits/2000.0) / 2.0
	}
	aimValue *= lengthBonus

	aimValue *= math.Pow(0.94, float64(score.CountMiss))

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
	}

	aimValue *= 0.5 + accuracy/2.0
	aimValue *= 0.98 + math.Pow(beatmap.OverallDifficulty, 2.0)/2500

	return aimValue
}

func speedPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
	speedValue := convertSR(beatmap.DifficultySpeed)
	totalHits := float64(beatmap.Circles + beatmap.Sliders + beatmap.Spinners)
	accuracy := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300))

	lengthBonus := 0.95 + 0.4*math.Min(1.0, totalHits/2000.0)
	if totalHits > 2000 {
		lengthBonus += math.Log10(totalHits/2000.0) / 2.0
	}
	speedValue *= lengthBonus

	speedValue *= math.Pow(0.94, float64(score.CountMiss))

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

func accPP(beatmap osuapi.Beatmap, score osuapi.Score) float64 {
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
