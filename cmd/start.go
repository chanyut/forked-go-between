package cmd

import (
	"net"

	"github.com/danielgatis/go-between/internal/api"
	"github.com/danielgatis/go-between/internal/logger"
	"github.com/danielgatis/go-between/internal/raft"
	"github.com/danielgatis/go-between/pkg/orderbook"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	id      string
	market  string
	port    string
	dataDir string
	join    string
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&id, "id", "", "the raft id")
	startCmd.PersistentFlags().StringVar(&market, "market", "USD/BTC", "the order book market")
	startCmd.PersistentFlags().StringVar(&port, "port", "", "the listen port")
	startCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "./data", "the raft storage dir")
	startCmd.PersistentFlags().StringVar(&join, "join", "", "when blank bootstrap a cluster otherwise joins")

	if err := startCmd.MarkPersistentFlagRequired("id"); err != nil {
		logger.GetLogrusInstance().Fatal(err)
	}

	if err := startCmd.MarkPersistentFlagRequired("market"); err != nil {
		logger.GetLogrusInstance().Fatal(err)
	}

	if err := startCmd.MarkPersistentFlagRequired("port"); err != nil {
		logger.GetLogrusInstance().Fatal(err)
	}

	if err := viper.BindPFlags(startCmd.PersistentFlags()); err != nil {
		logger.GetLogrusInstance().Fatal(err)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new node",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.GetLogrusInstance()
		book := orderbook.NewOrderBook(market)

		raftNode, err := raft.NewNode(id, port, dataDir, join, book)
		if err != nil {
			log.Fatal(err)
		}

		raftServer := raft.NewServer(raftNode)
		apiServer := api.NewServer(raftNode)

		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal(err)
		}

		mux := cmux.New(lis)
		grpcL := mux.Match(cmux.HTTP2())
		httpL := mux.Match(cmux.HTTP1Fast())

		g := new(errgroup.Group)
		g.Go(func() error { return raftServer.Serve(grpcL) })
		g.Go(func() error { return apiServer.Serve(httpL) })
		g.Go(func() error { return mux.Serve() })

		err = raftNode.JoinOrBootstrap(join, id, port)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(g.Wait())
	},
}
