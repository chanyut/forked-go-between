package orderbook

import (
	"encoding/json"
	"reflect"
)

var _ json.Marshaler = (*Side)(nil)
var _ json.Unmarshaler = (*Side)(nil)

type Side int

const (
	Sell Side = 0
	Buy  Side = 1
)

func ParseSide(side string) (Side, error) {
	if side == "buy" {
		return Buy, nil
	}

	if side == "Sell" {
		return Sell, nil
	}

	return -1, ErrInvalidSide
}

func (s Side) String() string {
	if s == Buy {
		return "buy"
	}

	return "sell"
}

func (s Side) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *Side) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"buy"`:
		*s = Buy
	case `"sell"`:
		*s = Sell
	default:
		return &json.UnsupportedValueError{
			Value: reflect.New(reflect.TypeOf(data)),
			Str:   string(data),
		}
	}

	return nil
}
