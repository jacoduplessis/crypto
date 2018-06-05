package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/jacoduplessis/crypto/asset"
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
		Pairs: []*asset.Pair{
			{Base: asset.Bitcoin, Quote: asset.Rand, Code: "XBTZAR"},
			{Base: asset.Ether, Quote: asset.Bitcoin, Code: "ETHXBT"},
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

	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		return nil, err
	}

	bids := Bids{}
	asks := Asks{}

	for _, bid := range d.Bids {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(bid.Volume, 64)
		if err != nil {
			return nil, err
		}
		bids = append(bids, [2]float64{price, volume})
	}

	for _, ask := range d.Asks {
		price, err := strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(ask.Volume, 64)
		if err != nil {
			return nil, err
		}
		asks = append(asks, [2]float64{price, volume})
	}

	return &OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil
}
