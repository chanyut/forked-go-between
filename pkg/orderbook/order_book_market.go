package orderbook

import (
	"strings"

	"github.com/shopspring/decimal"
)

func (ob *OrderBook) ProcessMarketOrder(orderID, traderID, market string, side Side, quantity, price decimal.Decimal) ([]*Trade, error) {
	defer ob.Unlock()
	ob.Lock()

	if strings.TrimSpace(orderID) == "" {
		return nil, ErrInvalidOrderID
	}

	if ob.orders[orderID] != nil {
		return nil, ErrOrderAlreadyExists
	}

	if strings.TrimSpace(traderID) == "" {
		return nil, ErrInvalidTraderID
	}

	if market != ob.market {
		return nil, ErrInvalidMarket
	}

	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidQuantity
	}

	if price.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidPrice
	}

	var (
		sideToProcess *OrderSide
		level         *OrderQueue
		next          func(decimal.Decimal) *OrderQueue
	)

	if side == Buy {
		sideToProcess = ob.asks
		level = ob.asks.MinPriceQueue()
		next = ob.asks.GreaterThan
	} else {
		sideToProcess = ob.bids
		level = ob.bids.MaxPriceQueue()
		next = ob.bids.LessThan
	}

	quantityToTrade := quantity
	priceToTrade := price
	trades := make([]*Trade, 0)

	for level != nil && quantityToTrade.GreaterThan(decimal.Zero) && priceToTrade.GreaterThan(decimal.Zero) {
		headOrderEl := level.Head()

		for headOrderEl != nil && quantityToTrade.GreaterThan(decimal.Zero) && priceToTrade.GreaterThan(decimal.Zero) {
			headOrder := headOrderEl.Value.(*Order)

			if headOrder.traderID == traderID {
				headOrderEl = headOrderEl.Next()
				continue
			}

			if quantity.GreaterThanOrEqual(headOrder.quantity) {
				trades = append(trades, NewTrade(orderID, headOrder.id, headOrder.quantity, headOrder.price))
				quantityToTrade = quantityToTrade.Sub(headOrder.quantity)
				priceToTrade = priceToTrade.Sub(headOrder.price)

				headOrderEl = headOrderEl.Next()
				ob.remove(headOrder.id)
			} else {
				trades = append(trades, NewTrade(orderID, headOrder.id, quantityToTrade, headOrder.price))
				sideToProcess.Update(headOrderEl, headOrder.quantity.Sub(quantityToTrade))
				quantityToTrade = decimal.Zero
				priceToTrade = decimal.Zero

				headOrderEl = headOrderEl.Next()
			}
		}

		level = next(level.price)
	}

	return trades, nil
}
