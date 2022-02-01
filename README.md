The main goal of this demo is to demonstrate some design principles rather than comprehensive business logic making a lot of sense.

# Some concepts used
- simple constructors
- random generators for testing
- pipelining via channels
- dependency injection
- graceful shutdown
- option pattern
- ring buffer

# Amendments to the task definition
## No errs channel in subscriptions
As it's supposed to be a service working online it does not make a lot of sense to me to expose errors in price streams. It's better to handle errors I can think of internally exposing just data stream. It may be different in more comprehensive cases.
## No use for ordering
While it may seem tempting to leverage the semantics of the streams implying events going with the time order, I don't really find a way to utilize it productively. The only possible way I could think of was stopping reading streams outrunning others, but it could lead to load spikes and buffer growing, and doesn't really give any advantages besides an interesting task to orchestrate all of this.
## Delay for index price
As data naturally may be late, it makes sense to publish index price with a delay allowing late data to be taken into account.
## Weighted average
To get a fair price it makes sense to use volume-weighted average price, so I introduced volumes in price stream.
## No ticker symbol
I didn't really find a way how to utilize symbol in this demo, so I dropped it.
## Previous price for no data
In the case of no data arrived in a particular period, index publishes last price it could calculate.

# Other clarifications
## Performance
From data I found on daily transaction numbers for crypto exchanges it doesn't look like we should expect a high rates which the current implementation can't handle according my checks with higher load. 
## Random generators
To kind of utilize fuzzy testing and not to make up test data, I implemented and used generators for a few distributions modelling realistically looking values. During the service initialization concrete generators (both random and not) are chosen randomly as well as their parameters to build a hundred of different sources of data
## Missed data
### Input
If the input is late for more than the delay we set it is dropped. If the input outruns index object processing a lot, it may lead to some drops, but shouldn't be a lot.
### Output
If the consumer of index price stream is slow, not consumed tick is overridden by the new one. It's not that hard to avoid, so I'd rather showed the concept of dropping.
### Period for data aggregation
is a few seconds for testing convenience, the timestamps of the index ticks are rounded to the duration multiplier, so it would be solid minutes in case of minute duration.

# The runnable demo
`go run cmd/print/main.go`
- from 100 random generators
- make a single stream
- consumed by the index object
- which periodically streams down the resulting value
- to the consumer printing the data to stdout
