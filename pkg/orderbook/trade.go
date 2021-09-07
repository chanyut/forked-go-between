package orderbook

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

var _ json.Marshaler = (*Trade)(nil)
var _ json.Unmarshaler = (*Trade)(nil)

type Trade struct {
	takerOrderID string
	makerOrderID string
	quantity     decimal.Decimal
	price        decimal.Decimal
}

func (t *Trade) TakerOrderID() string {
	return t.takerOrderID
}

func (t *Trade) MakerOrderID() string {
	return t.makerOrderID
}

func (t *Trade) Quantity() decimal.Decimal {
	return t.quantity
}

func (t *Trade) Price() decimal.Decimal {
	return t.price
}

func NewTrade(takerOrderID, makerOrderID string, quantity, price decimal.Decimal) *Trade {
	return &Trade{takerOrderID, makerOrderID, quantity, price}
}

func (t *Trade) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			TakerOrderID string          `json:"takerOrderId"`
			MakerOrderID string          `json:"makerOrderId"`
			Quantity     decimal.Decimal `json:"quantity"`
			Price        decimal.Decimal `json:"price"`
		}{
			t.takerOrderID,
			t.makerOrderID,
			t.quantity,
			t.price,
		},
	)
}

func (t *Trade) UnmarshalJSON(data []byte) error {
	obj := struct {
		TakerOrderID string          `json:"takerOrderId"`
		MakerOrderID string          `json:"makerOrderId"`
		Quantity     decimal.Decimal `json:"quantity"`
		Price        decimal.Decimal `json:"price"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	t.takerOrderID = obj.TakerOrderID
	t.makerOrderID = obj.MakerOrderID
	t.quantity = obj.Quantity
	t.price = obj.Price

	return nil
}
