package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
		Pairs: []*Pair{
			{Base: Bitcoin, Quote: Euro, Code: "XXBTZEUR"},
			{Base: Ripple, Quote: Euro, Code: "XXRPZEUR"},
			{Base: Ripple, Quote: Bitcoin, Code: "XXRPXXBT"},
			{Base: Litecoin, Quote: Euro, Code: "XLTCZEUR"},
			{Base: Litecoin, Quote: Bitcoin, Code: "XLTCXXBT"},
			{Base: Bitcoincash, Quote: Euro, Code: "BCHEUR"},
			{Base: Bitcoincash, Quote: Bitcoin, Code: "BCHXBT"},
			{Base: Ether, Quote: Euro, Code: "XETHZEUR"},
			{Base: Ether, Quote: Bitcoin, Code: "XETHXXBT"},
		},
	}
}

func (kr *Kraken) ParseOrderBookResponse(body io.Reader) (*OrderBook, error) {

	parse := func(responseEntries [][2]string) ([][2]float64, error) {
		entries := make([][2]float64, len(responseEntries))

		for i, e := range responseEntries {
			price, err := strconv.ParseFloat(e[0], 64)
			if err != nil {
				return nil, err
			}
			volume, err := strconv.ParseFloat(e[1], 64)
			if err != nil {
				return nil, err
			}
			entries[i] = [2]float64{price, volume}
		}

		return entries, nil
	}

	type Data struct {
		Asks [][2]string
		Bids [][2]string
	}
	var d struct {
		Errors []string
		Result map[string]Data
	}

	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		return nil, err
	}
	if len(d.Errors) != 0 {
		return nil, errors.New(strings.Join(d.Errors, ","))
	}

	var data Data

	for _, v := range d.Result {
		data = v
		break
	}

	bids, err := parse(data.Bids)
	if err != nil {
		return nil, err
	}
	asks, err := parse(data.Asks)
	if err != nil {
		return nil, err
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
