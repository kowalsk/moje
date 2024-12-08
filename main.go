package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Create a new game instance
	game := &Game{
		rectW: 32, // Player's width
		rectH: 32, // Player's height
		stationarySquares: []Square{
			{x: 100, y: 100, w: 64, h: 64, color: color.RGBA{255, 0, 0, 255}}, // Red square
			{x: 200, y: 200, w: 64, h: 64, color: color.RGBA{0, 255, 0, 255}}, // Green square
			{x: 300, y: 150, w: 64, h: 64, color: color.RGBA{0, 0, 255, 255}}, // Blue square
		},
	}

	ebiten.SetWindowTitle("Ebiten Rectangle Example")              // Set the window title
	ebiten.SetWindowSize(640, 480)                                 // Set initial window size
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled) // Enable window resizing
	ebiten.SetFullscreen(true)                                     // Enable fullscreen mode at startup

	if err := ebiten.RunGame(game); err != nil { // Run the game, with a 320x240 logical screen size and 60 FPS ?
		if err != ErrGameClosed { // Check if the error is the custom quit error, and don't log it
			log.Fatal(err)
		}
	}
}
