package raft

import (
	"encoding/json"
	"reflect"
)

type RaftCmdOp string

const (
	ProcessLimitOrderOp  RaftCmdOp = "ProcessLimitOrder"
	ProcessMarketOrderOp RaftCmdOp = "ProcessMarketOrder"
	CancelOrderOp        RaftCmdOp = "CancelOrder"
)

func ParseRaftCmdOp(op string) (RaftCmdOp, error) {
	if op == "ProcessLimitOrder" {
		return ProcessLimitOrderOp, nil
	}

	if op == "ProcessMarketOrder" {
		return ProcessMarketOrderOp, nil
	}

	if op == "CancelOrder" {
		return CancelOrderOp, nil
	}

	return "", ErrSyntaxParseCmdOp
}

func (op RaftCmdOp) String() string {
	if op == ProcessLimitOrderOp {
		return "ProcessLimitOrder"
	}

	if op == ProcessMarketOrderOp {
		return "ProcessMarketOrder"
	}

	if op == CancelOrderOp {
		return "CancelOrder"
	}

	return ""
}

func (op RaftCmdOp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + op.String() + `"`), nil
}

func (op *RaftCmdOp) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"ProcessLimitOrder"`:
		*op = ProcessLimitOrderOp
	case `"ProcessMarketOrder"`:
		*op = ProcessMarketOrderOp
	case `"CancelOrder"`:
		*op = CancelOrderOp
	default:
		return &json.UnsupportedValueError{
			Value: reflect.New(reflect.TypeOf(data)),
			Str:   string(data),
		}
	}

	return nil
}
