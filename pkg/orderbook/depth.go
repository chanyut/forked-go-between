package orderbook

import (
	"encoding/json"

	"github.com/rotisserie/eris"
)

var _ json.Marshaler = (*Depth)(nil)
var _ json.Unmarshaler = (*Depth)(nil)

type Depth struct {
	bids []*PriceLevel
	asks []*PriceLevel
}

func NewDepth(bids, asks []*PriceLevel) *Depth {
	return &Depth{bids, asks}
}

func (d *Depth) Bids() []*PriceLevel {
	return d.bids
}

func (d *Depth) Asks() []*PriceLevel {
	return d.asks
}

func (d *Depth) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Bids []*PriceLevel `json:"bids"`
			Asks []*PriceLevel `json:"asks"`
		}{
			d.bids,
			d.asks,
		},
	)
}

func (d *Depth) UnmarshalJSON(data []byte) error {
	obj := struct {
		Bids []*PriceLevel `json:"bids"`
		Asks []*PriceLevel `json:"asks"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	d.bids = obj.Bids
	d.asks = obj.Asks

	return nil
}
