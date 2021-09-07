package orderbook

import (
	"container/list"
	"encoding/json"
	"sync"

	"github.com/rotisserie/eris"
)

var _ json.Marshaler = (*OrderBook)(nil)
var _ json.Unmarshaler = (*OrderBook)(nil)

type OrderBook struct {
	sync.RWMutex
	market string
	orders map[string]*list.Element
	asks   *OrderSide
	bids   *OrderSide
}

func NewOrderBook(market string) *OrderBook {
	return &OrderBook{sync.RWMutex{}, market, make(map[string]*list.Element), NewOrderSide(), NewOrderSide()}
}

func (ob *OrderBook) restoreAsks(orders []*Order) {
	for _, order := range orders {
		ob.orders[order.id] = ob.asks.Append(order)
	}
}

func (ob *OrderBook) restoreBids(orders []*Order) {
	for _, order := range orders {
		ob.orders[order.id] = ob.bids.Append(order)
	}
}

func (ob *OrderBook) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Market string   `json:"market"`
			Bids   []*Order `json:"bids"`
			Asks   []*Order `json:"asks"`
		}{
			ob.market,
			ob.bids.Orders(),
			ob.asks.Orders(),
		},
	)
}

func (ob *OrderBook) UnmarshalJSON(data []byte) error {
	obj := struct {
		Market string   `json:"market"`
		Bids   []*Order `json:"bids"`
		Asks   []*Order `json:"asks"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	ob.market = obj.Market
	ob.bids = NewOrderSide()
	ob.asks = NewOrderSide()

	ob.orders = make(map[string]*list.Element)
	ob.restoreBids(obj.Bids)
	ob.restoreAsks(obj.Asks)

	return nil
}
