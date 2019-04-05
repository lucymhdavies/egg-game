package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type UI struct {
	// Pointer back to parent Game
	game *Game

	// All UI elements
	// For now, we only have buttons, so have this as a slice of buttons
	// TODO: In future, have a UIElements interface
	uiElements []*Button

	// TODO: need some way of referring to SPECIFIC buttons from other packages?
	// Or refer to game state within UI
	// e.g. when egg is dead, show respawn button
}

func (ui *UI) Update() error {

	// For testing...
	if ui.game.world.egg.state == StateDead {
		ui.uiElements[0].visible = true
	}

	// TODO: sort by Z-index, updating higher elements first

	// TODO: prevent overlapping buttons from doing both their actions
	// e.g. have a channel with capacity 1, and buttons push actions onto that?
	// maybe? i dunno.

	for _, e := range ui.uiElements {
		e.Update()
	}

	return nil
}

func (ui *UI) Draw(screen *ebiten.Image) error {

	// TODO: sort by Z-index, showing lower elements first

	for _, e := range ui.uiElements {
		e.Draw(screen)
	}
	return nil
}

func NewUI(g *Game) *UI {
	ui := &UI{
		game: g,
	}

	// For testing, draw an example button
	b := NewButton(ui, 100, 34)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = ScreenWidth/2 - b.size.W/2
	b.position.Y = ScreenHeight - b.size.H - 5
	b.text = "Respawn"
	b.textColor = color.RGBA{0, 0, 0, 255}
	b.action = func(w *World) {
		w.ReplaceEgg()
		b.visible = false
	}

	ui.uiElements = append(ui.uiElements, b)

	return ui
}
