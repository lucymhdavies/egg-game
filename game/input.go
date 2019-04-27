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

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		// Hide other windows
		i.game.ui.uiElements["itemsWindow"].SetVisible(false)

		// Toggle this window
		v := i.game.ui.uiElements["statsWindow"].IsVisible()
		i.game.ui.uiElements["statsWindow"].SetVisible(!v)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		// Hide other windows
		i.game.ui.uiElements["statsWindow"].SetVisible(false)

		// Toggle this window
		v := i.game.ui.uiElements["itemsWindow"].IsVisible()
		i.game.ui.uiElements["itemsWindow"].SetVisible(!v)
	}

	return nil
}
