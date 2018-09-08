package spriter

import (
	"math"
)

func ternary(c bool, a float64, b float64) float64 {
	if c {
		return a
	} else {
		return b
	}
}

func signum(f float64) float64 {
	return ternary(f == 0.0 || math.IsNaN(f), f, math.Copysign(1.0, f))
}

func angleDifference(a float64, b float64) float64 {
	return math.Min((2*math.Pi)-math.Abs(a-b), math.Abs(a-b))
}

// From https://github.com/Trixt0r/spriter
func solveCubic(a float64, b float64, c float64, d float64) float64 {
	if a == 0 {
		return solveQuadratic(b, c, d)
	}
	if d == 0 {
		return 0
	}

	b /= a
	c /= a
	d /= a
	squaredB := b * b
	q := (3.0*c - squaredB) / 9.0
	r := (-27.0*d + b*(9.0*c-2.0*squaredB)) / 54.0
	discriminant := (q * q * q) + (r * r)
	term1 := b / 3.0

	if discriminant > 0 {
		sqrtDisc := math.Sqrt(discriminant)
		s := r + sqrtDisc
		s = ternary(s < 0, -math.Cbrt(-s), math.Cbrt(s))
		t := r - sqrtDisc
		t = ternary(t < 0, -math.Cbrt(-t), math.Cbrt(t))

		result := -term1 + s + t
		if result >= 0 && result <= 1 {
			return result
		}
	} else if discriminant == 0 {
		r13 := ternary(r < 0, -math.Cbrt(-r), math.Cbrt(r))

		result := -term1 + 2.0*r13
		if result >= 0 && result <= 1 {
			return result
		}
		result = -(r13 + term1)
		if result >= 0 && result <= 1 {
			return result
		}
	} else {
		q = -q
		dum1 := q * q * q
		dum1 = math.Acos(r / math.Sqrt(dum1))
		r13 := 2.0 * math.Sqrt(q)

		result := -term1 + r13*math.Cos(dum1/3.0)
		if result >= 0 && result <= 1 {
			return result
		}

		result = -term1 + r13*math.Cos((dum1+2.0*math.Pi)/3.0)
		if result >= 0 && result <= 1 {
			return result
		}

		result = -term1 + r13*math.Cos((dum1+4.0*math.Pi)/3.0)
		if result >= 0 && result <= 1 {
			return result
		}
	}

	return -1
}

func solveQuadratic(a float64, b float64, c float64) float64 {
	squaredB := b * b
	twoA := 2 * a
	fourAC := 4 * a * c
	squareRoot := math.Sqrt(squaredB - fourAC)
	result := (-b + squareRoot) / twoA
	if result >= 0 && result <= 1 {
		return result
	}

	result = (-b - squareRoot) / twoA
	if result >= 0 && result <= 1 {
		return result
	}

	return -1
}

func Linear(a float64, b float64, t float64) float64 {
	return (b-a)*t + a
}

func InverseLinear(a float64, b float64, x float64) float64 {
	return (x - a) / (b - a)
}

func Quadratic(a float64, b float64, c float64, t float64) float64 {
	return Linear(Linear(a, b, t), Linear(b, c, t), t)
}

func Cubic(a float64, b float64, c float64, d float64, t float64) float64 {
	return Linear(Quadratic(a, b, c, t), Quadratic(b, c, d, t), t)
}

func Quartic(a float64, b float64, c float64, d float64, e float64, t float64) float64 {
	return Linear(Cubic(a, b, c, d, t), Cubic(b, c, d, e, t), t)
}
func Quintic(a float64, b float64, c float64, d float64, e float64, f float64, t float64) float64 {
	return Linear(Quartic(a, b, c, d, e, t), Quartic(b, c, d, e, f, t), t)
}

func LinearAngle(a float64, b float64, t float64) float64 {
	return a + angleDifference(b, a)*t
}

func QuadraticAngle(a float64, b float64, c float64, t float64) float64 {
	return LinearAngle(LinearAngle(a, b, t), LinearAngle(b, c, t), t)
}
func CubicAngle(a float64, b float64, c float64, d float64, t float64) float64 {
	return LinearAngle(QuadraticAngle(a, b, c, t), QuadraticAngle(b, c, d, t), t)
}

func QuarticAngle(a float64, b float64, c float64, d float64, e float64, t float64) float64 {
	return LinearAngle(CubicAngle(a, b, c, d, t), CubicAngle(b, c, d, e, t), t)
}

func QuinticAngle(a float64, b float64, c float64, d float64, e float64, f float64, t float64) float64 {
	return LinearAngle(QuarticAngle(a, b, c, d, e, t), QuarticAngle(b, c, d, e, f, t), t)
}

func Bezier(t float64, x1 float64, x2 float64, x3 float64, x4 float64) float64 {
	temp := t * t
	bezier0 := -temp*t + 3*temp - 3*t + 1
	bezier1 := 3*t*temp - 6*temp + 3*t
	bezier2 := -3*temp*t + 3*temp
	bezier3 := t * t * t

	return bezier0*x1 + bezier1*x2 + bezier2*x3 + bezier3*x4
}
