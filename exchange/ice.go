package exchange

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jacoduplessis/crypto/asset"
)

type ICE struct {
	APIKey    string
	APISecret string
}

func (ice *ICE) Meta() *Meta {

	return &Meta{
		Name: "ICE",
		Slug: "ice",
		API:  "https://ice3x.com/api/v1/",
		Pairs: []*asset.Pair{
			{Base: asset.Bitcoin, Quote: asset.Rand, Code: "3"},
			// {Base: asset.Litecoin, Quote: asset.Rand, Code: "6"},
			// {Base: asset.Ether, Quote: asset.Rand, Code: "11"},
			// {Base: asset.Ether, Quote: asset.Bitcoin, Code: "13"},
			// {Base: asset.Bitcoincash, Quote: asset.Bitcoin, Code: "14"},
			// {Base: asset.Bitcoincash, Quote: asset.Rand, Code: "15"},
			// {Base: asset.Litecoin, Quote: asset.Bitcoin, Code: "16"},
		},
	}
}

func (ice *ICE) GetOrderBookRequest(pairCode string) (*http.Request, error) {
	u := Build(ice, "orderbook/info", map[string]string{"pair_id": pairCode})
	return http.NewRequest("GET", u, nil)
}

func (ice *ICE) ParseOrderBookResponse(body io.Reader) (*OrderBook, error) {

	var d struct {
		Response struct {
			Entities struct {
				Bids []struct {
					Price  float64
					Amount float64
				}
				Asks []struct {
					Price  float64
					Amount float64
				}
			}
		}
	}

	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		return nil, err
	}

	var bids Bids
	var asks Asks

	for _, bid := range d.Response.Entities.Bids {
		bids = append(bids, [2]float64{bid.Price, bid.Amount})
	}

	for _, ask := range d.Response.Entities.Asks {
		asks = append(asks, [2]float64{ask.Price, ask.Amount})
	}

	return &OrderBook{
		Asks: asks,
		Bids: bids,
	}, nil

}
