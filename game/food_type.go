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

	// How many bites does it take to eat this type of food
	bites uint8

	// How much, if any, health does this food add (per bite)?
	health int8

	// How much, if any, hunger does this food add (per bite)?
	hunger int8

	// TODO:
	// Does this particular food leave the egg feeling fuller for longer?
	//saturation float64

	image *ebiten.Image
}

var foodTypes = map[string]FoodType{
	"cherry": FoodType{
		name: "cherry",

		health: 5,
		hunger: 15,
		bites:  1,

		image: loadImage(sprites.Food, "cherry"),
	},
	"pineapple": FoodType{
		name: "pineapple",

		health: 2,
		hunger: 5,
		bites:  5,

		image: loadImage(sprites.Food, "pineapple"),
	},
	"donut": FoodType{
		name: "donut",

		health: 0,
		hunger: 10,
		bites:  2,

		image: loadImage(sprites.Food, "donut"),
	},
	"definitely_not_beer": FoodType{
		name: "definitely_not_beer",

		health: -5,
		hunger: 0,
		bites:  2,
		// TODO: add to fun, when we implement that

		image: loadImage(sprites.Food, "definitely_not_beer"),
	},
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
