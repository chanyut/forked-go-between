package orderbook

import (
	"container/list"
	"sort"

	rbtx "github.com/emirpasic/gods/examples/redblacktreeextended"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type OrderSide struct {
	priceTree *rbtx.RedBlackTreeExtended
	prices    map[string]*OrderQueue
	quantity  decimal.Decimal
	size      int
	depth     int
}

func NewOrderSide() *OrderSide {
	priceTree := &rbtx.RedBlackTreeExtended{
		Tree: rbt.NewWith(func(a, b interface{}) int {
			return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
		}),
	}

	return &OrderSide{priceTree, make(map[string]*OrderQueue), decimal.Zero, 0, 0}
}

func (os *OrderSide) Append(order *Order) *list.Element {
	price := order.price
	strPrice := price.String()

	priceQueue, ok := os.prices[strPrice]
	if !ok {
		priceQueue = NewOrderQueue(price)
		os.prices[strPrice] = priceQueue
		os.priceTree.Put(price, priceQueue)
		os.depth++
	}

	os.size++
	os.quantity = os.quantity.Add(order.quantity)
	return priceQueue.Append(order)
}

func (os *OrderSide) Remove(e *list.Element) *Order {
	order := e.Value.(*Order)

	price := order.price
	strPrice := price.String()

	priceQueue := os.prices[strPrice]
	o := priceQueue.Remove(e)

	if priceQueue.Len() == 0 {
		delete(os.prices, strPrice)
		os.priceTree.Remove(price)
		os.depth--
	}

	os.size--
	os.quantity = os.quantity.Sub(o.Quantity())
	return o
}

func (os *OrderSide) Update(e *list.Element, quantity decimal.Decimal) *Order {
	order := e.Value.(*Order)

	price := order.price
	strPrice := price.String()

	os.quantity = os.quantity.Sub(order.quantity)
	os.quantity = os.quantity.Add(quantity)

	priceQueue := os.prices[strPrice]
	o := priceQueue.Update(e, quantity)

	return o
}

func (os *OrderSide) MaxPriceQueue() *OrderQueue {
	if os.depth <= 0 {
		return nil
	}

	if value, found := os.priceTree.GetMax(); found {
		return value.(*OrderQueue)
	}

	return nil
}

func (os *OrderSide) MinPriceQueue() *OrderQueue {
	if os.depth <= 0 {
		return nil
	}

	if value, found := os.priceTree.GetMin(); found {
		return value.(*OrderQueue)
	}

	return nil
}

func (os *OrderSide) LessThan(price decimal.Decimal) *OrderQueue {
	tree := os.priceTree.Tree
	node := tree.Root

	var floor *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) > 0 {
			floor = node
			node = node.Right
		} else {
			node = node.Left
		}
	}

	if floor != nil {
		return floor.Value.(*OrderQueue)
	}

	return nil
}

func (os *OrderSide) GreaterThan(price decimal.Decimal) *OrderQueue {
	tree := os.priceTree.Tree
	node := tree.Root

	var ceiling *rbt.Node
	for node != nil {
		if tree.Comparator(price, node.Key) < 0 {
			ceiling = node
			node = node.Left
		} else {
			node = node.Right
		}
	}

	if ceiling != nil {
		return ceiling.Value.(*OrderQueue)
	}

	return nil
}

func (os *OrderSide) Orders() []*Order {
	orders := make([]*Order, 0)

	for _, price := range os.prices {
		iter := price.Head()

		for iter != nil {
			orders = append(orders, iter.Value.(*Order))
			iter = iter.Next()
		}
	}

	sort.Slice(orders[:], func(i, j int) bool {
		if orders[i].side == Buy {
			return orders[i].price.GreaterThan(orders[j].price)
		} else {
			return orders[i].price.LessThan(orders[j].price)
		}
	})

	return orders
}
