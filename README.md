# go-dax-sdk
dax.btcc.com trading api sdk

### useage:
here is an example,you may edit it in github.com/JinAirsOs/go-dax-sdk/example/dax.go

```
package main

import (
	"log"

	"github.com/JinAirsOs/go-dax-sdk/dax"
)

const (
	DAX_API_KEY    = ""
	DAX_API_SECRET = ""
)

func main() {
	daxClient := dax.New(DAX_API_KEY, DAX_API_SECRET)

	//get account info
	resp := daxClient.GetAccountInfo()
	log.Println(string(resp.Body))

	//get my open orders
	resp = daxClient.GetMyOrders("ETH_BTC")
	log.Println(string(resp.Body))

	//place an order
	resp = daxClient.PlaceOrder(
		dax.Order{CurrencyPair: "ETH_BTC", // currency pair
			Price:     "0.11111", //order price
			Quantity:  "0.00001", // order quantity
			Type:      "LIMIT",   // "LIMIT"
			BuyOrSell: "SELL",
		}) //buy or sell
	log.Println(string(resp.Body))

	/*{"symbol":"ETH_BTC","price":"0.11111","quantity":"0.00001","type":"LIMIT","side":"SELL"}*/
	//cancel an order
	resp = daxClient.CancelOrder("ETH_BTC", "0123456789abcdef0123456789abcdef")
	log.Println(string(resp.Body))
}
```