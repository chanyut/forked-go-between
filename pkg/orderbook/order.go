package orderbook

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

var _ json.Marshaler = (*Order)(nil)
var _ json.Unmarshaler = (*Order)(nil)

type Order struct {
	id       string
	traderID string
	side     Side
	quantity decimal.Decimal
	price    decimal.Decimal
}

func NewOrder(ID, traderID string, side Side, quantity, price decimal.Decimal) *Order {
	return &Order{ID, traderID, side, quantity, price}
}

func (o *Order) ID() string {
	return o.id
}

func (o *Order) TraderID() string {
	return o.traderID
}

func (o *Order) Side() Side {
	return o.side
}

func (o *Order) Quantity() decimal.Decimal {
	return o.quantity
}

func (o *Order) Price() decimal.Decimal {
	return o.price
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			ID       string          `json:"id"`
			TraderID string          `json:"traderId"`
			Side     Side            `json:"side"`
			Quantity decimal.Decimal `json:"quantity"`
			Price    decimal.Decimal `json:"price"`
		}{
			o.id,
			o.traderID,
			o.side,
			o.quantity,
			o.price,
		},
	)
}

func (o *Order) UnmarshalJSON(data []byte) error {
	obj := struct {
		ID       string          `json:"id"`
		TraderID string          `json:"traderId"`
		Side     Side            `json:"side"`
		Quantity decimal.Decimal `json:"quantity"`
		Price    decimal.Decimal `json:"price"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	o.id = obj.ID
	o.side = obj.Side
	o.traderID = obj.TraderID
	o.side = obj.Side
	o.quantity = obj.Quantity
	o.price = obj.Price

	return nil
}
