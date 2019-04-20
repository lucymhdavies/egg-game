package game

import (
	"bytes"
	"image"
	"log"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/resources/sprites"
)

// Bar is a progress-bar type thing for displaying numerical values
type Bar struct {
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

	value float64
	max   float64

	visible bool

	// Color of the bar
	color string

	// TODO: bars are horizontal only. allow vertical later?

	images map[string]*ebiten.Image

	label Label

	updateFunc WorldFunc
}

// TODO, functional params
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewBar(p UIElement, width, height int, color string) *Bar {
	b := &Bar{
		parent:     p,
		game:       p.Game(),
		updateFunc: defaultWorldFunc,
		color:      color,
	}
	b.size.W = width
	b.size.H = height

	// TODO: create new label
	/*
		text:       text,
		textFormat: textFormat,
	*/

	//
	// Load all the images
	//

	b.images = make(map[string]*ebiten.Image)

	for _, color := range []string{"back", color} {
		for _, position := range []string{
			"mid", "left", "right",
		} {

			imageName := color + "_" + position

			img, _, err := image.Decode(bytes.NewReader(sprites.ProgressBarH[imageName]))
			if err != nil {
				// TODO: better error handling needed here
				log.Fatal(err)
			}
			b.images[imageName], _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		}
	}

	return b
}

func (b *Bar) IsVisible() bool {
	return b.visible
}
func (b *Bar) SetVisible(v bool) {
	b.visible = v
}
func (b *Bar) Game() *Game {
	return b.game
}

func (b *Bar) Position() r3.Vector {
	// TODO: return w.position, once we store that as an r3.Vector natively
	return r3.Vector{
		X: float64(b.position.X),
		Y: float64(b.position.Y),
		Z: float64(b.position.Z),
	}
}

func (b *Bar) Update() error {
	if !b.visible {
		return nil
	}

	b.updateFunc(b.game.world)

	return nil
}

func (b *Bar) Draw(screen *ebiten.Image) error {
	if !b.visible {
		return nil
	}

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Create a new empty image
	barImage, err := ebiten.NewImage(b.size.W, b.size.H, ebiten.FilterDefault)
	if err != nil {
		// TODO: better error handling needed here
		log.Fatal(err)
	}

	//
	// Bar Background
	//

	lSizeX, _ := b.images["back_left"].Size()
	rSizeX, _ := b.images["back_right"].Size()
	mSizeX, _ := b.images["back_mid"].Size()

	// Left
	barImage.DrawImage(b.images["back_left"], op)

	// Mid
	op.GeoM.Reset()
	scaleAmountX := (float64(b.size.W) - float64(lSizeX+rSizeX)) / float64(mSizeX)
	op.GeoM.Scale(scaleAmountX, 1.0)
	op.GeoM.Translate(float64(lSizeX), 0.0)
	barImage.DrawImage(b.images["back_mid"], op)

	// Right
	op.GeoM.Reset()
	op.GeoM.Translate(float64(b.size.W-rSizeX), 0.0)
	barImage.DrawImage(b.images["back_right"], op)

	if b.value > 0 {

		//
		// Bar Foreground
		//

		// For now, when displaying the bar foreground, we're only stretching
		// and squishing the middle image.
		// This is kinda okay, but makes the transition from 1% to 0% kinda weird.
		//
		// TODO: Need to do something about that. Not sure what though
		// One idea: squish the bar ends vertically as well as horizontally
		// works fine with the current placeholder art, but not sure what art
		// I'm going to have in future.

		lSizeX, _ = b.images[b.color+"_left"].Size()
		rSizeX, _ = b.images[b.color+"_right"].Size()
		mSizeX, _ = b.images[b.color+"_mid"].Size()

		// Left
		op.GeoM.Reset()
		barImage.DrawImage(b.images[b.color+"_left"], op)

		// Mid
		op.GeoM.Reset()
		scaleAmountX = (b.value / b.max) * (float64(b.size.W) - float64(lSizeX+rSizeX)) / float64(mSizeX)
		op.GeoM.Scale(scaleAmountX, 1.0)
		op.GeoM.Translate(float64(lSizeX), 0.0)
		barImage.DrawImage(b.images[b.color+"_mid"], op)

		// Right
		op.GeoM.Reset()
		barMidWidth := scaleAmountX * float64(mSizeX)
		op.GeoM.Translate(float64(lSizeX)+barMidWidth, 0.0)
		barImage.DrawImage(b.images[b.color+"_right"], op)

	}

	//
	// Draw final bar
	//

	op.GeoM.Reset()
	op.ColorM.Reset()
	// translate to parent X,Y
	op.GeoM.Translate(b.parent.Position().X, b.parent.Position().Y)
	op.GeoM.Translate(float64(b.position.X), float64(b.position.Y))

	screen.DrawImage(barImage, op)

	return nil
}
