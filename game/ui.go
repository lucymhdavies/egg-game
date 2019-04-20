package game

import (
	"fmt"
	"image/color"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

type UI struct {
	// Pointer back to parent Game
	game *Game

	// All UI elements (by name)
	uiElements map[string]UIElement
	// All UI elements (sorted by Z Index)
	zSortedUIElements []UIElement

	// TODO: need some way of referring to SPECIFIC buttons from other packages?
	// Or refer to game state within UI
	// e.g. when egg is dead, show respawn button
}

func (ui *UI) Update() error {

	// For testing...
	if ui.game.world.egg.state == StateDead {
		ui.uiElements["respawnButton"].SetVisible(true)
	} else {
		ui.uiElements["respawnButton"].SetVisible(false)
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

	for _, e := range ui.zSortedUIElements {
		e.Draw(screen)
	}
	return nil
}

func (ui *UI) IsVisible() bool {
	// TODO: do we even want to be able to show/hide the entire UI?
	return true
}

func (ui *UI) SetVisible(v bool) {
	// needs to exist to implement UIElement
	// but we can't hide the root UI, so do nothing
}
func (ui *UI) Position() r3.Vector {
	return r3.Vector{X: 0.0, Y: 0.0, Z: 0.0}
}

func (ui *UI) Game() *Game {
	return ui.game
}

func NewUI(g *Game) *UI {
	ui := &UI{
		game:       g,
		uiElements: make(map[string]UIElement),
	}

	ui.uiElements["respawnButton"] = ui.createRespawnButton()
	ui.uiElements["statsWindow"] = ui.createStatsWindow()

	// TODO: sort by Z-index, showing lower elements first
	// for now, just do this manually
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["respawnButton"])
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["statsWindow"])

	return ui
}

// TODO: move this to a different file?
func (ui *UI) createRespawnButton() *Button {
	// For testing, draw an example button
	b := NewButton(ui, 100, 34)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = ScreenWidth/2 - b.size.W/2
	b.position.Y = ScreenHeight - b.size.H - 5

	// Z-Index
	b.position.Z = 10

	b.text = "Respawn"
	b.textColor = color.RGBA{0, 0, 0, 255}
	b.action = func(w *World) {
		w.ReplaceEgg()
		b.visible = false
	}

	return b
}

func (ui *UI) createStatsWindow() *Window {

	w := NewWindow(ui, ScreenWidth-20, ScreenHeight-20)
	w.position.X = 10
	w.position.Y = 10

	// Z-Index
	w.position.Z = 20

	w.text = "Stats"
	w.textColor = color.RGBA{0, 0, 0, 255}

	ageLabel := NewLabel(w, "Age:", "Age: %v")
	ageLabel.textColor = color.RGBA{0, 0, 0, 255}
	ageLabel.SetVisible(true)
	ageLabel.centered = false
	ageLabel.position.X = 10
	ageLabel.position.Y = standardFont.Metrics().Height.Ceil()
	ageLabel.updateFunc = func(w *World) {
		ageLabel.text = fmt.Sprintf(ageLabel.textFormat, int(w.egg.stats.age))
	}

	w.uiElements = append(w.uiElements, ageLabel)

	return w
}
