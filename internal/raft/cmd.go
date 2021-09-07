package raft

import (
	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/shopspring/decimal"
)

type RaftCmd struct {
	Op       RaftCmdOp       `json:"op"`
	OrderID  string          `json:"orderId"`
	TraderID string          `json:"traderId"`
	Market   string          `json:"market"`
	Side     orderbook.Side  `json:"side"`
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
}
