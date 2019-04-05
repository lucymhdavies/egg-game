package game

import (
	"image/color"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	// TODO: use some other font than this in future
)

type ButtonFunc func(w *World)

var defaultButtonFunc ButtonFunc = func(w *World) {
}

var (
	// Keep track of whether there have ever been touches
	haveBeenTouches bool
)

type Button struct {
	// Pointer back to parent UI
	ui *UI

	// Where it is on screen
	// TODO: maybe just use an r3.Vector, given we need floats anyway
	position struct {
		// Position on screen
		X, Y int

		// Z-index. Keep stuff on top of other stuff
		Z int
	}
	size struct {
		W, H int
	}

	visible bool

	images struct {
		// Final images
		normal *ebiten.Image
		pushed *ebiten.Image
	}

	pushed bool

	text      string
	textColor color.RGBA

	// Thing that happens when you push the button
	action ButtonFunc
}

// IsMouseOver returns whether or not the mouse cursor is currently
// over the button
func (button *Button) IsMouseOver(x, y int) bool {

	return (x >= button.position.X &&
		x <= button.position.X+button.size.W &&
		y >= button.position.Y &&
		y <= button.position.Y+button.size.H)
}

func (button *Button) Update() error {
	if !button.visible {
		return nil
	}

	// TODO: use parent ui as a wrapper around raw ebiten events?
	// Might help with tests

	//
	// Handle Mouse Clicks
	//

	// if mouse button pressed, and cursor is over button...
	// change state to pressed
	// else, state is normal
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if button.IsMouseOver(ebiten.CursorPosition()) {
			// Only press the button if we clicked while the mouse was over
			// don't press the button if we clicked elsewhere, then dragged
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				button.pushed = true
			}
		} else {
			button.pushed = false
		}
	}

	// if mouse button is JUST unpressed, and cursor is still over button...
	// do button action
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
		button.IsMouseOver(ebiten.CursorPosition()) &&
		button.pushed == true {
		button.action(button.ui.game.world)
	}

	//
	// Handle Touches
	//

	touches := ebiten.TouchIDs()

	// if we are not touching the button...
	// and the button is in the pushed state...
	// and we have previously touched the screen...
	if (touches == nil || len(touches) == 0) && button.pushed == true && haveBeenTouches {
		// then we must have just released a touch which was over the button
		button.action(button.ui.game.world)
	}
	if touches != nil && len(touches) == 1 {

		// As long as there is precisely one touch
		// i.e. don't try to handle multi-touch for now

		if button.IsMouseOver(ebiten.TouchPosition(touches[0])) {
			// Similar logic to the mouse click
			// As long as the touch started this frame...
			if !haveBeenTouches {
				button.pushed = true
			}
		} else {
			button.pushed = false
		}

		// Register that we have initited a touch this frame
		haveBeenTouches = true
	}

	//
	// Reset button on no mouse click or touch
	//

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && len(touches) == 0 {
		button.pushed = false
		haveBeenTouches = false
	}

	return nil
}

func (button *Button) Draw(screen *ebiten.Image) error {
	if !button.visible {
		return nil
	}

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Top left corner, which changes depending on whether button is
	// pressed or not
	buttonX := float64(button.position.X)
	buttonY := float64(button.position.Y)

	if !button.pushed {
		// if button not pushed, then move everything up 4 pixels
		// to show the side at the bottom
		buttonY -= 4.0
	}

	// Move drawing to top corner
	op.GeoM.Translate(buttonX, buttonY)

	if !button.pushed {
		screen.DrawImage(button.images.normal, op)
	} else {
		screen.DrawImage(button.images.pushed, op)
	}

	// TODO: center this text in the button
	// need to know how big the text is
	// for now, we know it's width is 5 px per letter
	textWidth := len(button.text) * 5
	// see: ebiten/examples/blocks/blocks/font.go for how to do this better
	textPos := struct{ X, Y int }{
		button.position.X + (button.size.W / 2) - (textWidth / 2),
		button.position.Y + 15,
	}

	if button.pushed {
		textPos.Y += 4
	}

	// TODO: use something else to draw text, but for now this is fine
	text.Draw(screen,
		button.text,
		// TODO: use some other font face
		bitmapfont.Gothic12r,

		textPos.X, textPos.Y,

		button.textColor,
	)

	return nil
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewButton(ui *UI, width, height int) *Button {
	b := &Button{
		ui:     ui,
		action: defaultButtonFunc,
	}

	// Default button size, if unspecified
	b.size.W = width
	b.size.H = height

	// Reset drawing functions
	op.GeoM.Reset()
	op.ColorM.Reset()

	b.images.normal = NewBox(width, height, "ButtonBlueOutline").Image
	b.images.pushed = NewBox(width, height-4, "ButtonBlueOutlinePushed").Image

	return b
}
