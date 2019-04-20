package game

import (
	"bytes"
	"image"
	"log"
	"math"

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

	// Number of tiles needed to fill world
	xNum int
	yNum int

	egg *Egg
}

func NewWorld(sizeX, sizeY int) *World {
	xNum := int(math.Ceil(float64(sizeX) / float64(tileSize)))
	yNum := int(math.Ceil(float64(sizeY) / float64(tileSize)))

	w := &World{
		Width:  sizeX,
		Height: sizeY,
		xNum:   xNum,
		yNum:   yNum,
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

func (world *World) Update() error {
	err := world.egg.Update()
	if err != nil {
		return err
	}

	return nil
}

func (world *World) Draw(screen *ebiten.Image) error {

	for i := 0; i < world.xNum*world.yNum; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%world.xNum)*tileSize), float64((i/world.xNum)*tileSize))

		screen.DrawImage(imageGameBG, op)
	}

	err := world.egg.Draw(screen)
	if err != nil {
		return err
	}

	return nil
}

func (world *World) ReplaceEgg() error {

	if world.egg.state == StateDead {
		world.egg.state = StateRespawning
	} else {
		world.egg = NewEgg(world)
	}

	return nil
}

// WorldFunc is a generic function which interacts with the world
type WorldFunc func(w *World)

var defaultWorldFunc WorldFunc = func(w *World) {
}
