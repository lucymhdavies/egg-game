package game

import (
	"bytes"
	"image"
	"log"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

var (
	op = &ebiten.DrawImageOptions{}
)

type Egg struct {
	position r2.Point
	size     r2.Point
	images   struct {
		body   *ebiten.Image
		eyes   *ebiten.Image
		shadow *ebiten.Image
		emote  *ebiten.Image
	}
	name string
	// TODO: state
	// idle, bounce, sleep
}

func NewEgg(w *World) *Egg {
	e := &Egg{
		position: r2.Point{float64(w.Width) / 2, float64(w.Height) / 2},
		name:     "Egg",
	}

	img, _, err := image.Decode(bytes.NewReader(sprites.EggBody_png))
	if err != nil {
		log.Fatal(err)
	}
	bodyImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	sizeX, sizeY := bodyImg.Size()
	e.size = r2.Point{float64(sizeX), float64(sizeY)}

	e.images.body = bodyImg
	return e
}

func (egg *Egg) Draw(screen *ebiten.Image) error {
	op.GeoM.Reset()
	op.ColorM.Reset()

	// Move it to its position
	op.GeoM.Translate(egg.position.X-egg.size.X/2, egg.position.Y-egg.size.Y/2)

	// Draw it
	screen.DrawImage(egg.images.body, op)
	return nil
}
