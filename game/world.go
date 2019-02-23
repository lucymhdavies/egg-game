package game

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"

	"github.com/lucymhdavies/egg-game/resources/sprites"
)

var (
	imageGameBG *ebiten.Image
	tileSize    int
)

type World struct {
	Width  int
	Height int

	egg *Egg
}

func NewWorld(size int) *World {
	w := &World{
		Width:  size,
		Height: size,
	}

	// Create egg
	w.egg = NewEgg(w)

	return w
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(sprites.BGTile_png))
	if err != nil {
		log.Fatal(err)
	}
	imageGameBG, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	w, _ := imageGameBG.Size()
	tileSize = w
}

func (world *World) Draw(screen *ebiten.Image) error {
	xNum := world.Width / tileSize
	yNum := world.Height / tileSize

	for i := 0; i < xNum*yNum; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%xNum)*tileSize), float64((i/xNum)*tileSize))

		screen.DrawImage(imageGameBG, op)
	}

	err := world.egg.Draw(screen)
	if err != nil {
		return err
	}

	return nil
}
