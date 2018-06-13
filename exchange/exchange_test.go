package exchange

import (
	"reflect"
	"testing"
)

func TestTrade(t *testing.T) {

	ob := &OrderBook{
		Bids: [][2]float64{
			{10, 10},
			{12, 1},
			{20, 10},
		},
		Asks: [][2]float64{
			{10, 10},
			{12, 1},
			{20, 10},
		},
	}

	p := ob.Prepare()

	expectedBids := [][5]float64{
		{10, 10, 100, 10, 100},
		{12, 1, 12, 11, 112},
		{20, 10, 200, 21, 312},
	}
	if !reflect.DeepEqual(p.Bids, expectedBids) {
		t.Error("Incorrect preparation of order book")
	}

	s1 := Trade{OrderBook: p, Amount: 15, Type: SELL}
	if r, _ := s1.Simulate(); r.Gross != 192 {
		t.Error("Expected gross trade result if 192")
	}

	s2 := Trade{OrderBook: p, Amount: 192, Type: SELL, Quote: true}
	if r, _ := s2.Simulate(); r.Gross != 15 {
		t.Error("Expected gross trade result if 15")
	}

	b1 := Trade{OrderBook: p, Amount: 17, Type: SELL}
	if r, _ := b1.Simulate(); r.Gross != 232 {
		t.Error("Expected gross trade result if 232")
	}

	b2 := Trade{OrderBook: p, Amount: 232, Type: SELL, Quote: true}
	if r, _ := b2.Simulate(); r.Gross != 17 {
		t.Error("Expected gross trade result if 17")
	}

}
