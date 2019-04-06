package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/game"
	flag "github.com/spf13/pflag"
)

var (
	phoneMode = flag.Bool("phone", false, "Set resolution to approximate an iPhone XR")
)

func init() {
	flag.Parse()
}

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

	if *phoneMode {
		// https://www.paintcodeapp.com/news/ultimate-guide-to-iphone-resolutions
		game.ScreenWidth = int(414.0 / scaleFactor)
		// Includes 100 pt to account for notch
		// Meaured with the totally accurate "put phone against screen" technique
		game.ScreenHeight = int((896.0 - 100.0) / scaleFactor)
	}

	g := game.NewGame()
	if err := ebiten.Run(g.Update, game.ScreenWidth, game.ScreenHeight, scaleFactor, "Egg Garden"); err != nil && err != game.RegularTermination {
		log.Fatal(err)
	}
}
