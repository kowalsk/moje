package main

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var ErrGameClosed = errors.New("game closed by user") // Define a custom error to quit the game

type Game struct { // Game represents the game state
	rectX float64 // player position?
	rectY float64 // player position?
	rectW float64 // Player width
	rectH float64 // Player height

	targetX           float64 // where is player going?
	targetY           float64 // where is player going?
	isPlayerMovingToMouseClickPoint  bool
	stationarySquares []Square // List of stationary squares
}

func (g *Game) Update() error { // Update is called every frame (60 times per second by default)
	originalX := g.rectX
	originalY := g.rectY

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrGameClosed
	}

	// Key-based movement cancels the target movement
	moveSpeed := 2.0
	moveX := 0.0
	moveY := 0.0

	// Keyboard movement
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveY -= moveSpeed
		g.isPlayerMovingToMouseClickPoint = false // Don't move to target when using keyboard
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveY += moveSpeed
		g.isPlayerMovingToMouseClickPoint = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveX += moveSpeed
		g.isPlayerMovingToMouseClickPoint = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveX -= moveSpeed
		g.isPlayerMovingToMouseClickPoint = false
	}

	// Move on the X-axis first and check for collisions
	g.rectX += moveX
	collisionX := false
	for _, square := range g.stationarySquares {
		if checkCollision(g.rectX, g.rectY, g.rectW, g.rectH, square.x, square.y, square.w, square.h) {
			collisionX = true
			break
		}
	}
	// If there's a collision on the X-axis, revert the movement
	if collisionX {
		g.rectX = originalX
	}

	// Move on the Y-axis and check for collisions
	g.rectY += moveY
	collisionY := false
	for _, square := range g.stationarySquares {
		if checkCollision(g.rectX, g.rectY, g.rectW, g.rectH, square.x, square.y, square.w, square.h) {
			collisionY = true
			break
		}
	}
	// If there's a collision on the Y-axis, revert the movement
	if collisionY {
		g.rectY = originalY
	}

	// Mouse-based movement: set the target and enable movement to target
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		g.targetX = float64(mouseX)
		g.targetY = float64(mouseY)
		g.isPlayerMovingToMouseClickPoint = true
	}

	// Move towards the target only if isMovingToTarget is true
	if g.isPlayerMovingToMouseClickPoint {
		originalX = g.rectX // Store original position for collision checks
		originalY = g.rectY

		// Move towards the target (mouse-based)
		speed := moveSpeed

		// Move X towards the target and check for collisions
		if g.rectX < g.targetX {
			g.rectX += speed
			if g.rectX > g.targetX {
				g.rectX = g.targetX
			}
		} else if g.rectX > g.targetX {
			g.rectX -= speed
			if g.rectX < g.targetX {
				g.rectX = g.targetX
			}
		}

		// Check for X-axis collision
		collisionX = false
		for _, square := range g.stationarySquares {
			if checkCollision(g.rectX, g.rectY, g.rectW, g.rectH, square.x, square.y, square.w, square.h) {
				collisionX = true
				break
			}
		}
		// If there's a collision on the X-axis, revert the movement
		if collisionX {
			g.rectX = originalX
		}

		// Move Y towards the target and check for collisions
		if g.rectY < g.targetY {
			g.rectY += speed
			if g.rectY > g.targetY {
				g.rectY = g.targetY
			}
		} else if g.rectY > g.targetY {
			g.rectY -= speed
			if g.rectY < g.targetY {
				g.rectY = g.targetY
			}
		}

		// Check for Y-axis collision
		collisionY = false
		for _, square := range g.stationarySquares {
			if checkCollision(g.rectX, g.rectY, g.rectW, g.rectH, square.x, square.y, square.w, square.h) {
				collisionY = true
				break
			}
		}
		// If there's a collision on the Y-axis, revert the movement
		if collisionY {
			g.rectY = originalY
		}
	}

	return nil
}

// Draw is called every frame after Update and is responsible for rendering
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with a solid color (optional, black by default)
	screen.Fill(color.RGBA{10, 10, 10, 255}) // Very dark grey instead of pure black

	// Define border thickness
	borderThickness := 1.0

	// Get screen dimensions from Layout
	w, h := g.Layout(0, 0)
	sw, sh := float32(w), float32(h)

	// Grid settings
	cellSize := float32(32)                 // Each cell is 32x32 (same size as the rectangle)
	gridColor := color.RGBA{50, 50, 50, 30} // dark grey with low alpha transparency

	// Draw the grid lines
	for x := cellSize; x < sw; x += cellSize {
		vector.DrawFilledRect(screen, x, 0, float32(borderThickness), sh, gridColor, false)
	}
	for y := cellSize; y < sh; y += cellSize {
		vector.DrawFilledRect(screen, 0, y, sw, float32(borderThickness), gridColor, false)
	}

	// Draw the yellow rectangle (your player)
	vector.DrawFilledRect(screen, float32(g.rectX), float32(g.rectY), 32, 32, color.RGBA{200, 200, 0, 255}, false)

	// Draw the stationary squares
	for _, square := range g.stationarySquares {
		vector.DrawFilledRect(screen, float32(square.x), float32(square.y), float32(square.w), float32(square.h), square.color, false)
	}

	// Draw the borders
	borderColor := color.RGBA{128, 128, 128, 255} // Grey color

	// Top border
	vector.DrawFilledRect(screen, 0, 0, sw, float32(borderThickness), borderColor, false)

	// Bottom border
	vector.DrawFilledRect(screen, 0, sh-float32(borderThickness), sw, float32(borderThickness), borderColor, false)

	// Left border
	vector.DrawFilledRect(screen, 0, 0, float32(borderThickness), sh, borderColor, false)

	// Right border
	vector.DrawFilledRect(screen, sw-float32(borderThickness), 0, float32(borderThickness), sh, borderColor, false)
}

// Layout is used to define the screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the game's logical screen size
	// return 640, 480
	return outsideWidth, outsideHeight
}
