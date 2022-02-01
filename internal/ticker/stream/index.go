package stream

import (
	"context"
	"sync"
	"time"

	"github.com/eshch/index_price/internal/ticker"
)

type NowFunc func() time.Time

type Index struct {
	price   Stream
	ticker  *ticker.Index
	nowFunc NowFunc

	out <-chan ticker.Price
}

func NewIndex(price Stream, ticker *ticker.Index, nowFunc NowFunc) *Index {
	return &Index{price: price, ticker: ticker, nowFunc: nowFunc}
}

func (idx *Index) timeTicks(ctx context.Context) <-chan struct{} {
	c := make(chan struct{})
	go func(out chan<- struct{}) {
		defer close(out)
		from, every := idx.ticker.TickAt()
		at := from
		for {
			at = at.Add(every)
			select {
			case <-time.After(at.Sub(idx.nowFunc())):
			case <-ctx.Done():
				return
			}

			select {
			case out <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}(c)
	return c
}

func (idx *Index) Start(ctx context.Context) <-chan ticker.Price {
	if idx.out != nil {
		return idx.out
	}
	c := make(chan ticker.Price)
	idx.out = c
	wg := &sync.WaitGroup{}
	wg.Add(1)
	in := idx.price.Start(ctx)
	wg.Add(1)
	timeTicks := idx.timeTicks(ctx)
	d := make(chan struct{})
	go func(done chan<- struct{}) {
		defer close(d)
		wg.Wait()
	}(d)
	go func(out chan<- ticker.Price, done <-chan struct{}) {
		defer close(out)
		var priceTicks chan<- ticker.Price
		var priceTick ticker.Price
		for {
			select {
			case p, ok := <-in:
				if !ok {
					in = nil
					wg.Done()
					break
				}
				idx.ticker.Update(p)
			case _, ok := <-timeTicks:
				if !ok {
					timeTicks = nil
					wg.Done()
					break
				}
				if priceTicks == nil {
					wg.Add(1)
					priceTicks = out
				}
				priceTick = idx.ticker.Tick()
			case priceTicks <- priceTick:
				priceTicks = nil
				wg.Done()
			case <-done:
				return
			}
		}
	}(c, d)
	return idx.out
}
