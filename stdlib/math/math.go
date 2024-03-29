package math

import (
	"github.com/peter-mount/go-script/packages"
	"math"
)

func init() {
	// Register instance populated with constants from math package
	packages.Register("math", &Math{
		E:                      math.E,
		Pi:                     math.Pi,
		Phi:                    math.Phi,
		Sqrt2:                  math.Sqrt2,
		SqrtE:                  math.SqrtE,
		SqrtPi:                 math.SqrtPi,
		SqrtPhi:                math.SqrtPhi,
		Ln2:                    math.Ln2,
		Log2E:                  math.Log2E,
		Ln10:                   math.Ln10,
		Log10E:                 math.Log10E,
		MaxFloat32:             math.MaxFloat32,
		SmallestNonzeroFloat32: math.SmallestNonzeroFloat32,
		MaxFloat64:             math.MaxFloat64,
		SmallestNonzeroFloat64: math.SmallestNonzeroFloat64,
		MaxInt:                 math.MaxInt,
		MinInt:                 math.MinInt,
		MaxInt8:                math.MaxInt8,
		MinInt8:                math.MinInt8,
		MaxInt16:               math.MaxInt16,
		MinInt16:               math.MinInt16,
		MaxInt32:               math.MaxInt32,
		MinInt32:               math.MinInt32,
		MaxInt64:               math.MaxInt64,
		MinInt64:               math.MinInt64,
		MaxUint:                math.MaxUint,
		MaxUint8:               math.MaxUint8,
		MaxUint16:              math.MaxUint16,
		MaxUint32:              math.MaxUint32,
		MaxUint64:              math.MaxUint64,
	})
}

// Math exposes constants and functions from the standard math package to scripts
type Math struct {
	E                      float64
	Pi                     float64
	Phi                    float64
	Sqrt2                  float64
	SqrtE                  float64
	SqrtPi                 float64
	SqrtPhi                float64
	Ln2                    float64
	Log2E                  float64
	Ln10                   float64
	Log10E                 float64
	MaxFloat32             float32
	SmallestNonzeroFloat32 float32
	MaxFloat64             float64
	SmallestNonzeroFloat64 float64
	MaxInt                 int
	MinInt                 int
	MaxInt8                int8
	MinInt8                int8
	MaxInt16               int16
	MinInt16               int16
	MaxInt32               int32
	MinInt32               int32
	MaxInt64               int64
	MinInt64               int64
	MaxUint                uint
	MaxUint8               uint8
	MaxUint16              uint16
	MaxUint32              uint32
	MaxUint64              uint64
}

func (_ Math) Abs(a float64) float64                   { return math.Abs(a) }
func (_ Math) Acos(a float64) float64                  { return math.Acos(a) }
func (_ Math) Acosh(a float64) float64                 { return math.Acosh(a) }
func (_ Math) Asin(a float64) float64                  { return math.Asin(a) }
func (_ Math) Asinh(a float64) float64                 { return math.Asinh(a) }
func (_ Math) Atan(a float64) float64                  { return math.Atan(a) }
func (_ Math) Atan2(a, b float64) float64              { return math.Atan2(a, b) }
func (_ Math) Atanh(a float64) float64                 { return math.Atanh(a) }
func (_ Math) Cbrt(a float64) float64                  { return math.Cbrt(a) }
func (_ Math) Ceil(a float64) float64                  { return math.Ceil(a) }
func (_ Math) Copysign(a, b float64) float64           { return math.Copysign(a, b) }
func (_ Math) Cos(a float64) float64                   { return math.Cos(a) }
func (_ Math) Cosh(a float64) float64                  { return math.Cosh(a) }
func (_ Math) Dim(a, b float64) float64                { return math.Dim(a, b) }
func (_ Math) Erf(a float64) float64                   { return math.Erf(a) }
func (_ Math) Erfc(a float64) float64                  { return math.Erfc(a) }
func (_ Math) Erfcinv(a float64) float64               { return math.Erfcinv(a) }
func (_ Math) Exp(a float64) float64                   { return math.Exp(a) }
func (_ Math) Exp2(a float64) float64                  { return math.Exp2(a) }
func (_ Math) Expm1(a float64) float64                 { return math.Expm1(a) }
func (_ Math) Floor(a float64) float64                 { return math.Floor(a) }
func (_ Math) Frexp(a float64) (float64, interface{})  { return math.Frexp(a) }
func (_ Math) Gamma(a float64) float64                 { return math.Gamma(a) }
func (_ Math) Hypot(a, b float64) float64              { return math.Hypot(a, b) }
func (_ Math) J0(a float64) float64                    { return math.J0(a) }
func (_ Math) J1(a float64) float64                    { return math.J1(a) }
func (_ Math) Log(a float64) float64                   { return math.Log(a) }
func (_ Math) Logb(a float64) float64                  { return math.Logb(a) }
func (_ Math) Log1p(a float64) float64                 { return math.Log1p(a) }
func (_ Math) Log2(a float64) float64                  { return math.Log2(a) }
func (_ Math) Log10(a float64) float64                 { return math.Log10(a) }
func (_ Math) Lgamma(a float64) (float64, interface{}) { return math.Lgamma(a) }
func (_ Math) Max(a, b float64) float64                { return math.Max(a, b) }
func (_ Math) Min(a, b float64) float64                { return math.Min(a, b) }
func (_ Math) Mod(a, b float64) float64                { return math.Mod(a, b) }
func (_ Math) Modf(a float64) (float64, float64)       { return math.Modf(a) }
func (_ Math) Nextafter(a, b float64) float64          { return math.Nextafter(a, b) }
func (_ Math) Pow(a, b float64) float64                { return math.Pow(a, b) }
func (_ Math) Remainder(a, b float64) float64          { return math.Remainder(a, b) }
func (_ Math) Round(a float64) float64                 { return math.Round(a) }
func (_ Math) RoundToEven(a float64) float64           { return math.RoundToEven(a) }
func (_ Math) Sin(a float64) float64                   { return math.Sin(a) }
func (_ Math) Sincos(a float64) (float64, float64)     { return math.Sincos(a) }
func (_ Math) Sinh(a float64) float64                  { return math.Sinh(a) }
func (_ Math) Sqrt(a float64) float64                  { return math.Sqrt(a) }
func (_ Math) Tan(a float64) float64                   { return math.Tan(a) }
func (_ Math) Tanh(a float64) float64                  { return math.Tanh(a) }
func (_ Math) Trunc(a float64) float64                 { return math.Trunc(a) }
func (_ Math) Y0(a float64) float64                    { return math.Y0(a) }
func (_ Math) Y1(a float64) float64                    { return math.Y1(a) }
func (_ Math) Yn(n int, a float64) float64             { return math.Yn(n, a) }
