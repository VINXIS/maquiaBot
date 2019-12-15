package osutools

// ScoreRank gives you the score rank for a score.
func ScoreRank(percent50, percent300 float64, miss int, silver bool) (scoreRank string) {
	switch {
	case percent300 == 1:
		scoreRank = "SS"
	case percent300 > 0.9 && percent50 < 0.01 && miss == 0:
		scoreRank = "S"
	case percent300 > 0.8 && miss == 0, percent300 > 0.9:
		scoreRank = "A"
	case percent300 > 0.7 && miss == 0, percent300 > 0.8:
		scoreRank = "B"
	case percent300 > 0.6:
		scoreRank = "C"
	default:
		scoreRank = "D"
	}
	if silver && (scoreRank == "S" || scoreRank == "SS") {
		scoreRank += "H"
	}
	return scoreRank
}
