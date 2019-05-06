package game

import (
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

	// Can the egg crave this food?
	crave bool

	image  *ebiten.Image
	shadow *ebiten.Image
}

/* TODO
const (
	FoodCherry int = iota
	FoodPineapple
	etc.
)
*/

var foodTypes = map[string]FoodType{
	"cherry": FoodType{
		name: "cherry",

		health:     5,
		hunger:     15,
		bites:      1,
		saturation: true,
		crave:      true,

		image:  loadImage(sprites.Food, "cherry"),
		shadow: createShadowImage(sprites.Food, "cherry"),
	},
	"pineapple": FoodType{
		name: "pineapple",

		health:     2,
		hunger:     20,
		bites:      5,
		saturation: true,
		crave:      true,

		image:  loadImage(sprites.Food, "pineapple"),
		shadow: createShadowImage(sprites.Food, "pineapple"),
	},
	"egg": FoodType{
		// Yep, that's kinda creepy. eggs eating eggs
		name: "egg",

		health: 10,
		hunger: 0,
		bites:  5,

		image:  loadImage(sprites.Food, "egg"),
		shadow: createShadowImage(sprites.Food, "egg"),
	},
	"donut": FoodType{
		name: "donut",

		health: 0,
		hunger: 10,
		bites:  2,
		crave:  true,

		image:  loadImage(sprites.Food, "donut"),
		shadow: createShadowImage(sprites.Food, "donut"),
	},
	"definitely_not_beer": FoodType{
		name: "definitely_not_beer",

		health: -5,
		hunger: 0,
		bites:  2,
		// TODO: add to fun, when we implement that

		image:  loadImage(sprites.Food, "definitely_not_beer"),
		shadow: createShadowImage(sprites.Food, "definitely_not_beer"),
	},
}
