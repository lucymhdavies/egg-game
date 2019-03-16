package game

import (
	"bytes"
	"image"
	"math/rand"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"

	"github.com/lucymhdavies/egg-game/resources/sprites"
)

type State int

const (
	StateUnhatched State = iota // Unused
	StateIdle
	StateBounce
	StateDead
	StateRespawning
	StateSick  // Unused. Something like, if its health is low
	StateEat   // Unused
	StateSleep // Unused
)

var (
	op = &ebiten.DrawImageOptions{}

	// TODO: change both of these depending on mood?
	bounceChance   = 0.01
	maxBounceSpeed = 0.5

	// How old does it have to be before it potentially dies of old age
	// For now, hard code something really short
	minOldAgeDeath = 1.0

	// how likely it is that old age will lower health
	oldAgeSicknessChance = 0.1

	names = []string{
		"Aleggsandra", "Deggniel", "Eggberta", "Egglan", "Egglizabeth",
		"Eggsmerelda", "Gordeggn", "Heleggna", "Llywelegg", "Sabreggna",
		"Sveggn", "Eggsy",
	}
)

type Egg struct {
	velocity r3.Vector
	position r3.Vector
	size     r2.Point
	images   struct {
		body     *ebiten.Image
		eyeBall  *ebiten.Image
		eyePupil *ebiten.Image
		shadow   *ebiten.Image
		emote    *ebiten.Image
	}
	name  string
	world *World
	state State

	// Some stats we could use for these eggs
	stats struct {
		// The older it is, the more likely it is to die of old age
		// (i.e. randomly drop in health)
		age float64
		// Once it gets to health 0, it dies
		health uint8

		// Is this one of those creepy eggs?
		// Not sure what would trigger creepiness in eggs, but we shall see
		creepy bool

		//
		// Unused below, so far
		//

		// Once this reaches 255, start losing health
		hunger uint8
		// Increases based on food input + time
		bladder   uint8
		tiredness uint8
		comfort   uint8
		social    uint8
		hygene    uint8
	}
}

func NewEgg(w *World) *Egg {
	e := &Egg{
		position: r3.Vector{
			float64(w.Width) / 2,
			float64(w.Height) / 2,
			0,
		},
		world: w,
	}

	e.name = names[rand.Intn(len(names))]

	// Set default stats
	e.stats.health = 255

	// TODO: random chance of being a creepy egg
	//e.stats.creepy = true

	//
	// TODO: all this sprite loading stuff should be in init()
	//

	op.GeoM.Reset()
	op.ColorM.Reset()

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

	// Eyes
	img, _, err = image.Decode(bytes.NewReader(sprites.Eyeball_png))
	if err != nil {
		log.Fatal(err)
	}
	eyeBallImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	e.images.eyeBall = eyeBallImg

	// Pupils
	img, _, err = image.Decode(bytes.NewReader(sprites.EyePupil_png))
	if err != nil {
		log.Fatal(err)
	}
	eyePupilImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	e.images.eyePupil = eyePupilImg

	return e
}

func (egg *Egg) Update() error {

	if egg.state == StateDead {
		return nil
	}

	var deltaTime float64
	if ebiten.CurrentTPS() > 0 {
		deltaTime = 1 / ebiten.CurrentTPS()
	} else {
		return nil
	}

	if egg.state != StateRespawning {
		if egg.state != StateBounce && egg.stats.health == 0 {
			egg.state = StateDead
			egg.name = "DEAD"
			return nil
		}
		// Age the egg by deltatime
		egg.stats.age += deltaTime
	}

	// if we've entered "old age", start randomly losing health
	// if our health is already 0, do nothing, otherwise we'll get integer underflow
	if egg.stats.age > minOldAgeDeath && egg.stats.health > 0 {
		// the older the egg is beyond its "minimum old age" age,
		// the more health it should lose
		sicknessChance := (oldAgeSicknessChance * (egg.stats.age - minOldAgeDeath)) / minOldAgeDeath

		if sicknessChance > 1.0 || rand.Float64() <= sicknessChance {
			egg.stats.health--
		}
	}

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
			egg.position.X = egg.size.X/2 + 10
		}
		if egg.position.Y < egg.size.Y/2+10 {
			egg.position.Y = egg.size.Y/2 + 10
		}
		if egg.position.X > float64(egg.world.Width)-egg.size.X/2-10 {
			egg.position.X = float64(egg.world.Width) - egg.size.X/2 - 10
		}
		if egg.position.Y > float64(egg.world.Height)-egg.size.Y/2-10 {
			egg.position.Y = float64(egg.world.Height) - egg.size.Y/2 - 10
		}

		// Don't go through the floor!
		if egg.position.Z < 0 {
			egg.position.Z = 0
			egg.velocity = r3.Vector{0, 0, 0}
			egg.state = StateIdle
		} else {
			egg.velocity.Z -= 3 * deltaTime
		}

	case StateRespawning:
		// TODO: maybe accellerate, rather than just linearly going up
		// Also make these variables up top, rather than Magic Numbers
		egg.position.Z += 30 * deltaTime

		if egg.position.Z >= 30 {
			egg.world.ReplaceEgg()
		}
	}

	return nil
}

func (egg *Egg) Draw(screen *ebiten.Image) error {
	// TODO: split the drawing up into smaller functions
	// Maybe have types for different bits of the body?
	// Probably overkill.

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Draw shadow
	op.GeoM.Translate(egg.position.X-egg.size.X/2, egg.position.Y+egg.size.Y/3)
	if egg.position.Z > 0 {
		// TODO: maybe have the shadow shrink proportional to height instead
		op.ColorM.Scale(1, 1, 1, 0.5-egg.position.Z/50)
	}
	screen.DrawImage(egg.images.shadow, op)

	//
	// Draw body offscreen
	//

	// Create a temporary image to draw body + body parts
	bodyImg, _ := ebiten.NewImageFromImage(egg.images.body, ebiten.FilterDefault)

	// Sizes for body parts
	// TODO: maybe store these in egg struct?
	eyeSizeX, eyeSizeY := egg.images.eyeBall.Size()
	pupilSizeX, pupilSizeY := egg.images.eyePupil.Size()

	// Left eyeball
	op.GeoM.Reset()
	op.ColorM.Reset()
	op.GeoM.Translate(
		egg.size.X/2-float64(eyeSizeX)-2,
		egg.size.Y/2-float64(eyeSizeY),
	)
	bodyImg.DrawImage(egg.images.eyeBall, op)

	if !egg.stats.creepy {
		// TODO: move these, depending on direction of movement and/or mouse position
		// Have them either centred, or at a max distance from centre of eye
		// i.e. add:
		// +egg.velocity.Y*5
		// Doesn't quite look right though
		op.GeoM.Translate(
			float64(eyeSizeX)/2-float64(pupilSizeX)/2,
			float64(eyeSizeY)/2-float64(pupilSizeY)/2,
		)
		bodyImg.DrawImage(egg.images.eyePupil, op)
	}

	// Right eyeball
	op.GeoM.Reset()
	op.GeoM.Translate(
		egg.size.X/2+2,
		egg.size.Y/2-float64(eyeSizeY),
	)
	bodyImg.DrawImage(egg.images.eyeBall, op)

	if !egg.stats.creepy {
		// TODO: as above
		op.GeoM.Translate(
			float64(eyeSizeX)/2-float64(pupilSizeX)/2,
			float64(eyeSizeY)/2-float64(pupilSizeY)/2,
		)
		bodyImg.DrawImage(egg.images.eyePupil, op)
	}

	// Draw our temporary body image
	op.GeoM.Reset()
	op.ColorM.Reset()

	op.GeoM.Translate(egg.position.X-egg.size.X/2, egg.position.Y-egg.size.Y/2-egg.position.Z)
	if egg.state == StateDead {
		op.ColorM.Scale(255, 255, 255, 0.5)
	}
	if egg.state == StateRespawning {
		// TODO: scale transparency, inversely proportional to position.Z
		op.ColorM.Scale(255, 255, 255, 0.5-egg.position.Z/30)
	}

	screen.DrawImage(bodyImg, op)

	// Draw name?
	// Need to figure out how to center the text
	if egg.state != StateRespawning {
		ebitenutil.DebugPrintAt(screen, egg.name, int(egg.position.X-egg.size.X/2), int(egg.position.Y-egg.position.Z-egg.size.Y/2-20))
	}

	return nil
}
