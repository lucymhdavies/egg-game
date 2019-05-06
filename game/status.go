package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

// TODO: make this an interface?

// Status is a status effect on the egg
// Sometimes these come with images, which are emotes
type Status struct {
	// Text to show player
	name string

	// When egg has multiple status effects, which emote should we show?
	// lower wins
	priority uint8

	// TODO: maybe encode in the emote what the criteria is for it to show?
	// and maybe encode the effects in here too?

	image *ebiten.Image
}

const (
	StatusSaturated int = iota
	StatusHungry
	StatusStarving
)

var statuses = map[int]Status{
	StatusSaturated: Status{
		name:     "Saturated",
		priority: 5,
		image:    loadImage(sprites.Emotes, "sparkles"),
	},
	StatusHungry: Status{
		name:     "Hungry",
		priority: 2,
		image:    loadImage(sprites.Emotes, "blank"),
	},
	StatusStarving: Status{
		name:     "Starving",
		priority: 1,
		image:    loadImage(sprites.Emotes, "starvation"),
	},

	// TODO: Dying (low health)
	// TODO: Nearly Dead
}

// Custom status effects for the egg being hungry
// TODO: make these a different type that implements Status interface
// i.e. we need some way of tracking what the food is craving, within the Status
var hungerCravings = make(map[FoodType]Status)

func init() {
	// populate hungerCravings based on which food types can be craved
	for _, food := range foodTypes {
		if food.crave {
			hungerCravings[food] = Status{
				name:     "Hungry",
				priority: 2,
				image:    loadImage(sprites.Emotes, fmt.Sprintf("hungry_%s", food.name)),
			}
		}
	}
}
