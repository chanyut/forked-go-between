package raft

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	store "github.com/BBVA/raft-badger"
	transport "github.com/Jille/raft-grpc-transport"
	adminpb "github.com/Jille/raftadmin/proto"
	"github.com/danielgatis/go-between/internal/logger"
	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
)

type Node struct {
	*raft.Raft
	fsm *StateMachine
	tm  *transport.Manager
}

func NewNode(id, port, dir, join string, book *orderbook.OrderBook) (*Node, error) {
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	c := raft.DefaultConfig()
	c.Logger = logger.NewHCLogAdapter("main")
	c.LocalID = raft.ServerID(id)

	baseDir := filepath.Join(dir, id)

	ldbPath := filepath.Join(baseDir, "logs.dat")
	ldbOpts := badger.DefaultOptions(ldbPath).WithLogger(logger.NewBadgerLogAdapter())
	ldb, err := store.New(store.Options{Path: ldbPath, BadgerOptions: &ldbOpts})
	if err != nil {
		return nil, eris.Wrapf(err, `store.NewBadgerStore(%q)`, ldbPath)
	}

	sdbPath := filepath.Join(baseDir, "stable.dat")
	sdbOpts := badger.DefaultOptions(sdbPath).WithLogger(logger.NewBadgerLogAdapter())
	sdb, err := store.New(store.Options{Path: sdbPath, BadgerOptions: &sdbOpts})
	if err != nil {
		return nil, eris.Wrapf(err, `store.NewBadgerStore(%q)`, sdbPath)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, eris.Wrapf(err, `raft.NewFileSnapshotStore(%q, ...)`, baseDir)
	}

	fsm := NewStateMachine(book)
	tm := transport.New(raft.ServerAddress(addr), []grpc.DialOption{grpc.WithInsecure()})

	node, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return nil, eris.Wrapf(err, "raft.NewRaft")
	}

	return &Node{node, fsm, tm}, nil
}

func (n *Node) JoinOrBootstrap(join, id, port string) error {
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	if join == "" {
		n.BootstrapCluster(raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(id),
					Address:  raft.ServerAddress(addr),
				},
			},
		})
	} else {
		conn, err := grpc.Dial(join, grpc.WithInsecure())
		if err != nil {
			return eris.Wrapf(err, "grpc.Dial")
		}

		defer conn.Close()
		client := adminpb.NewRaftAdminClient(conn)

		req := &adminpb.AddVoterRequest{
			Id:            id,
			Address:       addr,
			PreviousIndex: 0,
		}

		_, err = client.AddVoter(context.Background(), req)
		if err != nil {
			return eris.Wrapf(err, "client.AddVoter")
		}
	}

	return nil
}

func (n *Node) IsLeader() bool {
	return n.State() == raft.Leader
}

func (n *Node) Book() *orderbook.OrderBook {
	return n.fsm.book
}
