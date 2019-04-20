package game

import (
	"image/color"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	// TODO: use some other font than this in future
)

type Label struct {
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

	// TODO: do we care?
	// Maybe. For now, just assume I'm sensible enough to have really short strings that do not overflow their bounds in any way, kinda like this comment
	size struct {
		W, H int
	}

	visible  bool
	centered bool

	text       string
	textFormat string
	textColor  color.RGBA

	updateFunc WorldFunc
}

func (l *Label) Update() error {
	if !l.visible {
		return nil
	}

	l.updateFunc(l.game.world)

	return nil
}

func (l *Label) Draw(screen *ebiten.Image) error {
	if !l.visible {
		return nil
	}

	op.GeoM.Reset()
	op.ColorM.Reset()

	// translate to parent X,Y
	/*
		op.GeoM.Translate(w.parent.Position().X, w.parent.Position().Y)
		op.GeoM.Translate(float64(w.position.X), float64(w.position.Y))
	*/

	// Default text position
	textPos := struct{ X, Y int }{
		l.position.X,
		l.position.Y + standardFont.Metrics().Ascent.Ceil(),
	}

	if l.centered {
		// TODO: center this text in the button
		// need to know how big the text is
		// for now, we know it's width is 5 px per letter
		textWidth := len(l.text) * 5
		// see: ebiten/examples/blocks/blocks/font.go for how to do this better
		textPos = struct{ X, Y int }{
			l.position.X + (l.size.W / 2) - (textWidth / 2),
			l.position.Y + standardFont.Metrics().Ascent.Ceil(),
		}
	}

	textPos.X += int(l.parent.Position().X)
	textPos.Y += int(l.parent.Position().Y)

	// TODO: use something else to draw text, but for now this is fine
	text.Draw(screen,
		l.text,
		// TODO: use some other font face
		standardFont,

		textPos.X, textPos.Y,

		l.textColor,
	)
	return nil
}

func (l *Label) IsVisible() bool {
	return l.visible
}
func (l *Label) SetVisible(v bool) {
	l.visible = v
}
func (l *Label) Game() *Game {
	return l.game
}

func (l *Label) Position() r3.Vector {
	// TODO: return w.position, once we store that as an r3.Vector natively
	return r3.Vector{
		X: float64(l.position.X),
		Y: float64(l.position.Y),
		Z: float64(l.position.Z),
	}
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewLabel(p UIElement, text, textFormat string) *Label {
	l := &Label{
		parent:     p,
		game:       p.Game(),
		text:       text,
		textFormat: textFormat,
		updateFunc: defaultWorldFunc,
	}

	return l
}
