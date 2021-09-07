package orderbook

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProcessLimitOrderBids(t *testing.T) {
	type input struct {
		OrderID  string
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
		price    decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Trades []*Trade
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "happy path - 1",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(10),
				price:    decimal.NewFromInt(100),
			},
		},
		{
			name: "happy path - 2",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.RequireFromString("0.3"),
				price:    decimal.NewFromInt(500),
			},
		},
		{
			name: "happy path - 3",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.RequireFromString("1.5"),
				price:    decimal.NewFromInt(300),
			},
		},
		{
			name: "happy path - 4",
			input: input{
				OrderID:  "4",
				traderID: "3",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.RequireFromString("1.5"),
				price:    decimal.NewFromInt(300),
			},
		},
		{
			name: "happy path - 5",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.RequireFromString("9"),
				price:    decimal.NewFromInt(1000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			given := []byte(`
				{
					"market": "USD/BTC",
					"bids": [],
					"asks": [
						{
							"id": "1",
							"traderId": "1",
							"side": "sell",
							"quantity": "5",
							"price": "500"
						},
						{
							"id": "2",
							"traderId": "2",
							"side": "sell",
							"quantity": "2",
							"price": "400"
						},
						{
							"id": "3",
							"traderId": "3",
							"side": "sell",
							"quantity": "1",
							"price": "300"
						}
					]
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			trades, err := book.ProcessLimitOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Trades: trades,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}

func TestProcessLimitOrderAsks(t *testing.T) {
	type input struct {
		OrderID  string
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
		price    decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Trades []*Trade
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "happy path - 1",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "happy path - 2",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(2),
				price:    decimal.NewFromInt(300),
			},
		},
		{
			name: "happy path - 3",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.RequireFromString("5.2"),
				price:    decimal.NewFromInt(500),
			},
		},
		{
			name: "happy path - 4",
			input: input{
				OrderID:  "4",
				traderID: "1",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.RequireFromString("0.5"),
				price:    decimal.NewFromInt(450),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			given := []byte(`
				{
					"market": "USD/BTC",
					"bids": [
						{
							"id": "1",
							"traderId": "1",
							"side": "buy",
							"quantity": "5",
							"price": "500"
						},
						{
							"id": "2",
							"traderId": "2",
							"side": "buy",
							"quantity": "1",
							"price": "400"
						},
						{
							"id": "3",
							"traderId": "3",
							"side": "buy",
							"quantity": "0.5",
							"price": "300"
						}
					],
					"asks": []
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			trades, err := book.ProcessLimitOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Trades: trades,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}

func TestProcessLimitOrderValidations(t *testing.T) {
	type input struct {
		OrderID  string
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
		price    decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Trades []*Trade
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "invalid order id",
			input: input{
				OrderID:  "",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "order already exists",
			input: input{
				OrderID:  "1",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "invalid trader id",
			input: input{
				OrderID:  "4",
				traderID: "",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "invalid market",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "xxx",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "invalid quantity",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(0),
				price:    decimal.NewFromInt(550),
			},
		},
		{
			name: "invalid price",
			input: input{
				OrderID:  "4",
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(5),
				price:    decimal.NewFromInt(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			given := []byte(`
				{
					"market": "USD/BTC",
					"bids": [
						{
							"id": "1",
							"traderId": "1",
							"side": "buy",
							"quantity": "5",
							"price": "500"
						},
						{
							"id": "2",
							"traderId": "2",
							"side": "buy",
							"quantity": "1",
							"price": "400"
						},
						{
							"id": "3",
							"traderId": "3",
							"side": "buy",
							"quantity": "0.5",
							"price": "300"
						}
					],
					"asks": []
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			trades, err := book.ProcessLimitOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Trades: trades,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}

func benchmarkProcessLimitOrder(l int, b *testing.B) {
	pickSide := func(j int) Side {
		if rand.Intn(100)%2 == 0 {
			return Buy
		} else {
			return Sell
		}
	}

	pickQty := func() decimal.Decimal {
		return decimal.NewFromInt(int64(rand.Intn(100)))
	}

	pickPrice := func() decimal.Decimal {
		return decimal.NewFromInt(int64(rand.Intn(100)))
	}

	book := NewOrderBook("USD/BTC")

	for n := 0; n < b.N; n++ {
		for j := 0; j < l; j++ {
			book.ProcessLimitOrder(strconv.Itoa(j), strconv.Itoa(j), "USD/BTC", pickSide(j), pickQty(), pickPrice())
		}
	}
}

func BenchmarkProcessLimitOrder100(b *testing.B)     { benchmarkProcessLimitOrder(100, b) }
func BenchmarkProcessLimitOrder1000(b *testing.B)    { benchmarkProcessLimitOrder(1000, b) }
func BenchmarkProcessLimitOrder10000(b *testing.B)   { benchmarkProcessLimitOrder(10000, b) }
func BenchmarkProcessLimitOrder100000(b *testing.B)  { benchmarkProcessLimitOrder(100000, b) }
func BenchmarkProcessLimitOrder1000000(b *testing.B) { benchmarkProcessLimitOrder(1000000, b) }
