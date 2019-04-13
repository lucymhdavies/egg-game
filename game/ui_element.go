package game

import "github.com/hajimehoshi/ebiten"

// Generic UI Element interface

type UIElement interface {
	Update() error
	Draw(*ebiten.Image) error
	IsVisible() bool
	SetVisible(bool)
}
