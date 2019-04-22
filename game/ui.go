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
		ui.uiElements["statsWindow"].SetVisible(false)
		ui.uiElements["statsIcon"].SetVisible(false)
		ui.uiElements["respawnButton"].SetVisible(true)
	} else {
		ui.uiElements["statsIcon"].SetVisible(true)
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

	ageLabel := NewLabel(w, "Age", "Age: %d")
	ageLabel.textColor = color.RGBA{0, 0, 0, 255}
	ageLabel.SetVisible(true)
	ageLabel.centered = false
	ageLabel.size.W = w.size.W - 20
	ageLabel.size.H = standardFont.Metrics().Height.Ceil()
	ageLabel.position.X = 10
	ageLabel.position.Y = standardFont.Metrics().Height.Ceil()
	ageLabel.updateFunc = func(w *World) {
		v, _ := w.egg.GetStat("age")
		ageLabel.text = fmt.Sprintf(ageLabel.textFormat, int(v))
	}

	w.uiElements = append(w.uiElements, ageLabel)

	//
	// Health Bar + Label
	//

	healthBar, healthLabel := ui.createLabeledBar(
		w,
		"Health",
		struct{ left, right, top, bottom int }{
			left: 10, right: 10,
			top: ageLabel.position.Y + ageLabel.size.H + 5,
		},
	)
	w.uiElements = append(w.uiElements, healthBar, healthLabel)

	//
	// Hunger Bar + Label
	//

	hungerBar, hungerLabel := ui.createLabeledBar(
		w,
		"Hunger",
		struct{ left, right, top, bottom int }{
			left: 10, right: 10,
			top: healthBar.position.Y + healthBar.size.H + 5,
		},
	)
	w.uiElements = append(w.uiElements, hungerBar, hungerLabel)

	return w
}

// TODO: this could be part of Layout?
// Create a generic bar with label
// TODO: functional options, you know, as I need to do everywhere else...
// TODO: allow any numeric type for value
// TODO: allow any parent (not just window)
//   possible solutions:
//   - give UIElement a Size() function
func (ui UI) createLabeledBar(w *Window, name string,
	offset struct{ left, right, top, bottom int }) (*Bar, *Label) {

	paddingX := (offset.left + offset.right)

	label := NewLabel(w, name, name+": %d")
	label.textColor = color.RGBA{0, 0, 0, 255}
	label.SetVisible(true)
	label.centered = false
	label.size.W = w.size.W - paddingX
	label.size.H = standardFont.Metrics().Height.Ceil()
	label.position.X = offset.left
	label.position.Y = offset.top

	label.updateFunc = func(w *World) {
		v, _ := w.egg.GetStat(name)
		label.text = fmt.Sprintf(label.textFormat, int(v))
	}

	bar := NewBar(w, w.size.W-paddingX, 18, "green")
	bar.SetVisible(true)
	bar.position.X = offset.left
	bar.position.Y = label.position.Y + label.size.H + 5
	bar.max = 255.0

	bar.updateFunc = func(w *World) {
		bar.value, _ = w.egg.GetStat(name)
	}

	return bar, label
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
