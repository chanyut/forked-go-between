package orderbook

import (
	"encoding/json"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestDethBids(t *testing.T) {
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

	result := book.Depth()
	s, err := json.Marshal(result)

	assert.Nil(t, err)
	cupaloy.SnapshotT(t, s)
}

func TestDepthAsks(t *testing.T) {
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

	result := book.Depth()
	s, err := json.Marshal(result)

	assert.Nil(t, err)
	cupaloy.SnapshotT(t, s)
}
