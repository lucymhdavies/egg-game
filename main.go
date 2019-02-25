package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/game"
)

func main() {
	ebiten.SetRunnableInBackground(true)

	scaleFactor := 2.0

	if runtime.GOARCH == "js" {
		scaleFactor = ebiten.DeviceScaleFactor()
		w, h := ebiten.ScreenSizeInFullscreen()
		ebiten.SetFullscreen(true)
		game.ScreenWidth = int(float64(w) / scaleFactor)
		game.ScreenHeight = int(float64(h) / scaleFactor)

	}

	g := game.NewGame()
	if err := ebiten.Run(g.Update, game.ScreenWidth, game.ScreenHeight, scaleFactor, "Egg Garden"); err != nil && err != game.RegularTermination {
		log.Fatal(err)
	}
}
