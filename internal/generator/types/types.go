package types

import (
	"math"
	"time"

	"github.com/shopspring/decimal"

	"github.com/eshch/index_price/internal/generator"
)

type Decimal struct {
	float64   generator.Generator
	precision int32
}

func NewDecimal(float64 generator.Generator, precision int32) *Decimal {
	return &Decimal{float64: float64, precision: precision}
}

func (g *Decimal) Next() decimal.Decimal {
	return decimal.NewFromFloatWithExponent(g.float64.Next(), -g.precision)
}

type Time struct {
	float64 generator.Generator
}

func NewTime(float64 generator.Generator) *Time {
	return &Time{float64: float64}
}

func (g *Time) Next() time.Time {
	sec, dec := math.Modf(g.float64.Next())
	return time.Unix(int64(sec), int64(dec*1e9))
}

type Duration struct {
	float64 generator.Generator
}

func NewDuration(float64 generator.Generator) *Duration {
	return &Duration{float64: float64}
}

func (g *Duration) Next() time.Duration {
	return time.Duration(g.float64.Next()) * time.Second
}
