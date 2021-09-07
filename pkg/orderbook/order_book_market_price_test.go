package orderbook

import (
	"encoding/json"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMarketPriceBids(t *testing.T) {
	type input struct {
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Result *MarketPrice
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "happy path - 1",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "happy path - 2",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(1),
			},
		},
		{
			name: "empty book",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "skip same trader",
			input: input{
				traderID: "2",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(6),
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
							"quantity": "2",
							"price": "200"
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
							"quantity": "2",
							"price": "600"
						}
					]
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			result, err := book.MarketPrice(tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Result: result,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}

func TestMarketPriceAsks(t *testing.T) {
	type input struct {
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Result *MarketPrice
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "happy path - 1",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "happy path - 2",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(1),
			},
		},
		{
			name: "empty book",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "skip same trader",
			input: input{
				traderID: "2",
				market:   "USD/BTC",
				side:     Sell,
				quantity: decimal.NewFromInt(6),
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
							"quantity": "2",
							"price": "200"
						},
						{
							"id": "2",
							"traderId": "2",
							"side": "buy",
							"quantity": "2",
							"price": "400"
						},
						{
							"id": "3",
							"traderId": "3",
							"side": "buy",
							"quantity": "2",
							"price": "600"
						}
					],
					"asks": []
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			result, err := book.MarketPrice(tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Result: result,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}

func TestMarketPriceValidations(t *testing.T) {
	type input struct {
		traderID string
		market   string
		side     Side
		quantity decimal.Decimal
	}

	type snapshot struct {
		Book   *OrderBook
		Result *MarketPrice
		Err    string
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "invalid trader id",
			input: input{
				traderID: "",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "invalid market",
			input: input{
				traderID: "4",
				market:   "BRL/BTC",
				side:     Buy,
				quantity: decimal.NewFromInt(6),
			},
		},
		{
			name: "invalid quantity",
			input: input{
				traderID: "4",
				market:   "USD/BTC",
				side:     Buy,
				quantity: decimal.Zero,
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
							"quantity": "2",
							"price": "200"
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
							"quantity": "2",
							"price": "600"
						}
					]
				}
			`)

			var book OrderBook
			err := json.Unmarshal(given, &book)
			assert.Nil(t, err)

			result, err := book.MarketPrice(tt.input.traderID, tt.input.market, tt.input.side, tt.input.quantity)

			var errorStr string
			if err != nil {
				errorStr = err.Error()
			}

			s, err := json.Marshal(&snapshot{
				Book:   &book,
				Result: result,
				Err:    errorStr,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}
}
