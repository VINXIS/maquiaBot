package osuapi

// Mods in the game.
const (
	ModNoFail Mods = 1 << iota
	ModEasy
	ModTouchDevice
	ModHidden
	ModHardRock
	ModSuddenDeath
	ModDoubleTime
	ModRelax
	ModHalfTime
	ModNightcore
	ModFlashlight
	ModAutoplay
	ModSpunOut
	ModAutopilot
	ModPerfect
	ModKey4
	ModKey5
	ModKey6
	ModKey7
	ModKey8
	ModFadeIn
	ModRandom
	ModCinema
	ModTarget
	ModKey9
	ModKeyCoop
	ModKey1
	ModKey3
	ModKey2
	ModScoreV2
	ModFreeModAllowed    = ModNoFail | ModEasy | ModHidden | ModHardRock | ModSuddenDeath | ModFlashlight | ModFadeIn | ModRelax | ModAutopilot | ModSpunOut | ModKeyMod
	ModKeyMod            = ModKey1 | ModKey2 | ModKey3 | ModKey4 | ModKey5 | ModKey6 | ModKey7 | ModKey8 | ModKey9 | ModKeyCoop
	ModScoreIncreaseMods = ModHidden | ModHardRock | ModDoubleTime | ModNightcore | ModFlashlight | ModFadeIn
)

// Mods is a bitwise enum of mods used in a score.
//
// Mods may appear complicated to use for a beginner programmer. Fear not!
// This is how hard they can get for creation of a mod combination:
//
//    myModCombination := osuapi.ModHardRock | osuapi.ModDoubleTime | osuapi.ModHidden | osuapi.ModSpunOut
//
// As for checking that an existing mod comination is enabled:
//
//    if modCombination&osuapi.ModHardRock != 0 {
//        // HardRock is enabled
//    }
//
// To learn more about bitwise operators, have a look at it on wikipedia:
// https://en.wikipedia.org/wiki/Bitwise_operation#Bitwise_operators
type Mods int

var modsString = [...]string{
	"NF",
	"EZ",
	"NV",
	"HD",
	"HR",
	"SD",
	"DT",
	"RX",
	"HT",
	"NC",
	"FL",
	"AU", // Auto.
	"SO",
	"AP", // Autopilot.
	"PF",
	"K4",
	"K5",
	"K6",
	"K7",
	"K8",
	"FI",
	"RN",
	"CN",
	"TR",
	"K9",
	"KC",
	"K1",
	"K3",
	"K2",
	"V2",
}

// ParseMods parse a string with mods in the format "HDHRDT"
func ParseMods(mods string) (m Mods) {
	if len(mods) == 0 || len(mods)%2 != 0 {
		return
	}
	modsSl := make([]string, len(mods)/2)
	for n, modPart := range mods {
		modsSl[n/2] += string(modPart)
	}
	for _, mod := range modsSl {
		for index, availableMod := range modsString {
			if availableMod == mod {
				m |= 1 << uint(index)
				break
			}
		}
	}
	return
}

func (m Mods) String() (s string) {
	if m == 0 {
		s = "NM"
	}
	for i := 0; i < len(modsString); i++ {
		activated := 1&m == 1
		if activated {
			s += modsString[i]
		}
		m >>= 1
	}
	return
}
