package math

import "math"

const (
	deg2Rad = math.Pi / 180.0
	rad2Dec = 180.0 / math.Pi
)

func (_ Math) Rad(a float64) float64 { return a * deg2Rad }

func (_ Math) Deg(a float64) float64 { return a * rad2Dec }
