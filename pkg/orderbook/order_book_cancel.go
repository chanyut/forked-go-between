package orderbook

func (ob *OrderBook) CancelOrder(orderID string) *Order {
	defer ob.Unlock()
	ob.Lock()

	return ob.remove(orderID)
}

func (ob *OrderBook) remove(orderID string) *Order {
	e, ok := ob.orders[orderID]
	if !ok {
		return nil
	}

	delete(ob.orders, orderID)

	if e.Value.(*Order).side == Buy {
		return ob.bids.Remove(e)
	}

	return ob.asks.Remove(e)
}
