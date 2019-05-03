package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

// TODO: rename from Emote to Status?

// Emote refers to a bubble above an egg displaying status effects
type Emote struct {
	// Text to show player
	name string

	// When egg has multiple status effects, which emote should we show?
	// lower wins
	priority uint8

	// TODO: maybe encode in the emote what the criteria is for it to show?
	// and maybe encode the effects in here too?

	image *ebiten.Image
}

var emotes = map[string]Emote{
	"saturated": Emote{
		name:     "Saturated",
		priority: 1,
		image:    loadImage(sprites.Emotes, "sparkles"),
	},

	// TODO: Hungry
	// TODO: Dying
	// TODO: Nearly Dead
}
