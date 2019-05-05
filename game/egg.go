package game

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"strings"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/lucymhdavies/egg-game/resources/sprites"
)

type State int

const (
	StateUnhatched State = iota // Unused
	StateIdle
	StateBounce
	StateEat
	StateDead
	StateRespawning
	StateSick  // Unused. Something like, if its health is low
	StateSleep // Unused
)

var (

	// TODO: change both of these depending on mood?
	defaultBounceChance = 0.01
	maxBounceSpeed      = 0.5

	maxEatDistance = 50.0

	// How old does it have to be before it potentially dies of old age
	// For now, hard code something really short
	minOldAgeDeath = 600.0 // 10 minutes

	// how likely it is that old age will lower health
	oldAgeSicknessChance = 0.1

	// How many seconds between bites of food
	timeBetweenBites = 0.5

	names = []string{
		"Aleggsandra", "Deggniel", "Eggberta", "Egglan", "Egglizabeth",
		"Eggsmerelda", "Gordeggn", "Heleggna", "Llywelegg", "Sabreggna",
		"Sveggn", "Eggsy",
	}
)

type Egg struct {
	velocity r3.Vector
	position r3.Vector // refers to midpoint of egg
	size     r2.Point
	images   struct {
		body     *ebiten.Image
		eyeBall  *ebiten.Image
		eyePupil *ebiten.Image
		// Tempoary, while the composite egg body does not change
		bodyFull *ebiten.Image

		shadow *ebiten.Image
	}
	name  string
	world *World
	state State

	// TODO: figure out how I want this to work!
	status map[string]Emote

	// How many seconds before egg can bite again
	// (only used during stateEat)
	timeUntilBite float64

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

		// How hungry is the egg?
		// 255 = not hungry, 0 = starving
		// Once this reaches 0, start losing health
		hunger uint8

		// How long before egg can eat again once full?
		saturation uint8

		//
		// Unused below, so far
		//

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
		world:  w,
		status: make(map[string]Emote),
	}

	e.name = names[rand.Intn(len(names))]

	// Set default stats
	e.stats.health = 255
	e.stats.hunger = 255

	// TODO: random chance of being a creepy egg
	//e.stats.creepy = true

	//
	// TODO: all this sprite loading stuff should be in init()
	//

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Pick a body image at random
	var keys []string
	for k := range sprites.Eggs {
		keys = append(keys, k)
	}
	n := rand.Int() % len(keys)
	key := keys[n]

	// Get body image
	img, _, err := image.Decode(bytes.NewReader(sprites.Eggs[key]))
	if err != nil {
		log.Fatal(err)
	}
	bodyImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	sizeX, sizeY := bodyImg.Size()
	e.size = r2.Point{float64(sizeX), float64(sizeY)}
	e.images.body = bodyImg

	// Shadow = body, but squashed vertically
	shadowImgRaw, _ := ebiten.NewImage(sizeX, sizeY/4, ebiten.FilterDefault)
	shadowImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault) // TODO: just use bodyImg
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

	//
	// Hunger
	//

	// decrement hunger bar over time
	// TODO: sometimes decrement this slower
	// e.g. if asleep, or saturated
	if egg.stats.hunger > 0 {
		hungerChance := 0.1
		if rand.Float64() <= hungerChance {
			// If egg is saturated, decrement that first
			if egg.stats.saturation > 0 {
				egg.status["saturated"] = emotes["saturated"]

				egg.stats.saturation--
			} else {
				delete(egg.status, "saturated")

				egg.stats.hunger--

			}
		}
	} else {
		if egg.stats.health > 0 {
			starveChace := 0.1
			if rand.Float64() <= starveChace {
				egg.stats.health--
			}
		}
	}

	//
	// Egg States
	// TODO: refactor these out as a FSM
	//

	switch egg.state {
	case StateUnhatched:
		// Hatching not yet implemented, so just go straight to idle
		egg.state = StateIdle

	case StateIdle:
		// look for nearest food
		var nearestFood *Food

		// Initialise our bounce chance, based on default
		bounceChance := defaultBounceChance

		// If the egg is hungry, and there is food in the world...
		if egg.stats.hunger < 255 {
			nearestFood = egg.world.NearestFood(egg.position)

			if nearestFood != nil {

				vecFromEggToFood := nearestFood.position.Sub(egg.position)

				distance := vecFromEggToFood.Norm()

				if distance < maxEatDistance {
					bounceChance = 0.0
					egg.state = StateEat
				} else {
					bounceChance = 1.0
				}
			}
		}

		if rand.Float64() <= bounceChance {
			egg.state = StateBounce
			egg.velocity.Z += 1.0

			// If there is no nearestFood (or we've not checked...)
			if nearestFood == nil {
				// Random direction
				egg.velocity.X = maxBounceSpeed * (rand.Float64()*2 - 1)
				egg.velocity.Y = maxBounceSpeed * (rand.Float64()*2 - 1)
			} else {
				// Bounce towards food
				vecFromEggToFood := nearestFood.position.Sub(egg.position)

				desiredVelocity := vecFromEggToFood.Normalize()
				desiredVelocity = desiredVelocity.Mul(maxBounceSpeed)

				egg.velocity.X = desiredVelocity.X
				egg.velocity.Y = desiredVelocity.Y
			}
		}

	case StateBounce:
		egg.position = egg.position.Add(egg.velocity)

		// Don't go outside bounds
		if egg.position.X < egg.size.X/2+float64(egg.world.Padding().Top) {
			egg.position.X = egg.size.X/2 + float64(egg.world.Padding().Top)
		}
		if egg.position.Y < egg.size.Y/2+float64(egg.world.Padding().Left) {
			egg.position.Y = egg.size.Y/2 + float64(egg.world.Padding().Left)
		}
		if egg.position.X > float64(egg.world.Width)-egg.size.X/2-float64(egg.world.Padding().Bottom) {
			egg.position.X = float64(egg.world.Width) - egg.size.X/2 - float64(egg.world.Padding().Bottom)
		}
		if egg.position.Y > float64(egg.world.Height)-egg.size.Y/2-float64(egg.world.Padding().Right) {
			egg.position.Y = float64(egg.world.Height) - egg.size.Y/2 - float64(egg.world.Padding().Right)
		}

		// Don't go through the floor!
		if egg.position.Z < 0 {
			egg.position.Z = 0
			egg.velocity = r3.Vector{0, 0, 0}
			egg.state = StateIdle
		} else {
			egg.velocity.Z -= 3 * deltaTime
		}

	case StateEat:
		return egg.updateEat()

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
	// Draw body offscreen unless it already exists
	//
	if egg.images.bodyFull == nil {

		// Create a temporary image to draw body + body parts
		op.GeoM.Reset()
		op.ColorM.Reset()
		bodyImg, err := ebiten.NewImage(int(egg.size.X), int(egg.size.Y), ebiten.FilterDefault)
		if err != nil {
			// TODO: better error handling needed here
			log.Fatal(err)
		}
		bodyImg.DrawImage(egg.images.body, op)

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

		egg.images.bodyFull = bodyImg
	}

	// Draw our temporary body image
	op.GeoM.Reset()
	op.ColorM.Reset()

	op.GeoM.Translate(egg.position.X-egg.size.X/2,
		egg.position.Y-egg.size.Y/2-egg.position.Z)
	if egg.state == StateDead {
		op.ColorM.Scale(255, 255, 255, 0.5)
	}
	if egg.state == StateRespawning {
		// TODO: scale transparency, inversely proportional to position.Z
		op.ColorM.Scale(255, 255, 255, 0.5-egg.position.Z/30)
	}

	screen.DrawImage(egg.images.bodyFull, op)

	if len(egg.status) > 0 {
		var emote Emote

		// find highest priorty status
		for _, status := range egg.status {
			if emote.name == "" {
				emote = status
			}

			if status.priority < emote.priority {
				emote = status
			}
		}

		statusW, statusH := emote.image.Size()

		op.GeoM.Reset()
		op.GeoM.Translate(
			egg.position.X-
				float64(statusW/4),
			egg.position.Y-
				float64(statusH)-
				egg.size.Y/2-
				5)

		screen.DrawImage(emote.image, op)

	}

	// Draw name?
	// TODO: Need to figure out how to center the text
	if egg.state != StateRespawning {
		ebitenutil.DebugPrintAt(screen, egg.name, 5, 0)
	}

	return nil
}

// TODO: GetStat(name string)

func (egg *Egg) GetStat(name string) (float64, error) {
	switch strings.ToLower(name) {
	case "health":
		return float64(egg.stats.health), nil
	case "hunger":
		return float64(egg.stats.hunger), nil
	case "saturation":
		return float64(egg.stats.saturation), nil
	case "age":
		return float64(egg.stats.age), nil
	}

	return 0, fmt.Errorf("No such stat: %s", name)
}

func (egg *Egg) updateEat() error {

	var deltaTime float64
	if ebiten.CurrentTPS() > 0 {
		deltaTime = 1 / ebiten.CurrentTPS()
	} else {
		return nil
	}

	// If egg needs to wait before taking another bite...
	if egg.timeUntilBite > 0 {

		egg.timeUntilBite = math.Max(egg.timeUntilBite-deltaTime, 0)

	} else {

		nearestFood := egg.world.NearestFood(egg.position)

		if nearestFood == nil {
			// Trivially, there's no food in the world
			// Back to Idle state
			egg.state = StateIdle
		} else {
			// Is nearest food within eating range?
			vecFromEggToFood := nearestFood.position.Sub(egg.position)
			distance := vecFromEggToFood.Norm()
			if distance >= maxEatDistance {
				// No, it's too far away.
				// Back to Idle state
				egg.state = StateIdle
			} else {
				// There is food, and it is close enough

				// Bite it!
				if nearestFood.bitesLeft > 1 {
					// there's enough food left to eat more later
					nearestFood.bitesLeft--
				} else {
					// last bite, so delete the food
					egg.world.RemoveFood(nearestFood)
				}

				// Cooldown between bites
				// (also cooldown before egg can move)
				egg.timeUntilBite += timeBetweenBites

				//
				// Food modifies egg stats
				//

				// Prevent overflow by adding them as floats, then using
				// math.Min (which needs floats) to set the final value
				newHealth := float64(egg.stats.health) + float64(nearestFood.foodType.health)
				egg.stats.health = uint8(math.Max(math.Min(255.0, newHealth), 0))

				newHunger := float64(egg.stats.hunger) + float64(nearestFood.foodType.hunger)
				egg.stats.hunger = uint8(math.Max(math.Min(255.0, newHunger), 0))

				// If we've eaten more than enough to fill the egg...
				// and this is the kind of food that leaves the egg saturated
				if newHunger > 255.0 && nearestFood.foodType.saturation {
					newSaturation := float64(egg.stats.saturation) + float64(nearestFood.foodType.hunger)
					egg.stats.saturation = uint8(math.Min(255.0, newSaturation))
				}
			}

		}

	}

	return nil
}
