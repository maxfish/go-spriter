package spriter

import (
	"fmt"
	"math"
)

type CurveType int

const (
	TypeLinear    CurveType = 0
	TypeInstant   CurveType = 1
	TypeQuadratic CurveType = 2
	TypeCubic     CurveType = 3
	TypeQuartic   CurveType = 4
	TypeQuintic   CurveType = 5
	TypeBezier    CurveType = 6
)

type Curve struct {
	curveType         CurveType
	constraints       [4]float64
	lastCubicSolution float64
}

func MakeCurve() *Curve {
	return &Curve{}
}

func MakeCurveWithType(curveType CurveType) *Curve {
	return &Curve{
		curveType:   curveType,
	}
}

func (c *Curve) String() string {
	return fmt.Sprintf("Curve [c1:%f, c2:%f, c3%f, c4%f]", c.constraints[0], c.constraints[1], c.constraints[2], c.constraints[3])
}

func getCurveTypeFromName(name string) CurveType {
	switch name {
	case "instant":
		return TypeInstant
	case "quadratic":
		return TypeQuadratic
	case "cubic":
		return TypeCubic
	case "quartic":
		return TypeQuartic
	case "quintic":
		return TypeQuintic
	case "bezier":
		return TypeBezier
	default:
		return TypeLinear
	}
}

func (c *Curve) interpolate(a float64, b float64, t float64) float64 {
	switch c.curveType {
	case TypeInstant:
		return a
	case TypeLinear:
		return Linear(a, b, t)
	case TypeQuadratic:
		return Quadratic(a, Linear(a, b, c.constraints[0]), b, t)
	case TypeCubic:
		return Cubic(a, Linear(a, b, c.constraints[0]), Linear(a, b, c.constraints[1]), b, t)
	case TypeQuartic:
		return Quartic(a, Linear(a, b, c.constraints[0]), Linear(a, b, c.constraints[1]), Linear(a, b, c.constraints[2]), b, t)
	case TypeQuintic:
		return Quintic(a, Linear(a, b, c.constraints[0]), Linear(a, b, c.constraints[1]), Linear(a, b, c.constraints[2]), Linear(a, b, c.constraints[3]), b, t)
	case TypeBezier:
		cubicSolution := solveCubic(3*(c.constraints[0]-c.constraints[2])+1, 3*(c.constraints[2]-2*c.constraints[0]), 3*c.constraints[0], -t)
		if cubicSolution == -1 {
			cubicSolution = c.lastCubicSolution
		} else {
			c.lastCubicSolution = cubicSolution
		}
		return Linear(a, b, Bezier(cubicSolution, 0, c.constraints[1], c.constraints[3], 1))

	default:
		return Linear(a, b, t)
	}
}

func (c *Curve) interpolatePoints(a *Point, b *Point, t float64, target *Point) {
	target[0] = c.interpolate(a.X(), b.X(), t)
	target[1] = c.interpolate(a.Y(), b.Y(), t)
}

func (c *Curve) interpolateAngleWithSpin(a float64, b float64, t float64, spin int) float64 {
	if spin == 0 {
		return a
	}
	if spin > 0 {
		if b-a < 0 {
			b += math.Pi*2
		}
	} else if spin < 0 {
		if b-a > 0 {
			b -= math.Pi*2
		}
	}

	return c.interpolate(a, b, t)
}

func (c *Curve) interpolateAngle(a float64, b float64, t float64) float64 {
	switch c.curveType {
	case TypeInstant:
		return a
	case TypeLinear:
		return LinearAngle(a, b, t)
	case TypeQuadratic:
		return QuadraticAngle(a, LinearAngle(a, b, c.constraints[0]), b, t)
	case TypeCubic:
		return CubicAngle(a, LinearAngle(a, b, c.constraints[0]), LinearAngle(a, b, c.constraints[1]), b, t)
	case TypeQuartic:
		return QuarticAngle(a, LinearAngle(a, b, c.constraints[0]), LinearAngle(a, b, c.constraints[1]), LinearAngle(a, b, c.constraints[2]), b, t)
	case TypeQuintic:
		return QuinticAngle(a, LinearAngle(a, b, c.constraints[0]), LinearAngle(a, b, c.constraints[1]), LinearAngle(a, b, c.constraints[2]), LinearAngle(a, b, c.constraints[3]), b, t)
	case TypeBezier:
		cubicSolution := solveCubic(3*(c.constraints[0]-c.constraints[2])+1, 3*(c.constraints[2]-2*c.constraints[0]), 3*c.constraints[0], -t)
		if cubicSolution == -1 {
			cubicSolution = c.lastCubicSolution
		} else {
			c.lastCubicSolution = cubicSolution
		}
		return LinearAngle(a, b, Bezier(cubicSolution, 0, c.constraints[1], c.constraints[3], 1))
	default:
		return LinearAngle(a, b, t)
	}
}
