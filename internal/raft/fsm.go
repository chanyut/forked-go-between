package raft

import (
	"encoding/json"
	"io"

	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/hashicorp/raft"
	"github.com/rotisserie/eris"
)

var _ raft.FSM = (*StateMachine)(nil)

type StateMachine struct {
	book *orderbook.OrderBook
}

func NewStateMachine(book *orderbook.OrderBook) *StateMachine {
	return &StateMachine{book}
}

func (s *StateMachine) Apply(log *raft.Log) interface{} {
	var cmd RaftCmd

	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		return eris.Wrap(err, "json.Unmarshal")
	}

	switch cmd.Op {
	case ProcessLimitOrderOp:
		result, err := s.book.ProcessLimitOrder(
			cmd.OrderID,
			cmd.TraderID,
			cmd.Market,
			cmd.Side,
			cmd.Quantity,
			cmd.Price,
		)

		if err != nil {
			return eris.Wrap(err, "s.book.ProcessLimitOrde")
		}

		return result
	case ProcessMarketOrderOp:
		result, err := s.book.ProcessMarketOrder(
			cmd.OrderID,
			cmd.TraderID,
			cmd.Market,
			cmd.Side,
			cmd.Quantity,
			cmd.Price,
		)

		if err != nil {
			return eris.Wrap(err, "s.book.ProcessMarketOrder")
		}

		return result
	case CancelOrderOp:
		return s.book.CancelOrder(cmd.OrderID)
	default:
		return eris.Wrap(ErrCmdNotFound, "switch default")
	}
}

func (s *StateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{s.book}, nil
}

func (s *StateMachine) Restore(rc io.ReadCloser) error {
	defer rc.Close()

	var book orderbook.OrderBook
	err := json.NewDecoder(rc).Decode(&book)
	if err != nil {
		return eris.Wrap(err, "dec.Decode")
	}

	s.book = &book
	return nil
}
