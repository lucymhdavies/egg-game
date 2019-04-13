package game

import (
	"image/color"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/bitmapfont"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	// TODO: use some other font than this in future
)

type Window struct {
	// Pointer back to parent UIElement
	parent UIElement

	// Pointer back to the game
	game *Game

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

	text      string
	textColor color.RGBA

	image *ebiten.Image

	// All UI elements within this window
	uiElements []UIElement
}

func (w *Window) Update() error {
	if !w.visible {
		return nil
	}

	for _, e := range w.uiElements {
		e.Update()
	}

	return nil
}

func (w *Window) Draw(screen *ebiten.Image) error {
	if !w.visible {
		return nil
	}

	op.GeoM.Reset()
	op.ColorM.Reset()

	// TODO: translate to parent X,Y
	op.GeoM.Translate(w.parent.Position().X, w.parent.Position().Y)
	op.GeoM.Translate(float64(w.position.X), float64(w.position.Y))
	op.ColorM.Scale(1, 1, 1, 0.8)
	screen.DrawImage(w.image, op)

	// TODO: center this text in the button
	// need to know how big the text is
	// for now, we know it's width is 5 px per letter
	textWidth := len(w.text) * 5
	// see: ebiten/examples/blocks/blocks/font.go for how to do this better
	textPos := struct{ X, Y int }{
		w.position.X + (w.size.W / 2) - (textWidth / 2),
		w.position.Y + 15,
	}
	textPos.X += int(w.parent.Position().X)
	textPos.Y += int(w.parent.Position().Y)

	// TODO: use something else to draw text, but for now this is fine
	text.Draw(screen,
		w.text,
		// TODO: use some other font face
		bitmapfont.Gothic12r,

		textPos.X, textPos.Y,

		w.textColor,
	)

	for _, e := range w.uiElements {
		// TODO: sort by Z-index, showing lower elements first

		e.Draw(screen)
	}

	return nil
}

func (w *Window) IsVisible() bool {
	return w.visible
}
func (w *Window) SetVisible(v bool) {
	w.visible = v
}
func (w *Window) Game() *Game {
	return w.game
}

func (w *Window) Position() r3.Vector {
	// TODO: return w.position, once we store that as an r3.Vector natively
	return r3.Vector{
		X: float64(w.position.X),
		Y: float64(w.position.Y),
		Z: float64(w.position.Z),
	}
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewWindow(p UIElement, width, height int) *Window {
	w := &Window{
		parent: p,
		game:   p.Game(),
	}

	// Default button size, if unspecified
	w.size.W = width
	w.size.H = height

	// TODO: proper window image
	w.image = NewBox(width, height, "grey_panel").Image

	return w
}
