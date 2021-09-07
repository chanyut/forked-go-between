package orderbook

import (
	"encoding/json"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProcessMarketOrderBids(t *testing.T) {
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
	}{}

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

			trades, err := book.ProcessMarketOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

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

func TestProcessMarketOrderAsks(t *testing.T) {
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
	}{}

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

			trades, err := book.ProcessMarketOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

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

func TestProcessMarketOrderValidations(t *testing.T) {
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

			trades, err := book.ProcessMarketOrder(tt.input.OrderID, tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity, tt.input.price)

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
