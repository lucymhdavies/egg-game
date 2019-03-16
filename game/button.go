package game

import (
	"bytes"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/lucymhdavies/egg-game/resources/sprites"
	log "github.com/sirupsen/logrus"

	// TODO: use some other font than this in future
	"github.com/hajimehoshi/bitmapfont"
)

type ButtonFunc func(w *World)

var defaultButtonFunc ButtonFunc = func(w *World) {
}

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
		// TODO: have buttons made of sub-images, i.e. corners, edges, middles, etc.
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
func (button *Button) IsMouseOver() bool {
	x, y := ebiten.CursorPosition()

	return (x >= button.position.X &&
		x <= button.position.X+button.size.W &&
		y >= button.position.Y &&
		y <= button.position.Y+button.size.H)
}

func (button *Button) Update() error {
	if !button.visible {
		return nil
	}

	// TODO: detect click / tap, etc, then call arbitrary function

	// if mouse button pressed, and cursor is over button...
	// change state to pressed
	// else, state is normal
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && button.IsMouseOver() {
		button.pushed = true
	} else {
		button.pushed = false
	}

	// if mouse button is JUST unpressed, and cursor is over button...
	// do button action
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && button.IsMouseOver() {
		button.action(button.ui.game.world)
	}

	return nil
}

func (button *Button) Draw(screen *ebiten.Image) error {
	if !button.visible {
		return nil
	}

	op.GeoM.Reset()
	op.ColorM.Reset()

	op.GeoM.Translate(float64(button.position.X), float64(button.position.Y))

	if !button.pushed {
		screen.DrawImage(button.images.normal, op)
	} else {
		screen.DrawImage(button.images.pushed, op)
	}

	// TODO: center this text in the button
	textPos := struct{ X, Y int }{
		button.position.X + 20,
		button.position.Y + 20,
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
func NewButton(ui *UI) *Button {
	b := &Button{
		ui:     ui,
		action: defaultButtonFunc,
	}

	// TODO: something to specify which image to use
	// for now, just use the one we have

	img, _, err := image.Decode(bytes.NewReader(sprites.ButtonBlueOutline_png))
	if err != nil {
		log.Fatal(err)
	}
	normalImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	b.images.normal = normalImg

	// Working on the assumption that button images are the same size...
	b.size.W, b.size.H = normalImg.Size()

	img, _, err = image.Decode(bytes.NewReader(sprites.ButtonPushed_png))
	if err != nil {
		log.Fatal(err)
	}
	pushedImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	b.images.pushed = pushedImg

	return b
}
