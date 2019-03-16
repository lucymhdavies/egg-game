package game

import (
	"errors"
	"runtime"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var RegularTermination = errors.New("regular termination")

type Input struct {
	// Pointer back to parent Game
	game *Game
}

func (i *Input) Update() error {
	if runtime.GOARCH != "js" {
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return RegularTermination
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			i.game.world.ReplaceEgg()
		}
	}

	return nil
}
