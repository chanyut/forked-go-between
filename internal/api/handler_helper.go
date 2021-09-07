package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielgatis/go-between/internal/raft"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

func sendCmd(node *raft.Node, cmd *raft.RaftCmd) (interface{}, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(cmd)
	if err != nil {
		return nil, eris.Wrap(err, "json.NewEncoder")
	}

	r := node.Apply(buffer.Bytes(), 10*time.Second)
	if err := r.Error(); err != nil {
		return nil, eris.Wrap(err, "raftNode.Apply")
	}

	if err, ok := r.Response().(error); ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  err.Error(),
			Internal: err,
		}
	}

	return r.Response(), nil
}
