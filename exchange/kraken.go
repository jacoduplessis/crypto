package exchange

import (
	"net/http"
	"github.com/jacoduplessis/crypto/asset"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Kraken struct {
	APIKey    string
	APISecret string
}

func (kr *Kraken) Meta() *Meta {
	return &Meta{
		Name: "Kraken",
		Slug: "kraken",
		API:  "https://api.kraken.com/0/",
		Pairs: []*asset.Pair{
			{Base: asset.Bitcoin, Quote: asset.Euro, Code: "XXBTZEUR"},
			{Base: asset.Ripple, Quote: asset.Euro, Code: "XXRPZEUR"},
			{Base: asset.Litecoin, Quote: asset.Euro, Code: "XLTCZEUR"},
		},
	}
}

func (kr *Kraken) ParseOrderBookResponse(b []byte) (*OrderBook, error) {

	type Data struct {
		Asks [][2]string
		Bids [][2]string
	}
	var d struct {
		Errors []string
		Result map[string]Data
	}

	err := json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}
	if len(d.Errors) != 0 {
		return nil, errors.New(strings.Join(d.Errors, ","))
	}

	bids := Bids{}
	asks := Asks{}

	var data Data

	for _, v := range d.Result {
		data = v
		break
	}

	for _, bid := range data.Bids {
		price, err := strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, err
		}
		bids = append(bids, [2]float64{price, volume})
	}

	for _, ask := range data.Asks {
		price, err := strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(ask[1], 64)
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

func (kr *Kraken) GetOrderBookRequest(pairCode string) (*http.Request, error) {

	u := Build(kr, "public/Depth", map[string]string{"pair": pairCode})
	return http.NewRequest("GET", u, nil)
}
