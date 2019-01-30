package timeline

// https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go
// https://developers.google.com/web/tools/chrome-devtools/network-performance/understanding-resource-timing?hl=ko

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Item struct {
	Start float64
	End   float64
	URL   string
	Type  string
}

func generateColorClass(category map[string]int) map[string]*color.RGBA {
	cmap := make(map[string]*color.RGBA)

	for k, _ := range category {
		r := rand.Intn(256)
		g := rand.Intn(256)
		b := rand.Intn(256)
		cmap[k] = &color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}

	return cmap
}

func DrawTimeline(max float64, domComplete float64, loadEvent float64, category map[string]int, data []*Item, output string) {
	rand.Seed(time.Now().UnixNano())

	colorMap := generateColorClass(category)

	SCALE := 10.0
	MAX := len(data)
	LEFT_TEXT := 5
	LEFT_MARGIN := 30
	// RIGHT_MARGIN := 10
	MARGIN := 20
	GAP := 2
	BAR_H := 13
	MAX_X := LEFT_MARGIN + int(max/SCALE) + 300
	MAX_Y := MAX*(BAR_H+GAP) + MARGIN*2

	target := image.NewRGBA(image.Rect(0, 0, MAX_X, MAX_Y))
	draw.Draw(target, target.Bounds(), image.White, image.ZP, draw.Src)

	var y int = 5
	for i := 0; i < MAX; i++ {
		d := data[i]

		x := int(d.Start / SCALE)
		w := int(d.End / SCALE)
		// rand.Intn(MAX_X - LEFT_MARGIN - RIGHT_MARGIN)

		drawText(target, LEFT_TEXT, y+MARGIN+10, fmt.Sprintf("%03d", i))
		draw.Draw(target, image.Rect(LEFT_MARGIN+x, y+MARGIN, LEFT_MARGIN+w, y+MARGIN+BAR_H), &image.Uniform{colorMap[d.Type]}, image.ZP, draw.Src)
		if len(d.URL) > 80 {
			d.URL = d.URL[:80]
		}
		drawText(target, 6+LEFT_MARGIN+w, y+MARGIN+10, d.URL)

		y += BAR_H + GAP
	}

	// draw.Draw(target, image.Rect(int(domComplete/SCALE), 0, 1+int(domComplete/SCALE), MAX_Y), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.ZP, draw.Src)
	draw.Draw(target, image.Rect(LEFT_MARGIN+int(loadEvent/SCALE), 0, LEFT_MARGIN+1+int(loadEvent/SCALE), MAX_Y), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.ZP, draw.Src)
	drawText(target, LEFT_MARGIN+int(loadEvent/SCALE)+3, 15, fmt.Sprintf("onLoad: %.2f sec", loadEvent/1000.0))

	draw.Draw(target, image.Rect(LEFT_MARGIN+int(max/SCALE), 0, LEFT_MARGIN+1+int(max/SCALE), MAX_Y), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.ZP, draw.Src)
	drawText(target, LEFT_MARGIN+int(max/SCALE)+3, 15, fmt.Sprintf("finished: %.2f sec", max/1000.0))

	// color agenda
	y = 25
	for n, c := range colorMap {
		draw.Draw(target, image.Rect(MAX_X-130, y, MAX_X-110, y+BAR_H), &image.Uniform{c}, image.ZP, draw.Src)
		drawText(target, MAX_X-105, y+10, n)
		y += BAR_H + GAP
	}

	fo, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	err = png.Encode(fo, target)
	if err != nil {
		log.Panic(err)
	}
}

func drawText(img *image.RGBA, x int, y int, text string) {
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)},
	}

	dr.DrawString(text)
}
