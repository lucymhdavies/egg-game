package game

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/lucymhdavies/egg-game/resources/sprites"
	// TODO: use some other font than this in future
)

// ButtonStyle determins what kind of button this is
type ButtonStyle struct {
	box    bool
	images struct {
		normal,
		pushed string
	}
	icon struct {
		imageMap *sprites.ImageMap
		image    string
	}
}

var defaultButtonStyle = ButtonStyle{
	box: true,
	images: struct{ normal, pushed string }{
		normal: "ButtonBlueOutline",
		pushed: "ButtonBlueOutlinePushed",
	},
}

// Button is a clickable button which does an action when clicked
type Button struct {
	// Pointer back to parent UI
	parent UIElement

	// Pointer back to the game
	game *Game

	// Where it is on screen
	// TODO: maybe just use an r3.Vector, given we need floats anyway
	position struct {
		// Position on screen
		// This refers to position of button when PUSHED
		// TODO: be consistent with size below...
		X, Y int

		// Z-index. Keep stuff on top of other stuff
		Z int
	}
	size struct {
		// Size of button
		// This refers to size of button when UNPUSHED
		// TODO: be consistant with position above...
		W, H int
	}
	// How much fo the image is border, and how much is interior
	// This refers to padding when buton is PUSHED
	// (when unpushed, add pushDepth to H)
	padding Padding

	visible bool

	images struct {
		// Final images
		normal *ebiten.Image
		pushed *ebiten.Image
	}

	pushed bool

	// have there been any touches on the screen while we've been around?
	// keep track of this separately to whether the buttton is pushed, as if a
	// touch was initiated outside the button, then the button should not be pushed
	haveBeenTouches bool

	// how much shorter is a pushed button compared to a normal button
	pushDepth int

	text      string
	textColor color.RGBA

	// Thing that happens when you push the button
	action WorldFunc
}

// IsMouseOver returns whether or not the mouse cursor is currently
// over the button
func (button *Button) IsMouseOver(x, y int) bool {

	// TODO: Detect whether there's anything in front of the button

	return (x >= button.position.X &&
		x <= button.position.X+button.size.W &&
		y >= button.position.Y &&
		y <= button.position.Y+button.size.H)
}

func (button *Button) Update() error {
	if !button.visible {
		return nil
	}

	// TODO: if UI has handled a click event this frame already
	// return nil
	// Maybe?

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

		// TODO: give it a better pointer back to the world
		button.action(button.game.world)
	}

	//
	// Handle Touches
	//

	touches := ebiten.TouchIDs()

	// if we are not touching the button...
	// and the button is in the pushed state...
	// and we have previously touched the screen...
	if (touches == nil || len(touches) == 0) && button.pushed == true && button.haveBeenTouches {
		// then we must have just released a touch which was over the button
		// TODO: give it a better pointer back to the world
		button.action(button.game.world)
	}
	if touches != nil && len(touches) == 1 {

		// As long as there is precisely one touch
		// i.e. don't try to handle multi-touch for now

		if button.IsMouseOver(ebiten.TouchPosition(touches[0])) {
			// Similar logic to the mouse click
			// As long as the touch started this frame...
			if !button.haveBeenTouches {
				button.pushed = true
			}
		} else {
			button.pushed = false
		}

		// Register that we have initited a touch this frame
		button.haveBeenTouches = true
	}

	//
	// Reset button on no mouse click or touch
	//

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && len(touches) == 0 {
		button.pushed = false
		button.haveBeenTouches = false
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
		buttonY -= float64(button.pushDepth)
	}

	// Move drawing to top corner
	op.GeoM.Translate(button.parent.Position().X, button.parent.Position().Y)
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
	textPos.X += int(button.parent.Position().X)
	textPos.Y += int(button.parent.Position().Y)

	if button.pushed {
		textPos.Y += button.pushDepth
	}

	// TODO: use something else to draw text, but for now this is fine
	text.Draw(screen,
		button.text,
		// TODO: use some other font face
		standardFont,

		textPos.X, textPos.Y,

		button.textColor,
	)

	return nil
}

func (button *Button) IsVisible() bool {
	return button.visible
}
func (button *Button) SetVisible(v bool) {
	button.visible = v
}

func (button *Button) Game() *Game {
	return button.game
}

func (button *Button) Position() r3.Vector {
	return r3.Vector{
		X: float64(button.position.X),
		Y: float64(button.position.Y),
		Z: float64(button.position.Z),
	}
}

func (b *Button) Padding() Padding {
	return b.padding
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewButton(p UIElement, width, height int, style ButtonStyle) *Button {
	b := &Button{
		parent: p,
		action: defaultWorldFunc,
		game:   p.Game(),
	}

	// Default button size, if unspecified
	b.size.W = width
	b.size.H = height

	// Box buttons are simple
	// Create a box, then return
	if style.box {
		b.pushDepth = 4 // TODO: get this atuomatically from the images?
		b.padding = Padding{Left: 5, Right: 5, Top: 5, Bottom: 5}
		b.images.normal = NewBox(width, height, style.images.normal).Image
		b.images.pushed = NewBox(width, height-b.pushDepth, style.images.pushed).Image

		return b
	}

	// Otherwise, this is an icon

	img, _, err := image.Decode(bytes.NewReader(sprites.Icons[style.images.normal]))
	if err != nil {
		log.Fatal(err)
	}
	b.images.normal, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(sprites.Icons[style.images.pushed]))
	if err != nil {
		log.Fatal(err)
	}
	b.images.pushed, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// If a button has an icon, draw it onto the images
	if style.icon.image != "" && style.icon.imageMap != nil {
		op.GeoM.Reset()
		op.ColorM.Reset()
		iconImg := loadImage(*style.icon.imageMap, style.icon.image)
		iconW, iconH := iconImg.Size()

		x := float64(b.size.W-iconW) / 2
		y := float64(b.size.H-iconH) / 2

		op.GeoM.Translate(x, y)

		b.images.normal.DrawImage(iconImg, op)
		b.images.pushed.DrawImage(iconImg, op)
	}

	return b
}
