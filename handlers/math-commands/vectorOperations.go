package mathcommands

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	mathtools "maquiaBot/math-tools"
)

// DistanceDirection gives the distance and direction between 2 points
func DistanceDirection(s *discordgo.Session, m *discordgo.MessageCreate) {
	VS, err := parseVectors(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	V1 := VS[0]
	V2 := VS[1]

	// Get distance and direction
	distance := mathtools.Distance(V1, V2)
	direction := mathtools.Direction(V1, V2)
	s.ChannelMessageSend(m.ChannelID, "The distance between "+V1.ToString()+" and "+V2.ToString()+" is **"+strconv.FormatFloat(distance, 'f', 2, 64)+"** and the direction is **"+direction.ToString()+"**")
}

// VectorAdd adds one vector from another
func VectorAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	VS, err := parseVectors(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	V1 := VS[0]
	V2 := VS[1]

	text := "V1: **" + V1.ToString() + "**\n"
	text += "V2: **" + V2.ToString() + "**\n\n"

	v1v2 := V1.Add(V2)
	text += "V1 + V2: **" + v1v2.ToString() + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// VectorSubtract subtracts one vector from another
func VectorSubtract(s *discordgo.Session, m *discordgo.MessageCreate) {
	VS, err := parseVectors(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	V1 := VS[0]
	V2 := VS[1]

	text := "V1: **" + V1.ToString() + "**\n"
	text += "V2: **" + V2.ToString() + "**\n\n"

	v1v2 := V1.Subtract(V2)
	v2v1 := V2.Subtract(V1)
	text += "V1 - V2: **" + v1v2.ToString() + "**\n"
	text += "V2 - V1: **" + v2v1.ToString() + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// VectorMultiply multiplies a scalar onto a vector
func VectorMultiply(s *discordgo.Session, m *discordgo.MessageCreate) {
	V, scalar, err := parseVectorScalar(m)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	text := "V: **" + V.ToString() + "**\n"
	text += "Scalar: **" + strconv.FormatFloat(scalar, 'f', 2, 64) + "**\n\n"

	VM := V.Multiply(scalar)
	text += "V * Scalar: **" + VM.ToString() + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// VectorDivide divides a vector by a scalar
func VectorDivide(s *discordgo.Session, m *discordgo.MessageCreate) {
	V, scalar, err := parseVectorScalar(m)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	text := "V: **" + V.ToString() + "**\n"
	text += "Scalar: **" + strconv.FormatFloat(scalar, 'f', 2, 64) + "**\n\n"

	VD := V.Divide(scalar)
	text += "V / Scalar: **" + VD.ToString() + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// VectorDot obtains the dot product of 2 vectors
func VectorDot(s *discordgo.Session, m *discordgo.MessageCreate) {
	VS, err := parseVectors(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	V1 := VS[0]
	V2 := VS[1]

	text := "V1: **" + V1.ToString() + "**\n"
	text += "V2: **" + V2.ToString() + "**\n\n"

	v1v2 := V1.Dot(V2)
	text += "V1 * V2: **" + strconv.FormatFloat(v1v2, 'f', 2, 64) + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// VectorCross obtains the dot product of 2 vectors
func VectorCross(s *discordgo.Session, m *discordgo.MessageCreate) {
	VS, err := parseVectors(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	V1 := VS[0]
	V2 := VS[1]

	text := "V1: **" + V1.ToString() + "**\n"
	text += "V2: **" + V2.ToString() + "**\n\n"

	v1v2 := V1.Cross(V2)
	v2v1 := V2.Cross(V1)
	text += "V1 x V2: **" + v1v2.ToString() + "**\n"
	text += "V2 x V1: **" + v2v1.ToString() + "**"
	s.ChannelMessageSend(m.ChannelID, text)
}

// UTILITY FUNCTIONS BELOW
// first parses a vector and scalar
// second parses 2 vectors

func parseVectorScalar(m *discordgo.MessageCreate) (V1 mathtools.Vector, s float64, parseError error) {
	vectorRegex, _ := regexp.Compile(`(?i)\((-?(\d|\.)+),\s*(-?(\d|\.)+)(,\s*(-?(\d|\.)+))?\)\s*(-?(\d|\.)+)`)

	if !vectorRegex.MatchString(m.Content) {
		return nil, 0, errors.New("please give 2 points with x and y (and z if 3 dimensions) coordinates in parentheses; refer to `help distance` for more info")
	}

	points := vectorRegex.FindStringSubmatch(m.Content)

	// 1st vector
	X1, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return nil, 0, errors.New("error in parsing the x coordinate")
	}
	Y1, err := strconv.ParseFloat(points[3], 64)
	if err != nil {
		return nil, 0, errors.New("error in parsing the y coordinate")
	}
	Z1 := 0.0
	if points[6] != "" {
		Z1, err = strconv.ParseFloat(points[6], 64)
		if err != nil {
			Z1 = 0
		}
	}

	if Z1 == 0 {
		V1 = mathtools.Vector2D{X1, Y1}
	} else {
		V1 = mathtools.Vector3D{mathtools.Vector2D{X1, Y1}, Z1}
	}

	s, err = strconv.ParseFloat(points[8], 64)
	if err != nil {
		return nil, 0, errors.New("error in parsing the scalar value")
	}
	return V1, s, nil
}

func parseVectors(m *discordgo.MessageCreate) ([]mathtools.Vector, error) {
	vectorRegex, _ := regexp.Compile(`(?i)\((-?(\d|\.)+)\s*,\s*(-?(\d|\.)+)(\s*,\s*(-?(\d|\.)+))?\)\s*\((-?(\d|\.)+)\s*,\s*(-?(\d|\.)+)(\s*,\s*(-?(\d|\.)+))?\)`)

	if !vectorRegex.MatchString(m.Content) {
		return nil, errors.New("please give 2 points with x and y (and z if 3 dimensions) coordinates in parentheses; refer to `help distance` for more info")
	}

	points := vectorRegex.FindStringSubmatch(m.Content)
	var (
		V1 mathtools.Vector
		V2 mathtools.Vector
	)

	// 1st vector
	X1, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return nil, errors.New("error in parsing the first x coordinate")
	}
	Y1, err := strconv.ParseFloat(points[3], 64)
	if err != nil {
		return nil, errors.New("error in parsing the first y coordinate")
	}
	Z1 := 0.0
	if points[6] != "" {
		Z1, err = strconv.ParseFloat(points[6], 64)
		if err != nil {
			Z1 = 0
		}
	}

	if Z1 == 0 {
		V1 = mathtools.Vector2D{X1, Y1}
	} else {
		V1 = mathtools.Vector3D{mathtools.Vector2D{X1, Y1}, Z1}
	}

	// 2nd vector
	X2, err := strconv.ParseFloat(points[8], 64)
	if err != nil {
		return nil, errors.New("error in parsing the second x coordinate")
	}
	Y2, err := strconv.ParseFloat(points[10], 64)
	if err != nil {
		return nil, errors.New("error in parsing the second y coordinate")
	}
	Z2 := 0.0
	if points[13] != "" {
		Z2, err = strconv.ParseFloat(points[13], 64)
		if err != nil {
			Z2 = 0
		}
	}

	if Z2 == 0 {
		V2 = mathtools.Vector2D{X2, Y2}
	} else {
		V2 = mathtools.Vector3D{mathtools.Vector2D{X2, Y2}, Z2}
	}
	return []mathtools.Vector{V1, V2}, nil
}
