package process

import (
	"github.com/eshch/index_price/internal/generator"
)

type ValidateFunc func(value float64) bool

type Option func(g *Generator)

type Generator struct {
	last         float64
	increment    generator.Generator
	validateFunc ValidateFunc
}

func New(last float64, increment generator.Generator, opts ...Option) *Generator {
	g := &Generator{last: last, increment: increment}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func Between(lower float64, upper float64) Option {
	return func(g *Generator) {
		g.validateFunc = func(value float64) bool {
			return value >= lower && value <= upper
		}
	}
}

func (g *Generator) Next() float64 {
	next := g.last + g.increment.Next()
	if g.validateFunc == nil || g.validateFunc(next) {
		g.last = next
	}
	return g.last
}
