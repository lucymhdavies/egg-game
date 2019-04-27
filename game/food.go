package game

import (
	"math/rand"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

type Food struct {
	position r3.Vector // refers to midpoint of egg
	size     r2.Point

	world *World

	// How much of this has the egg consumed?
	// TODO

	foodType FoodType
}

func NewFood(w *World, ft FoodType) *Food {
	// Calculate size based on FoodType
	sizeX, sizeY := ft.image.Size()

	// Random position in world
	xRand := rand.Float64()*
		float64((w.Width)- //      width of world
			(w.Padding().Left+w.Padding().Right)- //      minus world padding
			(sizeX)/2) + // then center on image
		float64(w.Padding().Left) // and add back the world padding
	yRand := rand.Float64()*
		float64((w.Height)- //     height of world
			(w.Padding().Top+w.Padding().Bottom)- //      minus world padding
			(sizeY)/2) + // then center on image
		float64(w.Padding().Top) // and add back the world padding

	f := &Food{
		position: r3.Vector{
			// TODO: random position within world
			xRand,
			yRand,
			0,
		},
		size: r2.Point{
			X: float64(sizeX),
			Y: float64(sizeY),
		},
		world:    w,
		foodType: ft,
	}

	return f
}

func (food *Food) Update() error {
	// TODO: if the food has been in the world too long
	// then it should spoil

	return nil
}

func (food *Food) Draw(screen *ebiten.Image) error {

	op.GeoM.Reset()
	op.ColorM.Reset()

	op.GeoM.Translate(food.position.X-food.size.X/2, food.position.Y-food.size.Y/2)
	screen.DrawImage(food.foodType.image, op)

	return nil
}
