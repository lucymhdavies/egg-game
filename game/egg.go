package game

import (
	"bytes"
	"image"
	"math/rand"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"

	"github.com/lucymhdavies/egg-game/resources/sprites"
)

type State int

const (
	StateUnhatched State = iota
	StateIdle
	StateBounce
	StateSleep
)

var (
	op = &ebiten.DrawImageOptions{}

	// TODO: change both of these depending on mood?
	bounceChance = 0.01
	//bounceChance   = 1.0
	maxBounceSpeed = 0.5
)

type Egg struct {
	velocity r3.Vector
	position r3.Vector
	size     r2.Point
	images   struct {
		body   *ebiten.Image
		eyes   *ebiten.Image
		shadow *ebiten.Image
		emote  *ebiten.Image
	}
	name  string
	world *World
	state State
}

func NewEgg(w *World) *Egg {
	e := &Egg{
		position: r3.Vector{
			float64(w.Width) / 2,
			float64(w.Height) / 2,
			0,
		},
		world: w,
		name:  "Egg",
	}

	// Get body image
	img, _, err := image.Decode(bytes.NewReader(sprites.EggBody_png))
	if err != nil {
		log.Fatal(err)
	}
	bodyImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	sizeX, sizeY := bodyImg.Size()
	e.size = r2.Point{float64(sizeX), float64(sizeY)}
	e.images.body = bodyImg

	// Shadow = body, but squashed vertically
	shadowImgRaw, _ := ebiten.NewImage(sizeX, sizeY/4, ebiten.FilterDefault)
	shadowImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	op.GeoM.Scale(1, 0.25)
	op.ColorM.Scale(0, 0, 0, 0.5)
	shadowImgRaw.DrawImage(shadowImg, op)
	e.images.shadow = shadowImgRaw

	return e
}

func (egg *Egg) Update() error {

	switch egg.state {
	case StateUnhatched:
		// Hatching not yet implemented, so just go straight to idle
		egg.state = StateIdle

	case StateIdle:
		if rand.Float64() <= bounceChance {
			egg.state = StateBounce
			egg.velocity.Z += 1

			// TODO: if there's food in the world, and hungry, go towards it

			// Random direction
			egg.velocity.X = maxBounceSpeed * (rand.Float64()*2 - 1)
			egg.velocity.Y = maxBounceSpeed * (rand.Float64()*2 - 1)
		}

	case StateBounce:
		egg.position = egg.position.Add(egg.velocity)

		// Don't go outside bounds
		if egg.position.X < egg.size.X/2+10 {
			log.Debugf("Left Boundary")
			egg.position.X = egg.size.X/2 + 10
		}
		if egg.position.Y < egg.size.Y/2+10 {
			log.Debugf("Top Boundary")
			egg.position.Y = egg.size.Y/2 + 10
		}
		if egg.position.X > float64(egg.world.Width)-egg.size.X/2-10 {
			log.Debugf("Right Boundary")
			egg.position.X = float64(egg.world.Width) - egg.size.X/2 - 10
		}
		if egg.position.Y > float64(egg.world.Height)-egg.size.Y/2-10 {
			log.Debugf("Bottom Boundary")
			egg.position.Y = float64(egg.world.Height) - egg.size.Y/2 - 10
		}

		// Don't go through the floor!
		if egg.position.Z < 0 {
			egg.position.Z = 0
			egg.velocity = r3.Vector{0, 0, 0}
			egg.state = StateIdle
		} else {
			egg.velocity.Z -= 0.05
		}
	}

	return nil
}

func (egg *Egg) Draw(screen *ebiten.Image) error {
	op.GeoM.Reset()
	op.ColorM.Reset()

	// Draw shadow
	op.GeoM.Translate(egg.position.X-egg.size.X/2, egg.position.Y+egg.size.Y/3)
	if egg.position.Z > 0 {
		// TODO: maybe have the shadow shrink proportional to height instead
		op.ColorM.Scale(1, 1, 1, 0.5-egg.position.Z/50)
	}
	screen.DrawImage(egg.images.shadow, op)

	// Draw body
	op.GeoM.Reset()
	op.ColorM.Reset()
	op.GeoM.Translate(egg.position.X-egg.size.X/2, egg.position.Y-egg.size.Y/2-egg.position.Z)
	screen.DrawImage(egg.images.body, op)
	return nil
}
