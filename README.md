# Go - Between

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-between?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-between)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-between/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-between)
[![Release](https://img.shields.io/github/release/danielgatis/go-between.svg?style=flat-square)](https://github.com/danielgatis/go-between/releases/latest)

<p align="center">
    <img width="250px" src="./logo.png">
</p>

A distributed fault-tolerant order book matching engine.

### Features

- Limit orders
- Market orders
- Order book depth
- Calculate market price for a given quantity
- Standard Price/time priority
- Distributed and fault-tolerant based on [Raft Consensus Algorithm](https://raft.github.io)

### Download binaries

You can download a pre-built binary [here](https://github.com/danielgatis/go-between/releases).

### Build from source

First, [install Go](https://golang.org/doc/install).

Next, fetch and build the binary.

```bash
go get -u github.com/danielgatis/go-between
```

### How to init a cluster

First of all, start a new cluster:
```sh
go-between start --id=1 --market=USD/BTC --port=3001
```

And add a two more nodes:

```sh
go-between start --id=1 --market=USD/BTC --port=3002 --join=localhost:3001
```

```sh
go-between start --id=1 --market=USD/BTC --port=3003 --join=localhost:3001
```

### API Endpoints

Each node publishes a REST API with the following endpoints.

#### Adding a new limit order

```sh
curl --request POST \
  --url http://127.0.0.1:3001/orders/limit \
  --header 'Content-Type: application/json' \
  --data '{
	"orderId": "1",
	"traderId": "1",
	"market": "USD/BTC",
	"side": "sell",
	"quantity": "1",
	"price": "400.00"
}'
```

#### Adding a new market order

```sh
curl --request POST \
  --url http://127.0.0.1:3001/orders/market \
  --header 'Content-Type: application/json' \
  --data '{
	"orderId": "2",
	"traderId": "1",
	"market": "USD/BTC",
	"side": "sell",
	"quantity": "1",
	"price": "400.00"
}'
```

#### Cancel an order

```sh
curl --request DELETE --url http://127.0.0.1:3001/orders/1
```

#### Calculate a market price

```sh
curl --request POST \
  --url http://127.0.0.1:3001/prices/market \
  --header 'Content-Type: application/json' \
  --data '{
	"traderId": "1",
	"market": "USD/BTC",
	"side": "sell",
	"quantity": "10"
}'
```

#### Get the order book price depth

```sh
curl --request GET --url http://127.0.0.1:3001/prices/depth
```

#### Get the full order book

```sh
curl --request GET --url http://127.0.0.1:3001/orderbook
```

#### Get the raft node stats

```sh
curl --request GET --url http://127.0.0.1:3001/raft/stats
```

### License

Copyright (c) 2021-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE)

### Buy me a coffee

Liked some of my work? Buy me a coffee (or more likely a beer)

<a href="https://www.buymeacoffee.com/danielgatis" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;"></a>
