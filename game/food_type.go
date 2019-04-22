package game

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

type FoodType struct {
	name string

	// How much, if any, health does this food add?
	health uint8

	// How much, if any, hunger does this food add?
	hunger uint8

	// TODO:
	// Does this particular food leave the egg feeling fuller for longer?
	//saturation float64

	// How many bites does it take to eat this type of food
	// TODO

	image *ebiten.Image
}

var foodCherry = FoodType{
	name: "cherry",

	health: 5,
	hunger: 15,

	image: loadImage(sprites.Food, "cherry"),
}

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
