# crypto

A Go library to parse cryptocurrency exchange data. WIP.

Currently supports limited order book retrieval for:

- Luno
- FNB
- Kraken
- ICE3x
- AltCoinTrader

If you want to use AltCoinTrader, use the included HTTP transport to
bypass the CloudFlare protection. See `example.go` for an example.

