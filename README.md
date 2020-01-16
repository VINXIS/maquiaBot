# maquiaBot
Discord bot that does a bunch of osu! stuff PROPERLY (sooner or later)

## Requirements
This project uses [go](https://golang.org/dl/). You will need to install [Tesseract](https://github.com/UB-Mannheim/tesseract/wiki) for image detection features, and all dependencies used.

Create a folder called data. In that folder, create a subfolder called osuFiles where .osu files will be stored. The plans are to minimize API calls, and to obtain information via the .osu files themselves instead after they are called once AND are ranked.

Duplicate `config.example.json` in the config folder and call it `config.json` fill in all the slots.

Inspiration from [owo](https://github.com/AznStevy/owo) and [BoatBot](https://github.com/0xg0ldpk3rx0/SupportBot)

cd to the directory and build using `go build -o maquiaBot core.go` and then run the program with `./maquiaBot`