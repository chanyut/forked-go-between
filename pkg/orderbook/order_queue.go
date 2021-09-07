package orderbook

import (
	"container/list"

	"github.com/shopspring/decimal"
)

type OrderQueue struct {
	quantity decimal.Decimal
	price    decimal.Decimal
	orders   *list.List
}

func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{decimal.Zero, price, list.New()}
}

func (oq *OrderQueue) Price() decimal.Decimal {
	return oq.price
}

func (oq *OrderQueue) Quantity() decimal.Decimal {
	return oq.quantity
}

func (oq *OrderQueue) Orders() *list.List {
	return oq.orders
}

func (oq *OrderQueue) Len() int {
	return oq.orders.Len()
}

func (oq *OrderQueue) Head() *list.Element {
	return oq.orders.Front()
}

func (oq *OrderQueue) Tail() *list.Element {
	return oq.orders.Back()
}

func (oq *OrderQueue) Append(order *Order) *list.Element {
	oq.quantity = oq.quantity.Add(order.quantity)
	return oq.orders.PushBack(order)
}

func (oq *OrderQueue) Update(e *list.Element, quantity decimal.Decimal) *Order {
	order := e.Value.(*Order)
	oq.quantity = oq.quantity.Sub(order.quantity)
	oq.quantity = oq.quantity.Add(quantity)
	order.quantity = quantity
	return order
}

func (oq *OrderQueue) Remove(e *list.Element) *Order {
	oq.quantity = oq.quantity.Sub(e.Value.(*Order).quantity)
	return oq.orders.Remove(e).(*Order)
}
