package orderbook

import (
	"strings"

	"github.com/shopspring/decimal"
)

func (ob *OrderBook) MarketPrice(traderID, market string, side Side, quantity decimal.Decimal) (*MarketPrice, error) {
	defer ob.RUnlock()
	ob.RLock()

	if strings.TrimSpace(traderID) == "" {
		return nil, ErrInvalidTraderID
	}

	if market != ob.market {
		return nil, ErrInvalidMarket
	}

	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	price := decimal.Zero

	var (
		level *OrderQueue
		next  func(decimal.Decimal) *OrderQueue
	)

	if side == Buy {
		level = ob.asks.MinPriceQueue()
		next = ob.asks.GreaterThan
	} else {
		level = ob.bids.MaxPriceQueue()
		next = ob.bids.LessThan
	}

	for level != nil && quantity.GreaterThan(decimal.Zero) {
		headOrderEl := level.Head()

		for headOrderEl != nil && quantity.GreaterThan(decimal.Zero) {
			headOrder := headOrderEl.Value.(*Order)

			if headOrder.traderID == traderID {
				headOrderEl = headOrderEl.Next()
				continue
			}

			if quantity.GreaterThanOrEqual(headOrder.quantity) {
				price = price.Add(headOrder.price.Mul(headOrder.quantity))
				quantity = quantity.Sub(headOrder.quantity)
			} else {
				price = price.Add(headOrder.price.Mul(quantity))
				quantity = decimal.Zero
			}

			headOrderEl = headOrderEl.Next()
		}

		level = next(level.price)
	}

	return NewMarketPrice(price, quantity), nil
}
