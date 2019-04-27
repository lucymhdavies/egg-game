package game

import (
	"fmt"
	"image/color"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

type UI struct {
	// Pointer back to parent Game
	game *Game

	// All UI elements (by name)
	uiElements map[string]UIElement
	// All UI elements (sorted by Z Index)
	zSortedUIElements []UIElement

	padding Padding

	// TODO: need some way of referring to SPECIFIC buttons from other packages?
	// Or refer to game state within UI
	// e.g. when egg is dead, show respawn button
}

func (ui *UI) Update() error {

	// TODO: do this elsewhere? some other way?
	if ui.game.world.egg.state == StateDead {
		// Hide all other UI elements
		// Maybe just loop through them all and set visible false?
		ui.uiElements["statsWindow"].SetVisible(false)
		ui.uiElements["itemsWindow"].SetVisible(false)
		ui.uiElements["statsIcon"].SetVisible(false)
		ui.uiElements["itemsIcon"].SetVisible(false)

		// SHow the respawn button
		ui.uiElements["respawnButton"].SetVisible(true)
	} else {
		// Show on-screen buttons
		// TODO: showAllButtons some other way?
		ui.uiElements["statsIcon"].SetVisible(true)
		ui.uiElements["itemsIcon"].SetVisible(true)

		// hide respawn button
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
func (ui *UI) Padding() Padding {
	return ui.padding
}

func (ui *UI) Game() *Game {
	return ui.game
}

func NewUI(g *Game) *UI {
	ui := &UI{
		game:       g,
		uiElements: make(map[string]UIElement),
		padding:    g.world.Padding(),
	}

	ui.uiElements["respawnButton"] = ui.createRespawnButton()
	ui.uiElements["statsIcon"] = ui.createStatsIcon()
	ui.uiElements["itemsIcon"] = ui.createItemsIcon()
	ui.uiElements["statsWindow"] = ui.createStatsWindow()
	ui.uiElements["itemsWindow"] = ui.createItemsWindow()

	// TODO: sort by Z-index, showing lower elements first
	// for now, just do this manually

	// on-screen buttons
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["statsIcon"])
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["itemsIcon"])

	// Respawn button should never be on screen at the same time as anything else
	// but put it here anyway
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["respawnButton"])

	// All windows, just add last
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["statsWindow"])
	ui.zSortedUIElements = append(ui.zSortedUIElements, ui.uiElements["itemsWindow"])

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

	// Add a 41px buffer at the bottom, so as not to overlap the icons
	// on the bottom row of the screen
	// 41 = height of button (36px) + 5px padding
	paddingX := ui.Padding().Left + ui.Padding().Right
	paddingY := ui.Padding().Top + ui.Padding().Bottom + 41

	w := NewWindow(ui, ScreenWidth-paddingX, ScreenHeight-paddingY)
	w.position.X = ui.Padding().Left
	w.position.Y = ui.Padding().Top

	// Z-Index
	w.position.Z = 20

	//
	// Label for Window Title
	//

	titleLabel := NewLabel(w, "Stats", "Stats")
	titleLabel.textColor = color.RGBA{0, 0, 0, 255}
	titleLabel.SetVisible(true)
	titleLabel.centered = true
	titleLabel.size.W = w.size.W - w.Padding().Left - w.Padding().Right
	titleLabel.size.H = standardFont.Metrics().Height.Ceil()
	titleLabel.position.X = w.Padding().Left
	titleLabel.position.Y = w.Padding().Top -
		(standardFont.Metrics().Height.Ceil() - standardFont.Metrics().Ascent.Ceil())

	w.uiElements = append(w.uiElements, titleLabel)

	//
	// Label to display Age
	//

	ageLabel := NewLabel(w, "Age", "Age: %d")
	ageLabel.textColor = color.RGBA{0, 0, 0, 255}
	ageLabel.SetVisible(true)
	ageLabel.centered = false
	ageLabel.size.W = w.size.W - w.Padding().Left - w.Padding().Right
	ageLabel.size.H = standardFont.Metrics().Height.Ceil()
	ageLabel.position.X = w.Padding().Left
	ageLabel.position.Y = titleLabel.position.Y + titleLabel.size.H + 5
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
		Padding{
			Left: w.Padding().Left, Right: w.Padding().Right,
			Top:    ageLabel.position.Y + ageLabel.size.H + 5,
			Bottom: 0,
		},
	)
	w.uiElements = append(w.uiElements, healthBar, healthLabel)

	//
	// Hunger Bar + Label
	//

	hungerBar, hungerLabel := ui.createLabeledBar(
		w,
		"Hunger",
		Padding{
			Left: w.Padding().Left, Right: w.Padding().Right,
			Top:    healthBar.position.Y + healthBar.size.H + 5,
			Bottom: 0,
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
	offset Padding) (*Bar, *Label) {

	paddingX := (offset.Left + offset.Right)

	label := NewLabel(w, name, name+": %d")
	label.textColor = color.RGBA{0, 0, 0, 255}
	label.SetVisible(true)
	label.centered = false
	label.size.W = w.size.W - paddingX
	label.size.H = standardFont.Metrics().Height.Ceil()
	label.position.X = offset.Left
	label.position.Y = offset.Top

	label.updateFunc = func(w *World) {
		v, _ := w.egg.GetStat(name)
		label.text = fmt.Sprintf(label.textFormat, int(v))
	}

	bar := NewBar(w, w.size.W-paddingX, 18, "green")
	bar.SetVisible(true)
	bar.position.X = offset.Left
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
				normal: "transparentDark_bars",
				pushed: "transparentLight_bars",
			},
		},
	)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = ScreenWidth - b.size.W - ui.Padding().Right
	b.position.Y = ScreenHeight - b.size.H - ui.Padding().Bottom
	b.visible = true

	b.action = func(world *World) {
		ui.uiElements["itemsWindow"].SetVisible(false)

		window := ui.uiElements["statsWindow"]
		window.SetVisible(!window.IsVisible())
	}
	// TODO: while statsWindow is visible, then this button should display as pushed?
	// Means it's gonna need an UpdateFunc as well as an Action Func
	// Or maybe some third Highlighted state?

	return b
}

func (ui *UI) createItemsWindow() *Window {

	//
	// The Window itself
	//

	// Add a 41px buffer at the bottom, so as not to overlap the icons
	// on the bottom row of the screen
	// 41 = height of button (36px) + 5px padding
	paddingX := ui.Padding().Left + ui.Padding().Right
	paddingY := ui.Padding().Top + ui.Padding().Bottom + 41

	w := NewWindow(ui, ScreenWidth-paddingX, ScreenHeight-paddingY)
	w.position.X = ui.Padding().Left
	w.position.Y = ui.Padding().Top

	// Z-Index
	w.position.Z = 20

	//
	// Label for Window Title
	//

	titleLabel := NewLabel(w, "Items", "Items")
	titleLabel.textColor = color.RGBA{0, 0, 0, 255}
	titleLabel.SetVisible(true)
	titleLabel.centered = true
	titleLabel.size.W = w.size.W - w.Padding().Left - w.Padding().Right
	titleLabel.size.H = standardFont.Metrics().Height.Ceil()
	titleLabel.position.X = w.Padding().Left
	titleLabel.position.Y = w.Padding().Top -
		(standardFont.Metrics().Height.Ceil() - standardFont.Metrics().Ascent.Ceil())

	w.uiElements = append(w.uiElements, titleLabel)

	//
	// Add a temporary button for testing
	// This will be part of a Grid later
	//

	initialOffset := Padding{
		Top:  titleLabel.position.Y + titleLabel.size.H + 5,
		Left: w.Padding().Left,
	}
	for foodType := range foodTypes {
		w.uiElements = append(w.uiElements, ui.createItemsWindowFoodIcon(w,
			Padding{
				Top:  initialOffset.Top,
				Left: initialOffset.Left,
			},
			foodType,
		))

		// TODO: implement Grid layout
		// 36 = width of icon
		// 5 for padding
		initialOffset.Left += 36 + 5

		// TODO: if this means we're going to spill off the right
		// set to parent.Padding().Left, then incremement initialOffset.Top
	}

	return w
}

// TODO: move this to a different file?
func (ui *UI) createItemsIcon() *Button {
	b := NewButton(ui, 36, 36,
		ButtonStyle{
			box: false,
			images: struct{ normal, pushed string }{
				normal: "transparentDark_star",
				pushed: "transparentLight_star",
			},
		},
	)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = ui.Padding().Left
	b.position.Y = ScreenHeight - b.size.H - ui.Padding().Top
	b.visible = true

	b.action = func(world *World) {
		ui.uiElements["statsWindow"].SetVisible(false)

		window := ui.uiElements["itemsWindow"]
		window.SetVisible(!window.IsVisible())
	}
	// TODO: while statsWindow is visible, then this button should display as pushed?
	// Means it's gonna need an UpdateFunc as well as an Action Func
	// Or maybe some third Highlighted state?

	return b
}

func (ui *UI) createItemsWindowFoodIcon(parent UIElement, offset Padding, foodType string) *Button {
	b := NewButton(parent, 36, 36,
		ButtonStyle{
			box: false,
			images: struct{ normal, pushed string }{
				normal: "transparentDark_blank",
				pushed: "transparentLight_blank",
			},
			icon: struct {
				imageMap *sprites.ImageMap
				image    string
			}{
				imageMap: &sprites.Food,
				image:    foodType,
			},
		},
	)

	// Center button horizontally, and stick at bottom of screen
	b.position.X = offset.Left
	b.position.Y = offset.Top
	b.visible = true

	b.action = func(world *World) {
		ui.uiElements["itemsWindow"].SetVisible(false)
		ui.game.world.AddFood(foodTypes[foodType])
	}

	return b
}
