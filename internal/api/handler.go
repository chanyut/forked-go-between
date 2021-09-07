package api

import (
	"net/http"

	"github.com/danielgatis/go-between/internal/raft"
	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

func deleteOrdersCancel(raftNode *raft.Node) echo.HandlerFunc {
	type form struct {
		OrderID string `json:"orderId"`
	}

	return func(c echo.Context) error {
		var f form
		f.OrderID = c.Param("id")

		cmd := &raft.RaftCmd{
			Op:      raft.CancelOrderOp,
			OrderID: f.OrderID,
		}

		_, err := sendCmd(raftNode, cmd)
		if err != nil {
			return eris.Wrap(err, "sendCmd(...)")
		}

		return c.NoContent(204)
	}
}

func postOrdersMarket(raftNode *raft.Node) echo.HandlerFunc {
	type form struct {
		OrderID  string          `json:"orderId"`
		TraderID string          `json:"traderId"`
		Market   string          `json:"market"`
		Side     orderbook.Side  `json:"side"`
		Quantity decimal.Decimal `json:"quantity"`
		Price    decimal.Decimal `json:"price"`
	}

	return func(c echo.Context) error {
		var f form
		if err := c.Bind(&f); err != nil {
			return err
		}

		cmd := &raft.RaftCmd{
			Op:       raft.ProcessMarketOrderOp,
			OrderID:  f.OrderID,
			TraderID: f.TraderID,
			Market:   f.Market,
			Side:     f.Side,
			Quantity: f.Quantity,
			Price:    f.Price,
		}

		r, err := sendCmd(raftNode, cmd)
		if err != nil {
			return eris.Wrap(err, "sendCmd(...)")
		}

		return c.JSON(201, r)
	}
}

func postOrdersLimit(raftNode *raft.Node) echo.HandlerFunc {
	type form struct {
		OrderID  string          `json:"orderId"`
		TraderID string          `json:"traderId"`
		Market   string          `json:"market"`
		Side     orderbook.Side  `json:"side"`
		Quantity decimal.Decimal `json:"quantity"`
		Price    decimal.Decimal `json:"price"`
	}

	return func(c echo.Context) error {
		var f form
		if err := c.Bind(&f); err != nil {
			return err
		}

		cmd := &raft.RaftCmd{
			Op:       raft.ProcessLimitOrderOp,
			OrderID:  f.OrderID,
			TraderID: f.TraderID,
			Market:   f.Market,
			Side:     f.Side,
			Quantity: f.Quantity,
			Price:    f.Price,
		}

		r, err := sendCmd(raftNode, cmd)
		if err != nil {
			return eris.Wrap(err, "sendCmd(...)")
		}

		return c.JSON(201, r)
	}
}

func postPricesMarket(raftNode *raft.Node) echo.HandlerFunc {
	type form struct {
		TraderID string          `json:"traderId"`
		Market   string          `json:"market"`
		Side     orderbook.Side  `json:"side"`
		Quantity decimal.Decimal `json:"quantity"`
	}

	return func(c echo.Context) error {
		var f form
		if err := c.Bind(&f); err != nil {
			return err
		}

		r, err := raftNode.Book().MarketPrice(
			f.TraderID,
			f.Market,
			f.Side,
			f.Quantity,
		)

		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnprocessableEntity,
				Message:  err.Error(),
				Internal: err,
			}
		}

		return c.JSON(201, r)
	}
}

func getPricesDepth(raftNode *raft.Node) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, raftNode.Book().Depth())
	}
}

func getBook(raftNode *raft.Node) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, raftNode.Book())
	}
}

func getRaftStats(raftNode *raft.Node) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, raftNode.Stats())
	}
}
