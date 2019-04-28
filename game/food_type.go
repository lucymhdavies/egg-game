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

	// Does this food give saturation?
	// i.e. if the egg fills up while eating this, will it take a while before
	// it's hungry again?
	saturation bool

	image *ebiten.Image
}

var foodTypes = map[string]FoodType{
	"cherry": FoodType{
		name: "cherry",

		health:     5,
		hunger:     15,
		bites:      1,
		saturation: true,

		image: loadImage(sprites.Food, "cherry"),
	},
	"pineapple": FoodType{
		name: "pineapple",

		health:     2,
		hunger:     20,
		bites:      5,
		saturation: true,

		image: loadImage(sprites.Food, "pineapple"),
	},
	"egg": FoodType{
		// Yep, that's kinda creepy. eggs eating eggs
		name: "egg",

		health: 10,
		hunger: 0,
		bites:  5,

		image: loadImage(sprites.Food, "egg"),
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
