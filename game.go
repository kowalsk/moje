// game.go
package main

import (
	"errors"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var ErrGameClosed = errors.New("game closed by user") // Define a custom error to quit the game

type Game struct { // Game represents the game state
	x, y          float64 // Player position
	width, height float64 // Player dimensions
	angle         float64 // Player rotation angle in degrees
	speed         float64 // Player movement speed

	stationarySquares []Square // List of stationary squares
}

func (g *Game) Update() error { // Update is called every frame (60 times per second by default)
	originalX := g.x
	originalY := g.y

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrGameClosed
	}

	// Rotational movement (left/right)
	rotationSpeed := 2.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.angle -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.angle += rotationSpeed
	}

	// Forward/backward movement
	moveSpeed := g.speed
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.x += moveSpeed * math.Cos(g.angle*math.Pi/180)
		g.y += moveSpeed * math.Sin(g.angle*math.Pi/180)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.x -= moveSpeed * math.Cos(g.angle*math.Pi/180)
		g.y -= moveSpeed * math.Sin(g.angle*math.Pi/180)
	}

	// Collision detection
	for _, square := range g.stationarySquares {
		if checkCollision(g.x, g.y, g.width, g.height, square.x, square.y, square.w, square.h) {
			g.x = originalX
			g.y = originalY
			break
		}
	}

	return nil
}

// Draw is called every frame after Update and is responsible for rendering
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 10, 255})

	// Draw player as a triangle
	triangle := []ebiten.Vertex{
		{DstX: float32(g.x), DstY: float32(g.y - g.height/2), ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},       // Tip
		{DstX: float32(g.x - g.width/2), DstY: float32(g.y + g.height/2), ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1}, // Left
		{DstX: float32(g.x + g.width/2), DstY: float32(g.y + g.height/2), ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1}, // Right
	}
	screen.DrawTriangles(triangle, []uint16{0, 1, 2}, ebiten.NewImage(1, 1), nil)

	// Draw obstacles
	for _, square := range g.stationarySquares {
		vector.DrawFilledRect(screen, float32(square.x), float32(square.y), float32(square.w), float32(square.h), square.color, false)
	}
}

// Layout defines the screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
