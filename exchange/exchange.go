package exchange

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/jacoduplessis/crypto/asset"
)

type Bids [][2]float64
type Asks [][2]float64

type OrderBook struct {
	Bids     Bids
	Asks     Asks
	Pair     *asset.Pair
	Exchange Exchange
}

type Exchange interface {
	Meta() *Meta
	GetOrderBookRequest(string) (*http.Request, error)
	ParseOrderBookResponse(io.Reader) (*OrderBook, error)
}

type RequestConfig struct {
	Method string
	Path   string
	Params map[string]string
}

type Meta struct {
	Name  string
	Slug  string
	Pairs []*asset.Pair
	API   string
	URL   string
}

func GetOrderBook(client http.Client, exc Exchange, pair *asset.Pair) (*OrderBook, error) {

	req, err := exc.GetOrderBookRequest(pair.Code)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ob, err := exc.ParseOrderBookResponse(resp.Body)
	if err != nil {
		return nil, err
	}
	ob.Pair = pair
	return ob, nil
}

func GetOrderBooks(client http.Client, exchanges ...Exchange) ([]*OrderBook, error) {
	var obs []*OrderBook
	numPairs := 0
	exchangePairs := map[Exchange][]*asset.Pair{}
	for _, e := range exchanges {
		p := e.Meta().Pairs
		numPairs += len(p)
		exchangePairs[e] = p
	}
	results := make(chan *OrderBook, numPairs)
	errors := make(chan error, numPairs)

	for exchange, pairs := range exchangePairs {

		for _, pair := range pairs {
			go func(e Exchange, p *asset.Pair) {
				ob, err := GetOrderBook(client, e, p)
				if err != nil {
					errors <- err
					return
				}
				ob.Exchange = e
				results <- ob
			}(exchange, pair)
		}

	}

	for i := 0; i < numPairs; i++ {
		select {
		case err := <-errors:
			return nil, err
		case ob := <-results:
			obs = append(obs, ob)
		}
	}
	return obs, nil
}

func Build(e Exchange, p string, params map[string]string) string {
	u, err := url.Parse(e.Meta().API)
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path.Join(u.Path, p)
	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String()
}
