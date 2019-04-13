package game

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

// Generic UI Element interface

type UIElement interface {
	Update() error
	Draw(*ebiten.Image) error
	IsVisible() bool
	SetVisible(bool)
	Position() r3.Vector
	Game() *Game
}
