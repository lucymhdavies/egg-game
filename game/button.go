package game

import (
	"bytes"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/lucymhdavies/egg-game/resources/sprites"
	log "github.com/sirupsen/logrus"
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

	// TODO: something to specify which image to use
	// for now, just use the one we have

	// TODO:
	// arguably this should be init()
	// or at least, have some way of storing common images used between
	// multiple buttons
	op.GeoM.Reset()
	op.ColorM.Reset()

	// Corner... store top left only
	img, _, err := image.Decode(bytes.NewReader(sprites.ButtonBlueCorner_png))
	if err != nil {
		log.Fatal(err)
	}
	cornerImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// Side corner... regular corner, but modified to #AAAAAA
	img, _, err = image.Decode(bytes.NewReader(sprites.ButtonSideCorner_png))
	if err != nil {
		log.Fatal(err)
	}
	sideCornerImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// Edge... store top edge only
	img, _, err = image.Decode(bytes.NewReader(sprites.ButtonBlueEdge_png))
	if err != nil {
		log.Fatal(err)
	}
	edgeImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// Side edge... regular edge, but modified to #AAAAAA
	img, _, err = image.Decode(bytes.NewReader(sprites.ButtonSideEdge_png))
	if err != nil {
		log.Fatal(err)
	}
	sideEdgeImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// Mid
	img, _, err = image.Decode(bytes.NewReader(sprites.ButtonBlueMid_png))
	if err != nil {
		log.Fatal(err)
	}
	midImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	//
	// Now cache the full images
	//

	op.GeoM.Reset()
	op.ColorM.Reset()

	// image when unpushed
	normalImg, _ := ebiten.NewImage(b.size.W, b.size.H, ebiten.FilterDefault)

	// image when pushed (4 pixels shorter)
	pushedImg, _ := ebiten.NewImage(b.size.W, b.size.H-4.0, ebiten.FilterDefault)

	// TODO: there's a bunch of 5s and 4s in this
	// Move these numbers to variables

	//
	// Side Corners
	// TODO: split this into its own function
	//

	// Start with the two corners
	sizeX, sizeY := sideCornerImg.Size()

	// Bottom Left
	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	// Now draw
	op.GeoM.Translate(0.0, float64(b.size.H-5))
	normalImg.DrawImage(sideCornerImg, op)

	// Bottom Right
	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi * 3 / 2)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw
	op.GeoM.Translate(float64(b.size.W-5), float64(b.size.H-5))
	normalImg.DrawImage(sideCornerImg, op)

	//
	// Side Edges
	//
	sizeX, sizeY = sideEdgeImg.Size()

	// Left Side
	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi / 2)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw
	op.GeoM.Translate(0.0, float64(b.size.H-sizeY-4))
	normalImg.DrawImage(sideEdgeImg, op)

	// Right side
	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi * 3 / 2)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw
	op.GeoM.Translate(float64(b.size.W-sizeX), float64(b.size.H-sizeY-4))
	normalImg.DrawImage(sideEdgeImg, op)

	// Bottom side
	// scale amount = (width of button - 2*width of corner (sizeX)) / width of corner
	// in this case, all tiles happen to be the same size, but...
	// TODO: refactor later ;)
	scaleAmountX := (float64(b.size.W) - float64(2*sizeX)) / float64(sizeX)
	op.GeoM.Reset()
	op.GeoM.Scale(scaleAmountX, 1.0)
	// Now draw
	op.GeoM.Translate(float64(sizeX), float64(b.size.H-sizeY))
	normalImg.DrawImage(sideEdgeImg, op)

	// TODO

	//
	// Corners
	// TODO: split this into its own function
	//

	op.GeoM.Reset()
	sizeX, sizeY = cornerImg.Size()

	// Top Left
	normalImg.DrawImage(cornerImg, op)
	pushedImg.DrawImage(cornerImg, op)

	// Top Right

	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi / 2)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw
	op.GeoM.Translate(float64(b.size.W-5), 0.0)
	normalImg.DrawImage(cornerImg, op)
	pushedImg.DrawImage(cornerImg, op)

	// Bottom Left

	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi * 3 / 2)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw
	//   Move down to bottom off button,
	//   up by 5 (height of corner),
	//   then up by 4 again (height of "side")
	op.GeoM.Translate(0.0, float64(b.size.H-5-4))
	normalImg.DrawImage(cornerImg, op)
	pushedImg.DrawImage(cornerImg, op)

	// Bottom Right

	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2)
	// Now draw (as above)
	op.GeoM.Translate(float64(b.size.W-5), float64(b.size.H-5-4))
	normalImg.DrawImage(cornerImg, op)
	pushedImg.DrawImage(cornerImg, op)

	//
	// Edges
	// TODO: split this into its own function
	//

	sizeX, sizeY = edgeImg.Size()

	// Top
	op.GeoM.Reset()
	// Stretch horizontally
	op.GeoM.Scale(scaleAmountX, 1.0)
	// Draw
	op.GeoM.Translate(5.0, 0.0)
	normalImg.DrawImage(edgeImg, op)
	pushedImg.DrawImage(edgeImg, op)
	op.GeoM.Translate(-5.0, 0.0)

	// Bottom
	op.GeoM.Translate(-float64(sizeX)/2*scaleAmountX, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi)
	op.GeoM.Translate(float64(sizeX)/2*scaleAmountX, float64(sizeY)/2)
	// Now draw (as above)
	op.GeoM.Translate(5.0, float64(b.size.H-5-4))
	normalImg.DrawImage(edgeImg, op)
	pushedImg.DrawImage(edgeImg, op)

	// Left
	// Similar to scaleAmountX, except were also taking into account the bottom "side" 4 pixels
	scaleAmountY := (float64(b.size.H) - float64(2*sizeY) - 4) / float64(sizeY)
	// Stretch vertically
	// Translate to midpoint, rotate, stretch, translate back
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2)
	op.GeoM.Rotate(math.Pi * 3 / 2)
	op.GeoM.Scale(1.0, scaleAmountY)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2*scaleAmountY)
	// Draw
	op.GeoM.Translate(0.0, 5.0)
	normalImg.DrawImage(edgeImg, op)
	pushedImg.DrawImage(edgeImg, op)
	op.GeoM.Translate(0.0, -5.0)

	// Right
	op.GeoM.Translate(-float64(sizeX)/2, -float64(sizeY)/2*scaleAmountY)
	op.GeoM.Rotate(math.Pi)
	op.GeoM.Translate(float64(sizeX)/2, float64(sizeY)/2*scaleAmountY)
	// Now draw (as above)
	op.GeoM.Translate(float64(b.size.W-5), 5.0)
	normalImg.DrawImage(edgeImg, op)
	pushedImg.DrawImage(edgeImg, op)

	//
	// Mid
	// TODO: split this out into its own function
	//
	sizeX, sizeY = midImg.Size()
	op.GeoM.Reset()
	// Stretch
	op.GeoM.Scale(scaleAmountX, scaleAmountY)
	// Now draw (as above)
	op.GeoM.Translate(5.0, 5.0)
	normalImg.DrawImage(midImg, op)
	pushedImg.DrawImage(midImg, op)

	//
	// Done
	//

	b.images.normal = normalImg
	b.images.pushed = pushedImg

	return b
}
