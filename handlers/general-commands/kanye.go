package gencommands

import (
	"crypto/rand"
	"math/big"

	"github.com/bwmarrin/discordgo"
)

// Kanye gives a funny line Kanye said in his songs
func Kanye(s *discordgo.Session, m *discordgo.MessageCreate) {
	KanyeLines := []string{
		"```\nGetting stupid ass straight out the jar / Pockets on Shrek, rockets on deck / Tell me what’s next, alien sex? / I’ma disrobe you, then I’ma probe you``` - Kanye West on ET",
		"```\nI am a God``` - Kanye West on I Am A God",
		"```\nHurry up with my damn massage / Hurry up with my damn ménage / Hurry up with my damn croissants``` - Kanye West on I Am A God",
		"```\nI just talked to Jesus, he said: what up Yeezus``` - Kanye West on I Am A God",
		"```\nPussy had me floatin', felt like Deepak Chopra / Pussy had me DEAD, might call Tupac over``` - Kanye West on Hold My Liquor",
		"```\nMy apologies are you into astrology 'cause, um, I'm trying to make it to Uranus``` - Kanye West on Gettin' It In",
		"```\nYour titties, let em out, free at last / Thank God almighty, they free at last``` - Kanye West on I'm In It",
		"```\nBlack Timbs all on your couch again / Black dick all in your spouse again``` - Kanye West on On Sight",
		"```\nI keep it 300, like the Romans / 300 bitches, where's the Trojans?``` - Kanye West on Black Skinhead",
		"```\nNow if I fuck this model / And she just bleached her asshole / And I get bleach on my T-shirt / I'mma feel like an asshole``` - Kanye West on Father Stretch My Hands Pt. 1",
		"```\nEatin' Asian pussy, all I need was sweet and sour sauce``` - Kanye West on I'm In It",
		"```\nNo more drugs for me, pussy and religion is all I need``` - Kanye West on Hell of a Life",
		"```\nOne day I’m gon' marry a porn star``` - Kanye West on Hell of a Life",
		"```\nHead of the class and she just want a swallow-ship``` - Kanye West on Monster",
		"```\nEverybody know I'm a motherfucking monster``` - Kanye West on Monster",
		"```\nI'm like a fly Malcolm X / Buy any jeans necessary``` - Kanye West on Good Morning",
		"```\nCause the same people that tried to black ball me / Forgot about two things, my black balls``` - Kanye West on Gorgeous",
		"```\nI don't need your pussy, bitch I'm on my own dick``` - Kanye West on POWER",
		"```\nAy, none of us would be here without cum``` - Kanye West on All Mine",
		"```\nI'm a sick fuck, I like a quick fuck``` - Kanye West on I Love It",
		"```\nI miss the old Kanye``` - Kanye West on I Love Kanye",
		"```\nPoopy-di scoop``` - Kanye West on Lift Yourself",
		"```\nYou know how many girls I took to the titty shop?``` - Kanye West on Yikes",
		"```\npremeditated murder``` - Kanye West on I Thought About Killing You",
		"```\nDrug dealer buy Jordans, crackhead buy crack / And a white man get paid off of all of that.``` - Kanye West on All Falls Down",
		"```\nCouldn't afford a car so she named her daughter Alexus``` - Kanye West on All Falls Down",
		"```\nThey be asking us questions, harass and arrest us / Saying 'we eat pieces of shit like you for breakfast' / Huh? Y'all eat pieces of shit? What's the basis?``` - Kanye West on Jesus Walks",
		"```\nNow even though I went to college and dropped out of school quick / I always had a Ph.D.: a Pretty Huge Dick``` - Kanye West on Breathe In, Breathe Out",
		"```\nSo I live by two words: 'fuck you, pay me'``` - Kanye West on Two Words",
		"```\nClosed on Sunday, you my Chick-fil-A``` - Kanye West on Closed On Sunday",
		"```\nLost in translation with a whole fuckin' nation / They say I was the abomination of Obama's nation / Well that's a pretty bad way to start the conversation``` - Kanye West on POWER",
		"```\nBa-ba-ba-ba-bwa-ba / Ga-ba-ba-ba / Rude-rude-rude-rude-woo!``` - Kanye West on Feel the Love",
		"```\nGet down girl, go 'head, get down``` - Kanye West on Gold Digger",
		"```\nMayonnaise colored Benz, I push Miracle Whips``` - Kanye West on Wack Niggaz",
		"```\nYou should be honored by my lateness / That I would even show up to this fake shit``` - Kanye West on Stronger",
	}
	roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(KanyeLines))))
	s.ChannelMessageSend(m.ChannelID, KanyeLines[roll.Int64()])
}
