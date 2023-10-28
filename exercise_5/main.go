package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	p := Point{
		x: x,
		y: y,
	}

	return &p
}

func (p1 Point) CalculateDistance(p2 *Point) float64 {
	rX := p2.x - p1.x
	rY := p2.y - p1.y

	return math.Sqrt((rX * rX) + (rY * rY))
}
func main() {
	point1 := NewPoint(2, 3)
	point2 := NewPoint(5, 7)

	result := point1.CalculateDistance(point2)
	fmt.Println(result)
}
