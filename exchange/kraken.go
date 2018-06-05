package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jacoduplessis/crypto/asset"
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
		Pairs: []*asset.Pair{
			{Base: asset.Bitcoin, Quote: asset.Euro, Code: "XXBTZEUR"},
			{Base: asset.Ripple, Quote: asset.Euro, Code: "XXRPZEUR"},
			{Base: asset.Ripple, Quote: asset.Bitcoin, Code: "XXRPXXBT"},
			{Base: asset.Litecoin, Quote: asset.Euro, Code: "XLTCZEUR"},
			{Base: asset.Litecoin, Quote: asset.Bitcoin, Code: "XLTCXXBT"},
			{Base: asset.Bitcoincash, Quote: asset.Euro, Code: "BCHEUR"},
			{Base: asset.Bitcoincash, Quote: asset.Bitcoin, Code: "BCHXBT"},
			{Base: asset.Ether, Quote: asset.Euro, Code: "XETHZEUR"},
			{Base: asset.Ether, Quote: asset.Bitcoin, Code: "XETHXXBT"},
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
