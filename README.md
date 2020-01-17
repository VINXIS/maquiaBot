# maquiaBot
Discord bot that does a bunch of osu! stuff PROPERLY (sooner or later)

## Installation
 1. [Install golang](https://golang.org/doc/install) and [Tesseract](https://github.com/UB-Mannheim/tesseract/wiki). Ideally you have Go version 1.13 or newer. 
 2. Clone the repository using `git clone https://github.com/VINXIS/Twitter-Discord-Feed.git` to wherever you want.
 3. Go to the folder and open a console. Install the dependencies with `go get`.
 4. Go to the config folder and duplicate `config.example.json`. Name the duplicate `config.json` and fill in the twitter API credentials, and discord information
	 1. You can obtain twitter API credentials [here](https://developer.twitter.com/en/docs).
	 2. For the discord information. You add the discord bot token which is obtained from creating a discord bot [here](https://discordapp.com/developers/applications). Put the username as anything you want, preferably your bot's username, and put the avatar field to some image link, preferably the same image as your discord bot's.
Duplicate `config.example.json` in the config folder and call it `config.json` fill in all the slots.
 5. Invite the bot to your server by replacing `PUT_CLIENT_ID_HERE` in the URL below with the discord application's client ID obtained here [here](https://discordapp.com/developers/applications). https://discordapp.com/api/oauth2/authorize?client_id=PUT_CLIENT_ID_HERE&permissions=536870912&scope=bot.
7. Run the program by running `go build -o bot core.go` and then `./bot` in your instance / computer.

Inspiration from [owo](https://github.com/AznStevy/owo) and [BoatBot](https://github.com/0xg0ldpk3rx0/SupportBot)