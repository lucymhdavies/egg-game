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

	// Load image from ImageMap
	img, _, err := image.Decode(bytes.NewReader(imageMap[name]))
	if err != nil {
		// TODO: better error handling needed here
		log.Fatalf("Error loading %s - %v", name, err)
	}

	// Create an ebiten Image from the go Image
	i, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		// TODO: better error handling needed here
		log.Fatalf("Error converting %s - %v", name, err)
	}
	return i
}

func createShadowImage(imageMap sprites.ImageMap, name string) *ebiten.Image {

	// Load the image we want to create a shadow from
	shadowSourceImage := loadImage(imageMap, name)
	sizeX, sizeY := shadowSourceImage.Size()

	// Create a new empty image for the shadow
	shadowImage, _ := ebiten.NewImage(sizeX, sizeY/4, ebiten.FilterDefault)

	// Squish and set color
	op.GeoM.Reset()
	op.ColorM.Reset()
	op.GeoM.Scale(1, 0.25)
	op.ColorM.Scale(0, 0, 0, 0.5)

	shadowImage.DrawImage(shadowSourceImage, op)

	return shadowImage
}

/*
	// Shadow = body, but squashed vertically
	shadowImgRaw, _ := ebiten.NewImage(sizeX, sizeY/4, ebiten.FilterDefault)
	shadowImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault) // TODO: just use bodyImg
	op.GeoM.Scale(1, 0.25)
	op.ColorM.Scale(0, 0, 0, 0.5)
	shadowImgRaw.DrawImage(shadowImg, op)
	e.images.shadow = shadowImgRaw
*/
