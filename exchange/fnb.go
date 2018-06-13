package exchange

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FNB struct {
}

func (fnb *FNB) Meta() *Meta {
	return &Meta{
		Name: "FNB",
		Slug: "fnb",
		API:  "https://www.fnb.co.za/",
		Pairs: []*Pair{
			{Base: Euro, Quote: Rand, Code: "EURZAR"},
		},
	}
}

func (fnb *FNB) GetOrderBookRequest(string) (*http.Request, error) {
	u := Build(fnb, "Controller", map[string]string{"nav": "rates.forex.list.ForexRatesList"})
	return http.NewRequest("GET", u, nil)
}

func (fnb *FNB) ParseOrderBookResponse(body io.Reader) (*OrderBook, error) {

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	row := doc.Find("table").Eq(0).Find("tr").Eq(1)
	cells := row.Find("td")
	// code := cells.Eq(1).Text()
	askString := strings.TrimSpace(cells.Eq(2).Text())
	bidString := strings.TrimSpace(cells.Eq(3).Text())

	ask, err := strconv.ParseFloat(askString, 64)
	if err != nil {
		return nil, err
	}
	bid, err := strconv.ParseFloat(bidString, 64)
	if err != nil {
		return nil, err
	}

	return &OrderBook{
		Asks: [][2]float64{{ask, 9999999}},
		Bids: [][2]float64{{bid, 9999999}},
	}, nil
}
