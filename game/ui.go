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
	ui.uiElements["statsIcon"] = ui.createStatsIcon()

	// TODO: sort by Z-index, showing lower elements first
	// for now, just do this manually
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["statsIcon"])
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["respawnButton"])
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["statsWindow"])

	return ui
}

// TODO: move this to a different file?
func (ui *UI) createRespawnButton() *Button {
	b := NewButton(ui, 100, 34, defaultButtonStyle)

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

	//
	// The Window itself
	//

	// Add a 36px buffer at the bottom, so as not to overlap the icons
	// on the bottom row of the screen
	w := NewWindow(ui, ScreenWidth-20, ScreenHeight-20-36)
	w.position.X = 10
	w.position.Y = 10
	//w.SetVisible(true)

	// Z-Index
	w.position.Z = 20

	w.text = "Stats"
	w.textColor = color.RGBA{0, 0, 0, 255}

	//
	// Label to display Age
	//

	ageLabel := NewLabel(w, "Age", "Age: %v")
	ageLabel.textColor = color.RGBA{0, 0, 0, 255}
	ageLabel.SetVisible(true)
	ageLabel.centered = false
	ageLabel.size.W = w.size.W - 20
	ageLabel.size.H = standardFont.Metrics().Height.Ceil()
	ageLabel.position.X = 10
	ageLabel.position.Y = standardFont.Metrics().Height.Ceil()
	ageLabel.updateFunc = func(w *World) {
		ageLabel.text = fmt.Sprintf(ageLabel.textFormat, int(w.egg.stats.age))
	}

	w.uiElements = append(w.uiElements, ageLabel)

	//
	// Label to display Health
	//

	healthLabel := NewLabel(w, "Health", "Health: %v")
	healthLabel.textColor = color.RGBA{0, 0, 0, 255}
	healthLabel.SetVisible(true)
	healthLabel.centered = false
	healthLabel.size.W = w.size.W - 20
	healthLabel.size.H = standardFont.Metrics().Height.Ceil()
	healthLabel.position.X = 10
	healthLabel.position.Y = ageLabel.position.Y + ageLabel.size.H + 5
	healthLabel.updateFunc = func(w *World) {
		healthLabel.text = fmt.Sprintf(healthLabel.textFormat, int(w.egg.stats.health))
	}

	w.uiElements = append(w.uiElements, healthLabel)

	//
	// Bar to display Health
	//

	healthBar := NewBar(w, w.size.W-20, 18, "green")
	healthBar.SetVisible(true)
	healthBar.position.X = 10
	healthBar.position.Y = healthLabel.position.Y + healthLabel.size.H + 5
	healthBar.max = 255.0
	healthBar.updateFunc = func(w *World) {
		healthBar.value = float64(w.egg.stats.health)
	}

	w.uiElements = append(w.uiElements, healthBar)

	//
	// Close Button
	//

	// 32x32, for icon, then 10x14, for button border
	hideStatsButton := NewButton(w, 36, 36,
		ButtonStyle{
			box: false,
			images: struct{ normal, pushed string }{
				normal: "red_boxCross",
				pushed: "grey_box",
			},
		},
	)
	hideStatsButton.SetVisible(true)

	// Bottom right of window (5px padding)
	hideStatsButton.position.X = w.size.W - hideStatsButton.size.W - 5
	hideStatsButton.position.Y = w.size.H - hideStatsButton.size.H - 5
	hideStatsButton.action = func(world *World) {
		w.SetVisible(false)
	}

	w.uiElements = append(w.uiElements, hideStatsButton)

	return w
}

// TODO: move this to a different file?
func (ui *UI) createStatsIcon() *Button {
	b := NewButton(ui, 36, 36,
		ButtonStyle{
			box: false,
			images: struct{ normal, pushed string }{
				normal: "transparentDark_star", // TODO: dedicated info icon
				pushed: "transparentLight_star",
			},
		},
	)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = ScreenWidth - b.size.W - 10
	b.position.Y = ScreenHeight - b.size.H - 5
	b.visible = true

	// Z-Index
	b.position.Z = 10

	b.action = func(world *World) {
		window := ui.uiElements["statsWindow"]
		window.SetVisible(!window.IsVisible())
	}
	// TODO: while statsWindow is visible, then this button should display as pushed?
	// Means it's gonna need an UpdateFunc as well as an Action Func
	// Or maybe some third Highlighted state?

	return b
}
