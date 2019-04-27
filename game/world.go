package game

import (
	"bytes"
	"image"
	"log"
	"math"
	"sort"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"

	"github.com/lucymhdavies/egg-game/resources/sprites"
)

var (
	imageGameBG *ebiten.Image
	tileSize    int
	maxFood     = 50
)

type World struct {
	Width  int
	Height int

	// Don't let anything enter the padding area
	padding r2.Point

	// Number of tiles needed to fill world
	xNum int
	yNum int

	egg *Egg

	// items in the world which are not the egg
	// for now, this is just Food
	food []*Food
}

func NewWorld(sizeX, sizeY int) *World {
	xNum := int(math.Ceil(float64(sizeX) / float64(tileSize)))
	yNum := int(math.Ceil(float64(sizeY) / float64(tileSize)))

	w := &World{
		Width:   sizeX,
		Height:  sizeY,
		padding: r2.Point{X: 10.0, Y: 10.0},

		xNum: xNum,
		yNum: yNum,
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

	for _, f := range world.food {
		err := f.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (world *World) Draw(screen *ebiten.Image) error {

	for i := 0; i < world.xNum*world.yNum; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%world.xNum)*tileSize), float64((i/world.xNum)*tileSize))

		screen.DrawImage(imageGameBG, op)
	}

	// TODO: this is where a generic "drawable entities" interface would be useful
	// Would be sorted by implicit and explicit Z-index

	// Draw food that is behind the egg
	for _, f := range world.food {
		if f.position.Y+f.size.Y/2 > world.egg.position.Y+world.egg.size.Y/2 {
			continue
		}

		err := f.Draw(screen)
		if err != nil {
			return err
		}
	}

	err := world.egg.Draw(screen)
	if err != nil {
		return err
	}

	// Draw food that is in front of the egg
	for _, f := range world.food {
		if f.position.Y+f.size.Y/2 <= world.egg.position.Y+world.egg.size.Y/2 {
			continue
		}

		err := f.Draw(screen)
		if err != nil {
			return err
		}
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

func (world *World) AddFood(foodType FoodType) error {

	if len(world.food) >= maxFood {
		return nil
	}

	f := NewFood(world, foodType)
	world.food = append(world.food, f)

	// Sort food by implicit Z-index
	// As food is always on the ground, this is just its Y-coord
	sort.Slice(world.food, func(i, j int) bool {
		return world.food[i].position.Y < world.food[j].position.Y
	})

	return nil
}

// NearestFood takes a point, and returns the nearest food to that
// If there is no food in the world, returns nil
func (world *World) NearestFood(v r3.Vector) *Food {
	if len(world.food) == 0 {
		return nil
	}

	// Initialise with nearest
	var nearestFood *Food
	var minDistance float64

	for i, food := range world.food {
		vecFromEggToFood := food.position.Sub(world.egg.position)
		distance := vecFromEggToFood.Norm()

		// if we're looking at the first food
		// or we've found a new nearest food
		if i == 0 || distance < minDistance {
			minDistance = distance
			nearestFood = food
		}
	}

	return nearestFood
}

// RemoveFood removes a specified Food from the world
func (world *World) RemoveFood(f *Food) {
	var foodIndex int
	for i, food := range world.food {
		if f == food {
			foodIndex = i
			break
		}
	}

	world.food = append(world.food[:foodIndex], world.food[foodIndex+1:]...)
}

// WorldFunc is a generic function which interacts with the world
type WorldFunc func(w *World)

var defaultWorldFunc WorldFunc = func(w *World) {
}
