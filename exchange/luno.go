package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Luno struct {
	APIKey    string
	APISecret string
}

func (ln *Luno) Meta() *Meta {
	return &Meta{
		Name: "Luno",
		Slug: "luno",
		API:  "https://api.mybitx.com/api/1/",
		Pairs: []*Pair{
			{Base: Bitcoin, Quote: Rand, Code: "XBTZAR"},
			{Base: Ether, Quote: Bitcoin, Code: "ETHXBT"},
		},
	}
}

func (ln *Luno) GetOrderBookRequest(pairCode string) (*http.Request, error) {

	u := Build(ln, "orderbook", map[string]string{"pair": pairCode})
	return http.NewRequest("GET", u, nil)
}

func (ln *Luno) ParseOrderBookResponse(body io.Reader) (*OrderBook, error) {

	type Entry struct {
		Volume string
		Price  string
	}

	var d struct {
		Bids []Entry
		Asks []Entry
	}

	parse := func(responseEntries []Entry) ([][2]float64, error) {
		entries := make([][2]float64, len(responseEntries))

		for i, e := range responseEntries {
			price, err := strconv.ParseFloat(e.Price, 64)
			if err != nil {
				return nil, err
			}
			volume, err := strconv.ParseFloat(e.Volume, 64)
			if err != nil {
				return nil, err
			}
			entries[i] = [2]float64{price, volume}
		}

		return entries, nil
	}

	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		return nil, err
	}

	bids, err := parse(d.Bids)
	if err != nil {
		return nil, err
	}
	asks, err := parse(d.Asks)
	if err != nil {
		return nil, err
	}

	return &OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil
}
