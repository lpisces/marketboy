package boot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"net/http"
)

var (
	orderBook map[string]Order10
)

type (
	OrderBook10 struct {
		Table  string
		Action string
		Data   []Order10
	}

	Order10 struct {
		Symbol    string `json:"symbol"`
		Bids      []Bid  `json:"bids"`
		Asks      []Ask  `json:"asks"`
		Timestamp string `json:"timestamp"`
	}
	Bid []float64
	Ask []float64
)

func dispatch(msg []byte) (err error) {
	message := string(msg)

	if message == "pong" {
		return handlePing(message)
	}

	topic := gjson.GetBytes(msg, "table")

	switch topic.String() {
	case "orderBook10":
		return handleOrderBook10(msg)
	default:
		return
	}

	return
}

func handlePing(msg string) (err error) {
	log.Debug(msg)
	return
}

func handleOrderBook10(msg []byte) (err error) {
	ob := &OrderBook10{}
	if err = json.Unmarshal(msg, ob); err != nil {
		return
	}

	if ob.Action == "partial" || ob.Action == "update" {
		orderBook = make(map[string]Order10)
		for _, order := range ob.Data {
			orderBook[order.Symbol] = order
		}
		return
	}
	return
}
