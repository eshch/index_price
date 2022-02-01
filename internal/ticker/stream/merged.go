package stream

import (
	"context"
	"sync"

	"github.com/eshch/index_price/internal/ticker"
)

type Stream interface {
	Start(ctx context.Context) <-chan ticker.Price
}

type MergedStream struct {
	streams []Stream

	out <-chan ticker.Price
}

func NewMergedStream(streams []Stream) *MergedStream {
	return &MergedStream{streams: streams}
}

func (m *MergedStream) Start(ctx context.Context) <-chan ticker.Price {
	if m.out != nil {
		return m.out
	}
	c := make(chan ticker.Price)
	m.out = c
	wg := &sync.WaitGroup{}
	for _, stream := range m.streams {
		wg.Add(1)
		in := stream.Start(ctx)
		go func(out chan<- ticker.Price) {
			defer wg.Done()
			for p := range in {
				out <- p
			}
		}(c)
	}
	go func(out chan<- ticker.Price) {
		defer close(c)
		wg.Wait()
	}(c)
	return m.out
}
