package osutools

import (
	"image"
	"image/color"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	osuapi "../osu-api"
	structs "../structs"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// ResultImage creates a result image for a score
func ResultImage(score osuapi.Score, beatmap osuapi.Beatmap, user osuapi.User, replay structs.ReplayData) (image.Image, error) {
	res, err := http.Get("https://assets.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapSetID) + "/covers/raw.jpg")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		f, err := os.Open("./osu-images/default.png")
		if err != nil {
			return nil, err
		}

		img, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
	}
	imgBounds := img.Bounds()
	if float64(imgBounds.Dx())/float64(imgBounds.Dy()) != 16/9 {
		if float64(imgBounds.Dx()) < float64(imgBounds.Dy())*16/9 {
			img = imaging.CropAnchor(img, imgBounds.Dx(), int(float64(imgBounds.Dx())*9/16), imaging.Center)
		} else {
			img = imaging.CropAnchor(img, int(float64(imgBounds.Dy())*16/9), imgBounds.Dy(), imaging.Center)
		}
	}
	img = imaging.Resize(img, 1920, 1080, imaging.Lanczos)

	ctx := gg.NewContextForImage(img)
	bounds := ctx.Image().Bounds()

	// Draw dark overlay
	ctx.SetRGBA(0, 0, 0, 0.35)
	ctx.DrawRectangle(0, 0, float64(bounds.Dx()), float64(bounds.Dy()))
	ctx.Fill()

	// Scaling based off of 1920 x 1080 sizes
	xScale := float64(bounds.Dx()) / 1920
	yScale := float64(bounds.Dy()) / 1080
	imgScale := yScale * 62
	fontScale := yScale * 112

	// Paths
	font := "./fonts/Aller-Light.ttf"

	// Draw results panel rectangle
	ctx.SetRGBA(0, 0, 0, 0.95)
	ctx.DrawRectangle(xScale*14, yScale*157, xScale*918, yScale*683)
	ctx.Fill()

	// Write score
	ctx.SetRGB255(234, 234, 234)
	ctx.LoadFontFace(font, fontScale*1.2)
	scoreText := strconv.FormatInt(score.Score, 10)
	for len(scoreText) < 8 {
		scoreText = "0" + scoreText
	}
	ctx.DrawStringAnchored(scoreText, xScale*(918*0.53+14), yScale*(210), 0.5, 0.5)

	// Write 300 100 50
	writeCount(strconv.Itoa(score.Count300), font, ctx, xScale*180, yScale*400, fontScale)
	writeCount(strconv.Itoa(score.Count100), font, ctx, xScale*180, yScale*525, fontScale)
	writeCount(strconv.Itoa(score.Count50), font, ctx, xScale*180, yScale*650, fontScale)

	// Add 100 50 images
	writeImage("./osu-images/100.png", ctx, int(xScale*45), int(yScale*525), imgScale) // 100
	writeImage("./osu-images/50.png", ctx, int(xScale*45), int(yScale*650), imgScale)  // 50

	// Write Geki Katu Miss
	writeCount(strconv.Itoa(score.CountGeki), font, ctx, xScale*635, yScale*400, fontScale)
	writeCount(strconv.Itoa(score.CountKatu), font, ctx, xScale*635, yScale*525, fontScale)
	writeCount(strconv.Itoa(score.CountMiss), font, ctx, xScale*635, yScale*650, fontScale)

	// Add Katu Miss images
	writeImage("./osu-images/katu.png", ctx, int(xScale*515), int(yScale*525), imgScale) // Katu
	writeImage("./osu-images/miss.png", ctx, int(xScale*515), int(yScale*650), imgScale) // Miss

	// Write Combo and accuracy
	writeCount(strconv.Itoa(score.MaxCombo), font, ctx, xScale*49, yScale*800, fontScale)
	ctx.LoadFontFace(font, fontScale) // Reload font
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	acc := strconv.FormatFloat(accCalc, 'f', 2, 64) + "%"
	ctx.DrawStringAnchored(acc, xScale*(918/2+14), yScale*800, 0, 0)

	// Put the Wheel thing and the grade icon
	f, err := os.Open("./osu-images/wheel.png")
	if err == nil {
		img, _, err := image.Decode(f)
		if err == nil {
			scale := yScale * 950 / float64(img.Bounds().Dy())
			img = imaging.Fit(img, int(scale*float64(img.Bounds().Dx())), int(scale*float64(img.Bounds().Dy())), imaging.Lanczos)
			img = imaging.Rotate(img, rand.Float64()*360, color.Transparent)
			ctx.DrawImageAnchored(img, int(xScale*1510), int(yScale*460), 0.5, 0.5)
		}
	}
	writeImage("./osu-images/"+score.Rank+".png", ctx, int(xScale*1375), int(yScale*760), yScale*600)

	// Draw upper rectangle now so it's above the wheel thing
	ctx.SetRGBA(0, 0, 0, 0.8)
	ctx.DrawRectangle(0, 0, float64(bounds.Dx()), float64(bounds.Dy())/8)
	ctx.Fill()

	// Write upper text
	ctx.SetRGB255(234, 234, 234)
	ctx.LoadFontFace(font, yScale*40)
	// first line
	if beatmap.Artist == "" {
		ctx.DrawStringAnchored(beatmap.Title+" ["+beatmap.DiffName+"]", xScale*7, yScale*40, 0, 0)
	} else {
		ctx.DrawStringAnchored(beatmap.Artist+" - "+beatmap.Title+" ["+beatmap.DiffName+"]", xScale*7, yScale*40, 0, 0)
	}

	ctx.LoadFontFace(font, yScale*30)
	ctx.DrawStringAnchored("Beatmap by "+beatmap.Creator, xScale*7, yScale*76, 0, 0)                                                                     // second line
	ctx.DrawStringAnchored("Played by "+user.Username+" on "+score.Date.GetTime().UTC().Format("2006-01-02 3:04:05 PM")+".", xScale*7, yScale*108, 0, 0) // third line

	// Draw watch rectangle and write watch
	ctx.LoadFontFace(font, fontScale)
	if score.Replay {
		ctx.SetRGBA(0, 0, 0, 0.5)
		ctx.DrawRectangle(xScale*1344, yScale*770, xScale*(1920-1344), yScale*120)
		ctx.Fill()

		ctx.SetRGBA(0.917647059, 0.917647059, 0.917647059, 0.5)
		w3, _ := ctx.MeasureString("Watch")
		ctx.DrawStringAnchored("Watch", xScale*1344+(xScale*(1920-1344)-w3)/2, yScale*770+yScale*60, 0, 0.5)
	}

	// Draw mods
	rawMods := score.Mods.String()
	rawMods = strings.Replace(rawMods, "DTNC", "NC", -1)
	mods := strings.Split(rawMods, "")
	i := 0
	x := float64(1831)
	for i < len(mods)-1 {
		f, err := os.Open("./osu-images/" + mods[i] + mods[i+1] + ".png")
		if err == nil {
			modImg, _, err := image.Decode(f)
			if err == nil {
				modY := yScale * 62
				modX := modY * 88 / 62
				modImg = imaging.Resize(modImg, int(1.5*modX), int(1.5*modY), imaging.Lanczos)
				ctx.DrawImageAnchored(modImg, int(xScale*x), int(yScale*590), 0.5, 0.5)
			}
		}
		x -= 50
		i += 2
	}

	// Performance Graph
	ctx.SetRGBA(0, 0, 0, 0.89)
	ctx.DrawRectangle(xScale*369, yScale*864, xScale*427, yScale*196)
	ctx.Fill()
	ctx.SetRGB255(234, 234, 234)
	if score.FullCombo {
		ctx.DrawStringAnchored("Perfect", xScale*(369+427/2), yScale*(864+196/2), 0.5, 0.5)
	}

	// Lifebar
	if score.Replay {
		ctx.SetRGB255(150, 204, 46)
		ctx.SetLineWidth(yScale * 4)
		ctx.DrawLine(xScale*369+xScale*3, yScale*865, xScale*369+3+xScale*150*(rand.Float64()*0.75+0.25), yScale*865)
		ctx.Stroke()
	}

	// Cursor
	cursorX := rand.Float64()*xScale*327 + xScale*469
	cursorY := rand.Float64()*yScale*109 + yScale*900
	boxX := cursorX + xScale*10
	boxY := cursorY + yScale*10
	if !score.Replay {
		cursorX = rand.Float64() * xScale * 1920
		cursorY = rand.Float64() * yScale * 1080
	}
	f, err = os.Open("./osu-images/cursor.png")
	if err == nil {
		img, _, err := image.Decode(f)
		if err == nil {
			scale := yScale * 100 / float64(img.Bounds().Dy())
			img = imaging.Fit(img, int(scale*float64(img.Bounds().Dx())), int(scale*float64(img.Bounds().Dy())), imaging.Lanczos)
			ctx.DrawImageAnchored(img, int(cursorX), int(cursorY), 0.5, 0.5)
		}
	}

	// UR stuff
	if score.Replay {
		// Border and Box
		ctx.SetRGB255(54, 54, 54)
		ctx.DrawRoundedRectangle(boxX-1, boxY-1, xScale*219, yScale*65, yScale*4)
		ctx.Fill()
		ctx.SetRGBA(0, 0, 0, 0.95)
		ctx.DrawRoundedRectangle(boxX, boxY, xScale*216, yScale*62, yScale*4)
		ctx.Fill()

		// Text
		ctx.SetRGB255(234, 234, 234)
		ctx.LoadFontFace(font, yScale*17)
		ctx.DrawStringAnchored("Accuracy:", boxX+xScale*3, boxY+yScale*4, 0, 1)
		ctx.DrawStringAnchored("Error: "+strconv.FormatFloat(replay.Early, 'f', 2, 64)+"ms - "+strconv.FormatFloat(replay.Late, 'f', 2, 64)+"ms avg", boxX+xScale*3, boxY+yScale*23, 0, 1)
		if score.Mods&osuapi.ModDoubleTime != 0 {
			replay.UnstableRate *= 1.5
		} else if score.Mods&osuapi.ModHalfTime != 0 {
			replay.UnstableRate *= 0.75
		}
		ctx.DrawStringAnchored("Unstable Rate: "+strconv.FormatFloat(replay.UnstableRate, 'f', 2, 64), boxX+xScale*3, boxY+yScale*43, 0, 1)
	}

	writeImage("./osu-images/back.png", ctx, 0, int(yScale*1080), imgScale*3)

	return ctx.Image(), nil
}

func writeCount(count, font string, ctx *gg.Context, x, y, size float64) {
	ctx.SetRGB255(234, 234, 234)
	ctx.LoadFontFace(font, size)
	w, _ := ctx.MeasureString(count)
	ctx.DrawStringAnchored(count, x, y, 0, 0)
	ctx.LoadFontFace(font, size*0.75)
	ctx.DrawStringAnchored("x", x+w, y, 0, 0)
}

func writeImage(path string, ctx *gg.Context, x, y int, imgScale float64) {
	f, err := os.Open(path)
	if err == nil {
		img, _, err := image.Decode(f)
		if err == nil {
			scale := imgScale / float64(img.Bounds().Dy())
			img = imaging.Fit(img, int(scale*float64(img.Bounds().Dx())), int(scale*float64(img.Bounds().Dy())), imaging.Lanczos)
			ctx.DrawImageAnchored(img, x, y, 0, 1)
		}
	}
}
