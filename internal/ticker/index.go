package ticker

import (
	"time"

	"github.com/shopspring/decimal"
)

type Index struct {
	start     time.Time
	every     time.Duration
	delay     time.Duration
	precision int32
	prevPrice decimal.Decimal

	from   time.Time
	wa     []wa
	cursor int
}

func NewIndex(start time.Time, every time.Duration, delay time.Duration, precision int32, prevPrice decimal.Decimal) *Index {
	return &Index{start: start, every: every, delay: delay, precision: precision, prevPrice: prevPrice}
}

func (idx *Index) Init() {
	idx.from = idx.start.Truncate(idx.every)
	// if `delay` is greater than `every` additional space is allocated
	idx.wa = make([]wa, 2+idx.delay/idx.every)
}

func (idx *Index) Update(p Price) bool {
	offset := int(p.Time.Sub(idx.from) / idx.every)
	if offset < 0 || offset >= len(idx.wa) {
		return false
	}
	idx.wa[(idx.cursor+offset)%len(idx.wa)].update(p.Price, p.Volume)
	return true
}

func (idx *Index) TickAt() (time.Time, time.Duration) {
	return idx.from.Add(idx.delay), idx.every
}

func (idx *Index) Tick() Price {
	tick := Price{}
	tick.Time = idx.from.Add(idx.every)
	cur := &idx.wa[idx.cursor]
	tick.Price = cur.value(idx.precision)
	if tick.Price.IsZero() {
		tick.Price = idx.prevPrice
	}
	tick.Volume = cur.w

	*cur = wa{}
	idx.cursor = (idx.cursor + 1) % len(idx.wa)

	idx.from = tick.Time
	idx.prevPrice = tick.Price

	return tick
}

// weighted average
type wa struct {
	wv decimal.Decimal
	w  decimal.Decimal
}

func (a *wa) update(v decimal.Decimal, w decimal.Decimal) {
	a.wv = a.wv.Add(w.Mul(v))
	a.w = a.w.Add(w)
}

func (a *wa) value(prec int32) decimal.Decimal {
	if a.w.IsZero() {
		return decimal.Decimal{}
	}
	return a.wv.DivRound(a.w, prec)
}
