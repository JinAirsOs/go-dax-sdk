package dax

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type SubscribeCmd struct {
	MsgType   string `json:"MsgType"`   //example "QuoteRequest"
	CRID      string `json:"CRID"`      //example unique id
	Symbol    string `json:"Symbol"`    //example ETH_BTC
	QuoteType int    `json:"QuoteType"` //example 2
}

const (
	DAX_WSS_API = "wss://ws-dax.btcc.com/ws/pub"
)

func (d *Dax) SubscribeExchange(market string, dataCh chan<- []byte, stop <-chan struct{}) {
	ws, _, err := websocket.DefaultDialer.Dial(DAX_WSS_API, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	crid := uuid.NewV4()
	subcmd, _ := json.Marshal(SubscribeCmd{
		MsgType:   "QuoteRequest",
		CRID:      crid.String(),
		Symbol:    market,
		QuoteType: 2,
	})

	err = ws.WriteMessage(websocket.TextMessage, subcmd)

	if err != nil {
		panic(err)
	}

	go func() {
		defer ws.Close()
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//log.Printf("recv: %s", message)
			dataCh <- message
		}
	}()

	select {
	case <-stop:
		log.Println("receive stop signal")
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		ws.Close()
	}
}
