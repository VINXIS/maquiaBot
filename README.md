# discord-osu
Discord bot that does a bunch of osu! stuff PROPERLY (sooner or later)

## Requirements
This project uses [go](https://golang.org/dl/). Alongside that, you will need to install [Tesseract](https://github.com/UB-Mannheim/tesseract/wiki) for image detection features.

Create a folder called data. In that folder, create a subfolder called osuFiles where .osu files will be stored. The plans are to minimize API calls, and to obtain information via the .osu files themselves instead after they are called once AND are ranked.
