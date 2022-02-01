package ticker

import (
	"time"

	"github.com/shopspring/decimal"
)

type Price struct {
	Time   time.Time
	Price  decimal.Decimal
	Volume decimal.Decimal
}
