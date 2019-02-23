package game

import (
	"errors"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var RegularTermination = errors.New("regular termination")

type Input struct {
}

func (i *Input) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return RegularTermination
	}

	return nil
}
