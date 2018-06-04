package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jacoduplessis/crypto/exchange"
)

func main() {

	client := http.Client{}
	luno := &exchange.Luno{}
	kraken := &exchange.Kraken{}
	ice := &exchange.ICE{}
	fnb := &exchange.FNB{}

	books, err := exchange.GetOrderBooks(client, fnb, ice, kraken, luno)
	if err != nil {
		log.Fatal(err)
	}

	for _, ob := range books {
		fmt.Println(ob.Exchange.Meta().Name, ob.Pair.Code, len(ob.Bids), len(ob.Asks))
	}

}
