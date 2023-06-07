package math

import (
	"github.com/peter-mount/go-script/executor"
	"math"
)

func init() {
	executor.RegisterFloat2("atan2", math.Atan2)
	executor.RegisterFloat2("copySign", math.Copysign)
	executor.RegisterFloat2("dim", math.Dim)
	executor.RegisterFloat2("hypot", math.Hypot)
	executor.RegisterFloat2("max", math.Max)
	executor.RegisterFloat2("min", math.Min)
	executor.RegisterFloat2("mod", math.Mod)
	executor.RegisterFloat2("nextAfter", math.Nextafter)
	executor.RegisterFloat2("pow", math.Pow)
	executor.RegisterFloat2("remainder", math.Remainder)

	executor.RegisterFloat1("abs", math.Abs)
	executor.RegisterFloat1("acos", math.Acos)
	executor.RegisterFloat1("acosh", math.Acosh)
	executor.RegisterFloat1("asin", math.Asin)
	executor.RegisterFloat1("asinh", math.Asinh)
	executor.RegisterFloat1("atan", math.Atan)
	executor.RegisterFloat1("atanh", math.Atanh)
	executor.RegisterFloat1("cbrt", math.Cbrt)
	executor.RegisterFloat1("ceil", math.Ceil)
	executor.RegisterFloat1("cos", math.Cos)
	executor.RegisterFloat1("cosh", math.Cosh)
	executor.RegisterFloat1("erf", math.Erf)
	executor.RegisterFloat1("erfc", math.Erfc)
	executor.RegisterFloat1("erfcinv", math.Erfcinv)
	executor.RegisterFloat1("exp", math.Exp)
	executor.RegisterFloat1("exp2", math.Exp2)
	executor.RegisterFloat1("expm1", math.Expm1)
	executor.RegisterFloat1("floor", math.Floor)
	executor.RegisterFloat1("gamma", math.Gamma)
	executor.RegisterFloat1("j0", math.J0)
	executor.RegisterFloat1("j1", math.J1)
	executor.RegisterFloat1("log", math.Log)
	executor.RegisterFloat1("logb", math.Logb)
	executor.RegisterFloat1("log1p", math.Log1p)
	executor.RegisterFloat1("log2", math.Log2)
	executor.RegisterFloat1("log10", math.Log10)
	executor.RegisterFloat1("round", math.Round)
	executor.RegisterFloat1("round2even", math.RoundToEven)
	executor.RegisterFloat1("sin", math.Sin)
	executor.RegisterFloat1("sinh", math.Sinh)
	executor.RegisterFloat1("sqrt", math.Sqrt)
	executor.RegisterFloat1("tan", math.Tan)
	executor.RegisterFloat1("tanh", math.Tanh)
	executor.RegisterFloat1("trunc", math.Trunc)
	executor.RegisterFloat1("y0", math.Y0)
	executor.RegisterFloat1("y1", math.Y1)
}
