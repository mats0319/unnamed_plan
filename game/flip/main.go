package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mats0319/unnamed_plan/game/flip/go"
)

func main() {
	flipGame, err := flip.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(flip.ScreenWidth, flip.ScreenHeight)
	ebiten.SetWindowTitle("Flip")

	err = ebiten.RunGame(flipGame)
	if err != nil {
		log.Fatal(err)
	}
}
