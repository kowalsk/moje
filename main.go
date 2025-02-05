// main.go
package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &Game{
		x:      100,
		y:      100,
		width:  32,
		height: 32,
		angle:  0,
		speed:  2.0,
		stationarySquares: []Square{
			{x: 200, y: 200, w: 64, h: 64, color: color.RGBA{255, 0, 0, 255}}, // Red square
		},
	}

	ebiten.SetWindowTitle("Spaceship Movement")
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(game); err != nil && err != ErrGameClosed {
		log.Fatal(err)
	}
}
