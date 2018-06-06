package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jacoduplessis/crypto/exchange"
	"github.com/jacoduplessis/crypto/scraper"
)

func main() {

	client := http.Client{Transport: scraper.NewTransport(http.DefaultTransport)}

	luno := &exchange.Luno{}
	kraken := &exchange.Kraken{}
	ice := &exchange.ICE{}
	fnb := &exchange.FNB{}
	alt := &exchange.AltCoinTrader{}

	books, err := exchange.GetOrderBooks(client, luno, kraken, ice, fnb, alt)
	if err != nil {
		log.Fatal(err)
	}

	for _, ob := range books {
		fmt.Println(ob.Exchange.Meta().Name, ob.Pair.Code, len(ob.Bids), ob.Bids[:1], len(ob.Asks), ob.Asks[:1])
	}

}
