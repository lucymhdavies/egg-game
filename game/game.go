package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
)

var (
	// TODO: common resolutions, e.g. iPhone size
	ScreenWidth  = 256
	ScreenHeight = 256
	//ScreenHeight = 384 // room for the buttons later

	debugMode = true
)

func init() {
	rand.Seed(time.Now().UnixNano())

	if debugMode {
		log.SetLevel(log.DebugLevel)
	}
}

type Game struct {
	world *World
	input *Input
	ui    *UI
}

func NewGame() *Game {
	g := &Game{
		world: NewWorld(ScreenWidth, ScreenHeight),
	}
	g.input = &Input{ // TODO: NewInput?
		game: g,
	}
	g.ui = NewUI(g)

	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	// TODO: scenemanager stuff, e.g.
	// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/scenemanager.go

	err := g.input.Update()
	if err != nil {
		return err
	}

	err = g.ui.Update()
	if err != nil {
		return err
	}

	err = g.update()
	if err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	err = g.draw(screen)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) update() error {
	g.world.Update()

	return nil
}

func (g *Game) draw(screen *ebiten.Image) error {
	screen.Fill(color.White)

	g.world.Draw(screen)
	g.ui.Draw(screen)

	if debugMode {
		msg := fmt.Sprintf(`Age: %0.2f
Health: %d`,
			g.world.egg.stats.age,
			g.world.egg.stats.health,
		)

		if runtime.GOARCH != "js" {
			msg = msg + `
Press Q to quit`
		}

		ebitenutil.DebugPrint(screen, msg)
	}

	return nil
}
