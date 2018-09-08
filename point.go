package spriter

import (
	"fmt"
	"math"
	)

type Point [2]float64

func MakePoint(x float64, y float64) *Point {
	return &Point{x, y}
}

func (p *Point) MakeCopy() *Point {
	return &Point{p[0], p[1]}
}

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) Y() float64 {
	return p[1]
}

func (p *Point) SetCoords(x float64, y float64) *Point {
	p[0] = x
	p[1] = y
	return p
}

func (p *Point) Set(other *Point) *Point {
	p[0] = other[0]
	p[1] = other[1]
	return p
}

func (p *Point) ScaleCoords(x float64, y float64) *Point {
	p[0] *= x
	p[1] *= y
	return p
}

func (p *Point) Scale(other *Point) *Point {
	p[0] *= other[0]
	p[1] *= other[1]
	return p
}

func (p *Point) Add(b *Point) *Point {
	p[0] += b[0]
	p[1] += b[1]
	return p
}

func (p *Point) Sub(b *Point) *Point {
	p[0] -= b[0]
	p[1] -= b[1]
	return p
}

func (p *Point) Rotate(degrees float64) *Point {
	if p[0] != 0 || p[1] != 0 {
		cos := math.Cos(degrees)
		sin := math.Sin(degrees)

		xx := p[0]*cos - p[1]*sin
		yy := p[0]*sin + p[1]*cos
		p[0] = xx
		p[1] = yy
	}
	return p
}

func (p *Point) String() string {
	return fmt.Sprintf("[%f,%f]", p[0], p[1])
}
