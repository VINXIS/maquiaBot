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
		"```\nYour girl don't like me, how long has she been gay?``` - Kanye West on Bring Me Down",
		"```\nAnd tell him that your mama had a fattie / He looked up at me / Said, 'Daddy, that's the reason why you had me?' / Yep, we was practicing / 'Til one day your ass bust through the packaging / You know what though? You my favorite accident``` - Kanye West on Celebration",
		"```\nOld folks talking 'bout, 'Back in my day' / But homie this is my day / Class started 2 hours ago, oh am I late?``` - Kanye West on Can't Tell Me Nothing",
		"```\nHey, you remember where we first met? / Okay, I don't remember where we first met``` - Kanye West on Bound 2",
		"```\nI'm living in the future so the present is my past / My presence is a present, kiss my ass``` - Kanye West on Monster",
		"```\nIma fix wolves``` - Kanye West on Twitter",
		"```\nI specifically ordered persian rugs with cherub imagery!!! What do I have to do to get a simple persian rug with cherub imagery uuuuugh``` - Kanye West on Twitter",
		"```\n...on another note, can brah be the girl version of bruh???``` - Kanye West on Twitter",
		"```\nI'M SO HYPE RIGHT NOW  EVERYTHING HAS CHANGED ... HAVE YA'LL EVER SEEN TRON? THE END OF THE TRON WHERE EVERYTHING LIGHT UP!!!!``` - Kanye West on Twitter",
		"```\nMark Zuckerberg invest 1 billion dollars into Kanye West ideas``` - Kanye West on Twitter",
		"```\nYou have distracted from my creative process``` - Kanye West on Twitter",
		"```\nI hate when I'm on a flight and I wake up with a water bottle next to me like oh great now I gotta be responsible for this water bottle``` - Kanye West on Twitter",
		"```\nI no longer have a manager. I can't be managed``` - Kanye West on Twitter",
		"```\nPlease no one text me or ask anything till Monday``` - Kanye West on Twitter",
		"```\nburn that excel spread sheet 🔥😂``` - Kanye West on Twitter",
		"```\nSometimes I push the door close button on people running towards the elevator. I just need my own elevator sometimes. My sanctuary.``` - Kanye West on Twitter",
		"```\nShe asked when is fashion week.... uuuum...I thought it was every week??!!``` - Kanye West on Twitter",
		"```\nyou may be talented, but you're not kanye west``` - Kanye West on Twitter",
		"```\nI wish I had a friend like me.``` - Kanye West on Twitter",
		"```\nim the kanye best hahaha``` - Kanye West on Twitter",
		"```\nSuper inspired by my visit to Ikea today , really amazing company... my mind is racing with the possibilities...``` - Kanye West on Twitter",
		"```\nI need a room full of mirrors so I can be surrounded by winners.``` - Kanye West on Twitter",
		"```\nI understand they you don't like me but I need you to understand that I don't care.``` - Kanye West on Twitter",
		"```\nI forgor 💀``` - Kanye West on Twitter",
		"```\nCome and get me ... this is the exodus``` - Kanye West on Twitter",
		"```\nJunya Watanabe on my WRI``` - Kanye West on Junya",
		"```\nMan, it's too early / What the hell you doin' wakin' me up at 5:30? / Why the hell are you worried? / Play somethin' that is very, very vibe-worthy``` - Kanye West on Believe What I Say",
	}
	roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(KanyeLines))))
	s.ChannelMessageSend(m.ChannelID, KanyeLines[roll.Int64()])
}
