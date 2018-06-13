package exchange

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type OrderType bool

const (
	BUY  OrderType = true
	SELL OrderType = false
)

// 0 - price
// 1 - volume
type OrderBook struct {
	Bids     [][2]float64
	Asks     [][2]float64
	Pair     *Pair
	Exchange Exchange
}

// 0 - price
// 1 - volume
// 2 - value
// 3 - cum_volume
// 4 - cum_value
type PreparedOrderBook struct {
	Bids [][5]float64
	Asks [][5]float64
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
	Pairs []*Pair
	API   string
	URL   string
}

// GetOrderBook fetches a single trading pair on an exchange.
func GetOrderBook(client http.Client, exc Exchange, pair *Pair) (*OrderBook, error) {

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

// GetOrderBooks fetches all trading pairs for the provided exchanges concurrently.
func GetOrderBooks(client http.Client, exchanges ...Exchange) ([]*OrderBook, error) {
	var obs []*OrderBook
	numPairs := 0
	exchangePairs := map[Exchange][]*Pair{}
	for _, e := range exchanges {
		p := e.Meta().Pairs
		numPairs += len(p)
		exchangePairs[e] = p
	}
	results := make(chan *OrderBook, numPairs)
	errors := make(chan error, numPairs)

	for exchange, pairs := range exchangePairs {

		for _, pair := range pairs {
			go func(e Exchange, p *Pair) {
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

// Build is a helper function to build a complete URL for an exchange.
// It accepts a path and query parameters.
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

func (ob *OrderBook) Prepare() *PreparedOrderBook {

	return &PreparedOrderBook{
		Bids: prepareEntries(ob.Bids),
		Asks: prepareEntries(ob.Asks),
	}
}

func prepareEntries(entries [][2]float64) [][5]float64 {

	var (
		cumValue  float64
		cumVolume float64
	)

	prepared := make([][5]float64, len(entries))

	for i, e := range entries {
		price := e[0]
		volume := e[1]
		value := price * volume
		cumValue += value
		cumVolume += volume
		prepared[i] = [5]float64{price, volume, value, cumVolume, cumValue}
	}

	return prepared

}

type Trade struct {
	OrderBook *PreparedOrderBook
	Amount    float64
	Type      OrderType
	Quote     bool
	Pair      *Pair
}

type TradeResult struct {
	Gross     float64
	Fee       float64
	Nett      float64
	GrossUnit float64
	NettUnit  float64
	Asset     *Asset
}

type Route struct {
	Legs []*RouteLeg
}

type RouteLeg struct {
	Pair      *Pair
	OrderBook *PreparedOrderBook
	Exchange  Exchange
}

type RouteRequest struct {
	Pair   *Pair
	Volume float64
	Type   OrderType
}

type RouteResult struct {
	Amount      float64
	Asset       *Asset
	Description string
	Requests    []*RouteRequest
}

func (t *Trade) Simulate() (*TradeResult, error) {

	var (
		gross     float64
		cumVolume float64
		cumValue  float64
		ass       *Asset
		fee       float64
		entries   [][5]float64
	)
	if t.Type == BUY {
		entries = t.OrderBook.Asks
	} else {
		entries = t.OrderBook.Bids
	}

	for i, entry := range entries {
		price := entry[0]
		cumVolume = entry[3]
		cumValue = entry[4]

		if t.Quote {
			if cumValue > t.Amount {
				prevIx := i - 1
				if prevIx == -1 {
					gross = t.Amount / price
				} else {
					prev := entries[prevIx]
					gross = prev[3] + (t.Amount-prev[4])/price
				}
				break
			}
		} else {
			if cumVolume > t.Amount {
				prevIx := i - 1
				if prevIx == -1 {
					gross = t.Amount * price
				} else {
					prev := entries[prevIx]
					gross = prev[4] + (t.Amount-prev[3])*price
				}
				break
			}
		}

	}

	if gross == 0 {
		total := cumVolume
		if t.Quote {
			total = cumValue
		}
		return nil, fmt.Errorf("Orderbook too small: wanted %f, total %f", t.Amount, total)
	}

	if t.Pair != nil {
		fee = t.Pair.TakerFee * gross
		if t.Type == BUY {
			ass = t.Pair.Base
		} else {
			ass = t.Pair.Quote
		}
	}

	nett := gross - fee

	return &TradeResult{
		Gross:     gross,
		Fee:       fee,
		Nett:      nett,
		GrossUnit: gross / t.Amount,
		NettUnit:  nett / t.Amount,
		Asset:     ass,
	}, nil

}

func (r *Route) Simulate(asset *Asset, amount float64) (*RouteResult, error) {

	var (
		description string
		requests    []*RouteRequest
	)

	for i, leg := range r.Legs {

		var ot OrderType

		switch asset {
		case leg.Pair.Quote:
			ot = BUY
		case leg.Pair.Base:
			ot = SELL
		default:
			return nil, fmt.Errorf("Invalid asset \"%s\" for simulation: must be \"%s\" or \"%s\"", asset.Slug, leg.Pair.Base.Slug, leg.Pair.Quote.Slug)
		}

		t := &Trade{Pair: leg.Pair, OrderBook: leg.OrderBook, Type: ot}
		res, err := t.Simulate()
		if err != nil {
			return nil, err
		}
		if ot == BUY {
			description += fmt.Sprintf("\nTRADE: buy %.4f %s on %s for %.4f %s (%.4f fee)",
				res.Gross, res.Asset.Slug, leg.Exchange.Meta().Slug, amount, asset, res.Fee)
			requests = append(requests, &RouteRequest{Pair: leg.Pair, Volume: res.Gross, Type: BUY})
		} else {
			description += fmt.Sprintf("\nTRADE: sell %.f %s on %s for %.4f %s (%.4f fee)",
				amount, asset, leg.Exchange.Meta().Slug, res.Gross, res.Asset.Slug, res.Fee)
			requests = append(requests, &RouteRequest{Pair: leg.Pair, Volume: amount, Type: SELL})
		}

		amount = res.Nett
		asset = res.Asset

		description += fmt.Sprintf("\nHOLDING: %.4f %s", amount, asset.Slug)

		if i+1 < len(r.Legs) {

			nextExchange := r.Legs[i+1].Exchange
			if nextExchange != leg.Exchange {

				withdrawalFee := 0.0 // TODO: implement exchange.calculate_withdrawal_fee(asset, amount)
				depositFee := 0.0    // TODO: implement nextExchange.calculate_deposit_fee(asset, amount - withdrawal_fee)

				description += fmt.Sprintf("FEE: %.4f %s withdrawal fee at %s", withdrawalFee, asset, leg.Exchange.Meta().Slug)
				description += fmt.Sprintf("FEE: %.4f %s deposit fee at %s", depositFee, asset, nextExchange.Meta().Slug)

				amount = amount - withdrawalFee - depositFee
				description += fmt.Sprintf("HOLDING: %.4f %s", amount, asset)
			}

		} else {
			description += "DONE"
		}
	}

	return &RouteResult{Amount: amount, Asset: asset, Description: description, Requests: requests}, nil
}
