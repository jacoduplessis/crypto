package asset

type Pair struct {
	Base          *Asset
	Quote         *Asset
	Code          string
	CryptoWatchID string
	TakerFee      float64
	MakerFee      float64
}

type Asset struct {
	Slug   string
	Name   string
	Code   string
	Symbol string
}

var Bitcoin = &Asset{"bitcoin", "Bitcoin", "xbt", "฿"}
var Ether = &Asset{"ether", "Ether", "eth", "Ξ"}
var Litecoin = &Asset{"litecoin", "Litecoin", "ltc", "Ł"}
var Bitcoincash = &Asset{"bitcoincash", "BitcoinCash", "bch", "฿"}
var Ripple = &Asset{"ripple", "Ripple", "xrp", "Ʀ"}

var Euro = &Asset{"euro", "Euro", "eur", "€"}
var Rand = &Asset{"rand", "Rand", "zar", "R"}

func GetAllCrypto() []*Asset {

	return []*Asset{
		Bitcoin,
		Ether,
		Litecoin,
		Bitcoincash,
		Ripple,
	}
}

func GetAllFiat() []*Asset {
	return []*Asset{
		Euro,
		Rand,
	}
}
