//Package dax provides a trading api sdk and subscribe websocket
package dax

import (
	"net/http"
	"net/url"

	"github.com/JinAirsOs/go-dax-sdk/aws"
)

const (
	API_BASE = "https://api-dax.btcc.com/tradeapi/" // Dax API endpoint
)

type Dax struct {
	ApiBase     string
	Credentials *aws.Credentials
	Client      *http.Client
}

func New(apiKey, apiSecret string) *Dax {
	var d Dax
	d.Client = &http.Client{}
	d.Credentials = &aws.Credentials{
		AccessKeyID:     apiKey,
		SecretAccessKey: apiSecret,
	}
	d.ApiBase = API_BASE
	return &d
}

func (d *Dax) GetAccountInfo() (*HttpResponse, error) {
	uri := "balance"
	resp, err := d.doRequest("GET", uri, nil)

	return resp, err
}

//for open orders go to subscribe websocket api
func (d *Dax) GetMyOrders(currencyPair string) (*HttpResponse, error) {
	uri := "orders/" + currencyPair + "/open"
	resp, err := d.doRequest("GET", uri, nil)
	return resp, err
}

func (d *Dax) PlaceOrder(order Order) (*HttpResponse, error) {
	uri := "orders"
	data := url.Values{}
	data.Set("symbol", order.CurrencyPair)
	data.Add("price", order.Price)
	data.Add("quantity", order.Quantity)
	data.Add("type", order.Type)
	data.Add("side", order.BuyOrSell)

	resp, err := d.doRequest("POST", uri, []byte(data.Encode()))
	return resp, err
}

func (d *Dax) CancelOrder(currencyPair, oid string) (*HttpResponse, error) {
	uri := "orders/" + currencyPair + "/" + oid
	resp, err := d.doRequest("DELETE", uri, nil)
	return resp, err
}
