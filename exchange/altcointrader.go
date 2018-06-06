package exchange

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jacoduplessis/crypto/asset"
)

type AltCoinTrader struct {
}

func (alt *AltCoinTrader) Meta() *Meta {

	return &Meta{
		Name: "AltCoinTrader",
		Slug: "alt",
		API:  "https://www.altcointrader.co.za",
		Pairs: []*asset.Pair{
			{Base: asset.Bitcoin, Quote: asset.Rand, Code: "/"},
			{Base: asset.Ripple, Quote: asset.Rand, Code: "/xrp"},
		},
	}
}

func (alt *AltCoinTrader) GetOrderBookRequest(pairCode string) (*http.Request, error) {

	u := Build(alt, pairCode, nil)
	return http.NewRequest("GET", u, nil)
}

func (alt *AltCoinTrader) ParseOrderBookResponse(body io.Reader) (*OrderBook, error) {

	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, err
	}

	var bids [][2]float64
	var asks [][2]float64

	doc.Find("tr.orderUdSell").Each(func(i int, s *goquery.Selection) {
		priceString := strings.TrimSpace(s.Find(".orderUdSPr").Text())
		volumeString := strings.TrimSpace(s.Find(".orderUdSAm").Text())

		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			return
		}
		volume, err := strconv.ParseFloat(volumeString, 64)
		if err != nil {
			return
		}

		asks = append(asks, [2]float64{price, volume})

	})

	doc.Find("tr.orderUdBuy").Each(func(i int, s *goquery.Selection) {
		priceString := strings.TrimSpace(s.Find(".orderUdBPr").Text())
		volumeString := strings.TrimSpace(s.Find(".orderUdBAm").Text())

		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			return
		}
		volume, err := strconv.ParseFloat(volumeString, 64)
		if err != nil {
			return
		}

		bids = append(bids, [2]float64{price, volume})

	})

	return &OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil
}
