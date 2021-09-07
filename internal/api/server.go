package api

import (
	"net"

	"github.com/danielgatis/go-between/internal/logger"
	"github.com/danielgatis/go-between/internal/raft"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(raftNode *raft.Node) *Server {
	e := echo.New()
	e.HideBanner = true
	e.DisableHTTP2 = true
	e.Logger = logger.NewEchoLogAdapter()

	server := &Server{e}
	configureRoutes(server, raftNode)

	return server
}

func (s *Server) Serve(l net.Listener) error {
	return s.Server.Serve(l)
}
