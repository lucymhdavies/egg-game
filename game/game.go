package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

const (
	ScreenWidth  = 256
	ScreenHeight = 384
)

type Game struct {
	world *World
	input Input
}

func NewGame() *Game {
	return &Game{
		world: NewWorld(256),
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	// TODO: scenemanager stuff, e.g.
	// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/scenemanager.go

	err := g.input.Update()
	if err != nil {
		return err
	}

	err = g.update()
	if err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	err = g.draw(screen)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) update() error {
	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	screen.Fill(color.White)

	g.world.Draw(screen)

	return nil
}
