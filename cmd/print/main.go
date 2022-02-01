package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/rand"

	"github.com/eshch/index_price/internal/generator"
	"github.com/eshch/index_price/internal/generator/process"
	"github.com/eshch/index_price/internal/generator/types"
	"github.com/eshch/index_price/internal/ticker"
	"github.com/eshch/index_price/internal/ticker/stream"
)

func makeGeneratedPriceStream(
	startTime time.Time,
	timeTick generator.Generator,
	startPrice float64,
	priceStep generator.Generator,
	lowerPrice float64,
	upperPrice float64,
	volume generator.Generator,
	delay generator.Generator,
) stream.Stream {
	return stream.NewGeneratedStream(
		time.Now,
		types.NewTime(process.New(float64(startTime.Unix()), timeTick)),
		types.NewDecimal(process.New(startPrice, priceStep, process.Between(lowerPrice, upperPrice)), 2),
		types.NewDecimal(volume, 8),
		types.NewDuration(delay),
	)
}

func makeGeneratedPriceStreams(n int) []stream.Stream {
	streams := []stream.Stream{}
	startTime := time.Now()
	for i := 0; i < n; i++ {
		seed := uint64(i)
		timeTicks := []generator.Generator{
			generator.NewConst(rand.Float64()),
			generator.NewPoisson(rand.NewSource(0), rand.Float64()),
			generator.NewConst(1),
		}
		priceSteps := []generator.Generator{
			generator.NewRandSlice(rand.NewSource(seed), []float64{-1, 0, 1}),
			generator.NewNorm(rand.NewSource(seed), 0, 1),
		}
		startPrice := 50 + rand.Float64()
		s := makeGeneratedPriceStream(
			startTime,
			timeTicks[rand.Intn(len(timeTicks))],
			startPrice,
			priceSteps[rand.Intn(len(priceSteps))],
			startPrice-10,
			startPrice+10,
			generator.NewErlang(rand.NewSource(seed), 1+rand.Intn(100), rand.Float64()),
			generator.NewErlang(rand.NewSource(seed), 1+rand.Intn(10), rand.Float64()*100),
		)
		streams = append(streams, s)
	}
	return streams
}

func main() {
	streams := makeGeneratedPriceStreams(100)
	mergedPriceStream := stream.NewMergedStream(streams)
	index := ticker.NewIndex(time.Now(), 2*time.Second, 1*time.Second, 4, decimal.Decimal{})
	index.Init()
	indexPriceStream := stream.NewIndex(mergedPriceStream, index, time.Now)
	priceStreamPrinter := stream.NewPrinter(indexPriceStream, stream.UTCPrice)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	priceStreamPrinter.Run(ctx)
}
