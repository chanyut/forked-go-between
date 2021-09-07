package raft

import (
	"bytes"
	"encoding/json"

	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/hashicorp/raft"
	"github.com/rotisserie/eris"
)

var _ raft.FSMSnapshot = (*snapshot)(nil)

type snapshot struct {
	book *orderbook.OrderBook
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		var buffer bytes.Buffer

		if err := json.NewEncoder(&buffer).Encode(s.book); err != nil {
			return eris.Wrap(err, "enc.Encode")
		}

		if _, err := sink.Write(buffer.Bytes()); err != nil {
			return eris.Wrap(err, "sink.Write")
		}

		if err := sink.Close(); err != nil {
			return eris.Wrap(err, "sink.Close")
		}

		return nil
	}()

	if err != nil {
		if err := sink.Cancel(); err != nil {
			return eris.Wrap(err, "sink.Cancel")
		}

		return err
	}

	return nil
}

func (s *snapshot) Release() { /* nope */ }
