package api

import (
	"github.com/danielgatis/go-between/internal/logger"
	"github.com/danielgatis/go-between/internal/raft"
	"github.com/labstack/echo/v4/middleware"
)

func configureRoutes(server *Server, raftNode *raft.Node) {
	server.Use(requestLogger(logger.GetLogrusInstance()))
	server.Use(middleware.Recover())

	leader := server.Group("", checkRaftLeader(raftNode))
	leader.DELETE("/orders/:id", deleteOrdersCancel(raftNode))
	leader.POST("/orders/limit", postOrdersLimit(raftNode))
	leader.POST("/orders/market", postOrdersMarket(raftNode))

	any := server.Group("")
	any.POST("/prices/market", postPricesMarket(raftNode))
	any.GET("/prices/depth", getPricesDepth(raftNode))
	any.GET("/orderbook", getBook(raftNode))
	any.GET("/raft/stats", getRaftStats(raftNode))
}
