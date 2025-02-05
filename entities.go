// entities.go
package main

import "image/color"

// Square represents a stationary obstacle in the game
type Square struct {
	x     float64
	y     float64
	w     float64
	h     float64
	color color.RGBA
}

// checkCollision checks if two rectangles overlap
func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2
}
