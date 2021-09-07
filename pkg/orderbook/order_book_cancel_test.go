package orderbook

import (
	"encoding/json"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestCancel(t *testing.T) {
	type input struct {
		OrderID string
	}

	type snapshot struct {
		Book  *OrderBook
		Order *Order
	}

	tests := []struct {
		name  string
		input input
	}{
		{
			name: "found",
			input: input{
				OrderID: "1",
			},
		},
		{
			name: "not found",
			input: input{
				OrderID: "foo",
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

			order := book.CancelOrder(tt.input.OrderID)

			s, err := json.Marshal(&snapshot{
				Book:  &book,
				Order: order,
			})

			assert.Nil(t, err)
			cupaloy.SnapshotT(t, s)
		})
	}

}
