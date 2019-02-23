package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/egg-game/game"
)

func main() {
	g := game.NewGame()
	if err := ebiten.Run(g.Update, game.ScreenWidth, game.ScreenHeight, 2, "Egg Garden"); err != nil && err != game.RegularTermination {
		log.Fatal(err)
	}
}
