package raft

import (
	"net"

	"github.com/Jille/raftadmin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	node   *Node
	server *grpc.Server
}

func NewServer(raftNode *Node) *Server {
	server := grpc.NewServer()

	raftNode.tm.Register(server)
	raftadmin.Register(server, raftNode.Raft)
	reflection.Register(server)

	return &Server{raftNode, server}
}

func (rn *Server) Serve(l net.Listener) error {
	return rn.server.Serve(l)
}
