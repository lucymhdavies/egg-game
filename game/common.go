package game

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

// TODO: put generic sprite interface here?

var (
	op = &ebiten.DrawImageOptions{}

	standardFont = bitmapfont.Gothic12r
)

// TODO: make this helper function a bit better
func loadImage(imageMap sprites.ImageMap, name string) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(imageMap[name]))
	if err != nil {
		// TODO: better error handling needed here
		log.Fatal(err)
	}
	i, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		// TODO: better error handling needed here
		log.Fatal(err)
	}
	return i
}
