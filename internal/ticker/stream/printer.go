package stream

import (
	"context"
	"fmt"
	"time"

	"github.com/eshch/index_price/internal/ticker"
)

type StringFunc func(price ticker.Price) string

func GoStructString(price ticker.Price) string {
	return fmt.Sprintf("%#v", price)
}

func UTCPrice(price ticker.Price) string {
	return fmt.Sprintf(
		"%s %s",
		price.Time.In(time.UTC).Format("2006-01-02T15:04:05.999999"),
		price.Price.StringFixed(4),
	)
}

func UTCMicroPriceVolume(price ticker.Price) string {
	return fmt.Sprintf(
		"%s %s %s",
		price.Time.In(time.UTC).Format("2006-01-02T15:04:05.000000"),
		price.Price.StringFixed(4),
		price.Volume.StringFixed(2),
	)
}

type Printer struct {
	stream     Stream
	stringFunc StringFunc
}

func NewPrinter(stream Stream, stringFunc StringFunc) *Printer {
	return &Printer{stream: stream, stringFunc: stringFunc}
}

func (p *Printer) Run(ctx context.Context) {
	for price := range p.stream.Start(ctx) {
		fmt.Println(p.stringFunc(price))
	}
}
