package game

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
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
	padding Padding
	size    struct {
		W, H int
	}

	visible bool

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

	// translate to parent X,Y
	op.GeoM.Translate(w.parent.Position().X, w.parent.Position().Y)
	op.GeoM.Translate(float64(w.position.X), float64(w.position.Y))
	op.ColorM.Scale(1, 1, 1, 0.8)
	screen.DrawImage(w.image, op)

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
func (w *Window) Padding() Padding {
	return w.padding
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewWindow(p UIElement, width, height int) *Window {
	w := &Window{
		parent: p,
		game:   p.Game(),
		padding: Padding{
			Top: 10, Bottom: 10, Right: 10, Left: 10,
		},
	}

	// Default window size, if unspecified
	w.size.W = width
	w.size.H = height

	// TODO: proper window image
	w.image = NewBox(width, height, "grey_panel").Image

	return w
}
