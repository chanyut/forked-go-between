package orderbook

import "github.com/rotisserie/eris"

var (
	ErrInvalidOrderID     = eris.New("Invalid order id")
	ErrInvalidTraderID    = eris.New("Invalid trader id")
	ErrInvalidMarket      = eris.New("Invalid market")
	ErrInvalidQuantity    = eris.New("Invalid quantity")
	ErrInvalidPrice       = eris.New("Invalid price")
	ErrInvalidSide        = eris.New("Invalid side")
	ErrOrderAlreadyExists = eris.New("Order already exists")
)
