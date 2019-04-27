package game

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

// Generic UI Element interface

type Padding struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

type UIElement interface {
	Update() error
	Draw(*ebiten.Image) error
	IsVisible() bool
	SetVisible(bool)
	Position() r3.Vector // TODO: this should just be a struct of ints
	// TODO: Size()
	// TOOD: Parent()
	Padding() Padding
	Game() *Game
}
