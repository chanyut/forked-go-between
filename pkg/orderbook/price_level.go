package orderbook

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

var _ json.Marshaler = (*PriceLevel)(nil)
var _ json.Unmarshaler = (*PriceLevel)(nil)

type PriceLevel struct {
	price    decimal.Decimal
	quantity decimal.Decimal
}

func NewPriceLevel(price, quantity decimal.Decimal) *PriceLevel {
	return &PriceLevel{price, quantity}
}

func (p *PriceLevel) Price() decimal.Decimal {
	return p.price
}

func (p *PriceLevel) Quantity() decimal.Decimal {
	return p.quantity
}

func (p *PriceLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Quantity decimal.Decimal `json:"quantity"`
			Price    decimal.Decimal `json:"price"`
		}{
			p.quantity,
			p.price,
		},
	)
}

func (p *PriceLevel) UnmarshalJSON(data []byte) error {
	obj := struct {
		Quantity decimal.Decimal `json:"quantity"`
		Price    decimal.Decimal `json:"price"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	p.quantity = obj.Quantity
	p.price = obj.Price

	return nil
}
