package exchange

import (
	"encoding/json"
	"io"
	"net/http"
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
		Pairs: []*Pair{
			{Base: Bitcoin, Quote: Rand, Code: "3"},
			{Base: Litecoin, Quote: Rand, Code: "6"},
			{Base: Ether, Quote: Rand, Code: "11"},
			{Base: Ether, Quote: Bitcoin, Code: "13"},
			{Base: Bitcoincash, Quote: Bitcoin, Code: "14"},
			{Base: Bitcoincash, Quote: Rand, Code: "15"},
			{Base: Litecoin, Quote: Bitcoin, Code: "16"},
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

	var bids [][2]float64
	var asks [][2]float64

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
