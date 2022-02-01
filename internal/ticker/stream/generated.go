package stream

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/eshch/index_price/internal/ticker"
)

type Time interface {
	Next() time.Time
}

type Duration interface {
	Next() time.Duration
}

type Decimal interface {
	Next() decimal.Decimal
}

type GeneratedStream struct {
	nowFunc NowFunc
	time    Time
	price   Decimal
	volume  Decimal
	delay   Duration

	out <-chan ticker.Price
}

func NewGeneratedStream(nowFunc NowFunc, time Time, price Decimal, volume Decimal, delay Duration) *GeneratedStream {
	return &GeneratedStream{nowFunc: nowFunc, time: time, price: price, volume: volume, delay: delay}
}

func (s *GeneratedStream) Start(ctx context.Context) <-chan ticker.Price {
	if s.out != nil {
		return s.out
	}
	c := make(chan ticker.Price)
	s.out = c
	go func(out chan<- ticker.Price) {
		defer close(out)
		for {
			priceTick := ticker.Price{
				Time:   s.time.Next(),
				Price:  s.price.Next(),
				Volume: s.volume.Next(),
			}
			select {
			case <-time.After(priceTick.Time.Sub(s.nowFunc()) + s.delay.Next()):
			case <-ctx.Done():
				return
			}
			select {
			case out <- priceTick:
			case <-ctx.Done():
				return
			}
		}
	}(c)
	return s.out
}
