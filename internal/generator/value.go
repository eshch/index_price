package generator

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Generator interface {
	Next() float64
}

type Func func() float64

func (f Func) Next() float64 {
	return f()
}

func NewConst(value float64) Generator {
	return Func(func() float64 { return value })
}

func NewUniform(source rand.Source) Generator {
	rnd := rand.New(source)
	return Func(func() float64 { return rnd.Float64() })
}

func NewRandSlice(source rand.Source, values []float64) Generator {
	rnd := rand.New(source)
	return Func(func() float64 { return values[rnd.Intn(len(values))] })
}

func NewNorm(source rand.Source, mean float64, stddev float64) Generator {
	rnd := rand.New(source)
	return Func(func() float64 { return mean + stddev*rnd.NormFloat64() })
}

func NewGamma(source rand.Source, k float64, rate float64) Generator {
	return Func(distuv.Gamma{Src: source, Alpha: k, Beta: rate}.Rand)
}

func NewErlang(source rand.Source, k int, rate float64) Generator {
	return NewGamma(source, float64(k), rate)
}

func NewPoisson(source rand.Source, rate float64) Generator {
	return NewErlang(source, 1, rate)
}
