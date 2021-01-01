package handlers

import (
	mathcommands "maquiaBot/handlers/math-commands"

	"github.com/bwmarrin/discordgo"
)

// MathHandle handles commands that are regarding math
func MathHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		case "ave", "average", "mean":
			go mathcommands.Average(s, m)
		case "d", "dist", "distance", "dir", "direction":
			go mathcommands.DistanceDirection(s, m)
		case "dr", "degrad", "degreesradians":
			go mathcommands.DegreesRadians(s, m)
		case "rd", "raddeg", "radiansdegrees":
			go mathcommands.RadiansDegrees(s, m)
		case "stddev", "standarddev", "stddeviation", "standarddeviation":
			go mathcommands.StandardDeviation(s, m)
		case "va", "vadd", "vectora", "vectoradd":
			go mathcommands.VectorAdd(s, m)
		case "vc", "vcross", "vectorc", "vectorcross":
			go mathcommands.VectorCross(s, m)
		case "vd", "vdiv", "vdivide", "vectord", "vectordiv", "vectordivide":
			go mathcommands.VectorDivide(s, m)
		case "vdot", "vectordot":
			go mathcommands.VectorDot(s, m)
		case "vm", "vmult", "vmultiply", "vectorm", "vectormult", "vectormultiply":
			go mathcommands.VectorMultiply(s, m)
		case "vs", "vsub", "vsubtract", "vectors", "vectorsub", "vectorsubtract":
			go mathcommands.VectorSubtract(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "Please provide a valid math argument to execute! Check `help` for more details.")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please provide a math argument to execute! Check `help` for more details.")
	}
}
