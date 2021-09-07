package orderbook

import (
	"strings"

	"github.com/shopspring/decimal"
)

func (ob *OrderBook) ProcessLimitOrder(orderID, traderID, market string, side Side, quantity, price decimal.Decimal) ([]*Trade, error) {
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
		sideToAdd     *OrderSide
		sideToProcess *OrderSide
		comparator    func(decimal.Decimal) bool
		best          func() *OrderQueue
		next          func(decimal.Decimal) *OrderQueue
	)

	if side == Buy {
		sideToAdd = ob.bids
		sideToProcess = ob.asks
		comparator = price.GreaterThanOrEqual
		best = ob.asks.MinPriceQueue
		next = ob.asks.GreaterThan
	} else {
		sideToAdd = ob.asks
		sideToProcess = ob.bids
		comparator = price.LessThanOrEqual
		best = ob.bids.MaxPriceQueue
		next = ob.bids.LessThan
	}

	trades := make([]*Trade, 0)
	quantityToTrade := quantity
	bestPrice := best()

	for bestPrice != nil && quantityToTrade.GreaterThan(decimal.Zero) && comparator(bestPrice.price) {
		headOrderEl := bestPrice.Head()
		bestPrice = next(bestPrice.price)

		for headOrderEl != nil && quantityToTrade.GreaterThan(decimal.Zero) {
			headOrder := headOrderEl.Value.(*Order)

			if headOrder.traderID == traderID {
				headOrderEl = headOrderEl.Next()
				continue
			}

			if quantityToTrade.GreaterThanOrEqual(headOrder.quantity) {
				trades = append(trades, NewTrade(orderID, headOrder.id, headOrder.quantity, headOrder.price))
				quantityToTrade = quantityToTrade.Sub(headOrder.quantity)

				headOrderEl = headOrderEl.Next()
				ob.remove(headOrder.id)
			} else {
				trades = append(trades, NewTrade(orderID, headOrder.id, quantityToTrade, headOrder.price))
				sideToProcess.Update(headOrderEl, headOrder.quantity.Sub(quantityToTrade))
				quantityToTrade = decimal.Zero

				headOrderEl = headOrderEl.Next()
			}
		}
	}

	if quantityToTrade.GreaterThan(decimal.Zero) {
		order := NewOrder(orderID, traderID, side, quantityToTrade, price)
		ob.orders[order.id] = sideToAdd.Append(order)
	}

	return trades, nil
}
