package dax

type Order struct {
	CurrencyPair string // currency pair
	Price        string //order price
	Quantity     string // order quantity
	Type         string // "LIMIT"
	BuyOrSell    string //"BUY" "SELL"
}
