package aggregate

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"polygon-websocket-aggregator/model/trade"
	"sync"
	"time"
)

type Aggregate struct {
	Symbol                string  `json:"sym"`
	OpenPrice             float64 `json:""`
	OpenPriceTimestamp    int64
	ClosingPrice          float64 `json:""`
	ClosingPriceTimestamp int64
	HighPrice             float64 `json:""`
	LowPrice              float64 `json:""`
	Volume                int     `json:"v"`
	Timestamp             int64   `json:""`
	MutexLock             *sync.Mutex
}

func (agg *Aggregate) Print() {
	unixTimestamp := time.UnixMilli(agg.Timestamp)
	fmt.Printf("%d:%d:%.2d - open: $%.2f, close: $%.2f, high: $%.2f, low: $%.2f, volume: %d\n",
		unixTimestamp.Hour(), unixTimestamp.Minute(), unixTimestamp.Second(), agg.OpenPrice, agg.ClosingPrice, agg.HighPrice, agg.LowPrice, agg.Volume)
}

func (agg *Aggregate) DebugAggregate() {
	unixTimestamp := time.Unix(agg.Timestamp, 0)
	logrus.Debugf("%d:%d:%.2d - open: $%.2f, close: $%.2f, high: $%.2f, low: $%.2f, volume: %d\n",
		unixTimestamp.Hour(), unixTimestamp.Minute(), unixTimestamp.Second(), agg.OpenPrice, agg.ClosingPrice, agg.HighPrice, agg.LowPrice, agg.Volume)
}

func (agg *Aggregate) Update(trade trade.TradeRequest) {
	agg.MutexLock.Lock()
	defer agg.MutexLock.Unlock()

	agg.Volume += trade.Size
	if trade.Timestamp < agg.OpenPriceTimestamp {
		agg.OpenPrice = trade.Price
		agg.OpenPriceTimestamp = trade.Timestamp
	}
	if trade.Timestamp > agg.ClosingPriceTimestamp {
		agg.ClosingPrice = trade.Price
		agg.ClosingPriceTimestamp = trade.Timestamp
	}
	if trade.Price > agg.HighPrice {
		agg.HighPrice = trade.Price
	}
	if trade.Price < agg.LowPrice {
		agg.LowPrice = trade.Price
	}
}

func Calculate(trades []trade.TradeRequest, tickerName string, timeStamp int64) Aggregate {
	agg := Aggregate{Symbol: tickerName, Timestamp: timeStamp, MutexLock: &sync.Mutex{}}
	var highestPrice float64
	var lowestPrice float64
	if len(trades) == 0 {
		// all other values will be initialized to 0 due to the way Golang creates primal types
		return agg
	} else {
		agg.OpenPriceTimestamp = trades[0].Timestamp
		agg.OpenPrice = trades[0].Price
		highestPrice = math.SmallestNonzeroFloat64
		lowestPrice = math.MaxFloat64
	}
	totalVolume := 0
	for _, t := range trades {
		totalVolume += t.Size
		if t.Price > highestPrice {
			highestPrice = t.Price
		}
		if t.Price < lowestPrice {
			lowestPrice = t.Price
		}
		if t.Timestamp < agg.OpenPriceTimestamp {
			agg.OpenPrice = t.Price
			agg.OpenPriceTimestamp = t.Timestamp
		}
		if t.Timestamp > agg.ClosingPriceTimestamp {
			agg.ClosingPrice = t.Price
			agg.ClosingPriceTimestamp = t.Timestamp
		}
	}
	agg.LowPrice = lowestPrice
	agg.HighPrice = highestPrice
	agg.Volume = totalVolume
	return agg
}
