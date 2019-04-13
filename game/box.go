package game

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
	log "github.com/sirupsen/logrus"
)

// Box is a generic box to be drawn on screen
type Box struct {
	W, H int

	// Final stitched together image
	Image *ebiten.Image
}

func NewBox(width, height int, name string) *Box {
	b := &Box{
		W: width,
		H: height,
	}

	//
	// Load all the images
	//

	images := make(map[string]*ebiten.Image)

	suffixes := []string{
		"topleft", "top", "topright",
		"left", "mid", "right",
		"bottomleft", "bottom", "bottomright",
	}

	for _, suffix := range suffixes {
		img, _, err := image.Decode(bytes.NewReader(sprites.Sprites_png[name+"_"+suffix]))
		if err != nil {
			// TODO: better error handling needed here
			log.Fatal(err)
		}
		images[suffix], _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	}

	// Reset drawing functions
	op.GeoM.Reset()
	op.ColorM.Reset()

	// Create a new empty image
	img, err := ebiten.NewImage(width, height, ebiten.FilterDefault)
	if err != nil {
		// TODO: better error handling needed here
		log.Fatal(err)
	}
	b.Image = img

	//
	// Start with corners
	//

	b.Image.DrawImage(images["topleft"], op)

	sizeX, _ := images["topright"].Size()
	op.GeoM.Translate(float64(b.W-sizeX), 0.0)
	b.Image.DrawImage(images["topright"], op)

	op.GeoM.Reset()
	_, sizeY := images["bottomleft"].Size()
	op.GeoM.Translate(0.0, float64(b.H-sizeY))
	b.Image.DrawImage(images["bottomleft"], op)

	sizeX, _ = images["bottomright"].Size()
	op.GeoM.Translate(float64(b.W-sizeX), 0.0)
	b.Image.DrawImage(images["bottomright"], op)

	//
	// Edges
	//

	// Get sizes of each image
	tlCornerSizeX, tlCornerSizeY := images["topleft"].Size()
	trCornerSizeX, trCornerSizeY := images["topright"].Size()
	topSizeX, topSizeY := images["top"].Size()
	blCornerSizeX, blCornerSizeY := images["bottomleft"].Size()
	brCornerSizeX, brCornerSizeY := images["bottomright"].Size()
	bottomSizeX, bottomSizeY := images["bottom"].Size()
	leftSizeX, leftSizeY := images["left"].Size()
	rightSizeX, rightSizeY := images["right"].Size()
	midSizeX, midSizeY := images["mid"].Size()

	// Calculate stretch amount and draw

	// Top
	op.GeoM.Reset()
	scaleAmountX := (float64(b.W) - float64(tlCornerSizeX+trCornerSizeX)) / float64(topSizeX)
	op.GeoM.Scale(scaleAmountX, 1.0)
	op.GeoM.Translate(float64(tlCornerSizeX), 0.0)
	b.Image.DrawImage(images["top"], op)

	// Bottom
	op.GeoM.Reset()
	scaleAmountX = (float64(b.W) - float64(blCornerSizeX+brCornerSizeX)) / float64(bottomSizeX)
	op.GeoM.Scale(scaleAmountX, 1.0)
	op.GeoM.Translate(float64(blCornerSizeX), float64(b.H-bottomSizeY))
	b.Image.DrawImage(images["bottom"], op)

	// Left
	op.GeoM.Reset()
	scaleAmountY := (float64(b.H) - float64(tlCornerSizeY+blCornerSizeY)) / float64(leftSizeY)
	op.GeoM.Scale(1.0, scaleAmountY)
	op.GeoM.Translate(0.0, float64(tlCornerSizeY))
	b.Image.DrawImage(images["left"], op)

	// Right
	op.GeoM.Reset()
	scaleAmountY = (float64(b.H) - float64(trCornerSizeY+brCornerSizeY)) / float64(rightSizeY)
	op.GeoM.Scale(1.0, scaleAmountY)
	op.GeoM.Translate(float64(b.W-rightSizeX), float64(trCornerSizeY))
	b.Image.DrawImage(images["right"], op)

	//
	// Finally, draw interior
	//

	op.GeoM.Reset()
	scaleAmountX = (float64(b.W) - float64(leftSizeX+rightSizeX)) / float64(midSizeX)
	scaleAmountY = (float64(b.H) - float64(topSizeY+bottomSizeY)) / float64(midSizeY)
	op.GeoM.Scale(scaleAmountX, scaleAmountY)
	op.GeoM.Translate(float64(tlCornerSizeX), float64(tlCornerSizeY))
	b.Image.DrawImage(images["mid"], op)

	return b
}
