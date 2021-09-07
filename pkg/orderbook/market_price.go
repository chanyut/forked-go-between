package orderbook

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

var _ json.Marshaler = (*MarketPrice)(nil)
var _ json.Unmarshaler = (*MarketPrice)(nil)

type MarketPrice struct {
	price             decimal.Decimal
	remainingQuantity decimal.Decimal
}

func NewMarketPrice(price, remainingQuantity decimal.Decimal) *MarketPrice {
	return &MarketPrice{price, remainingQuantity}
}

func (m *MarketPrice) Price() decimal.Decimal {
	return m.price
}

func (m *MarketPrice) RemainingQuantity() decimal.Decimal {
	return m.remainingQuantity
}

func (m *MarketPrice) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Price             decimal.Decimal `json:"price"`
			RemainingQuantity decimal.Decimal `json:"remainingQuantity"`
		}{
			m.price,
			m.remainingQuantity,
		},
	)
}

func (m *MarketPrice) UnmarshalJSON(data []byte) error {
	obj := struct {
		Price             decimal.Decimal `json:"price"`
		RemainingQuantity decimal.Decimal `json:"remainingQuantity"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	m.price = obj.Price
	m.remainingQuantity = obj.RemainingQuantity

	return nil
}
