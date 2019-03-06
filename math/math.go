package math

import (
	"math"

	"github.com/shopspring/decimal"
)

func Round(value float64, precision int32) float64 {
	if math.IsNaN(value) {
		return value
	}

	dValue := decimal.NewFromFloat(value)
	out, _ := dValue.Round(precision).Float64()
	return out
}

func Floor(value float64, precision int32) float64 {
	dMulti := decimal.NewFromFloat(math.Pow10(int(precision)))
	dValue := decimal.NewFromFloat(value)
	out, _ := dValue.Mul(dMulti).Floor().Div(dMulti).Float64()
	return out
}

func Ceil(value float64, precision int32) float64 {
	dMulti := decimal.NewFromFloat(math.Pow10(int(precision)))
	dValue := decimal.NewFromFloat(value)
	out, _ := dValue.Mul(dMulti).Ceil().Div(dMulti).Float64()
	return out
}

func Trunc(v float64, precision int) float64 {
	return math.Trunc(v*math.Pow10(precision)+0.5) * math.Pow10(-precision)
}

func Pow10(v int) float64 {
	return math.Pow10(v)
}
