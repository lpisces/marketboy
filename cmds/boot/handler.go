package boot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
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

func createOrder(symbol, side string, qty float64, price float64) (err error) {
	// endpoint
	path := "/order"
	endpoint := fmt.Sprintf("%s://%s%s%s", Conf.RestConfig.Scheme, Conf.RestConfig.Host, Conf.RestConfig.Prefix, path)

	// params
	params := make(map[string]interface{})
	params["symbol"] = symbol
	params["side"] = side
	params["orderQty"] = qty
	params["price"] = price

	// verb
	verb := "POST"

	// sign
	expires := time.Now().Unix() + 5
	sign := getSign(Conf.AuthConfig.Secret, verb, Conf.RestConfig.Prefix+path, expires, string(mustMarshal(params)))

	// header
	header := make(map[string]string)
	header["api-expires"] = time.Now().Unix() + 5
	header["api-key"] = Conf.AuthConfig.Key
	header["api-signature"] = sign

	// request
	request := resty.R()
	for k, v := range header {
		request.SetHeader(k, v)
	}

	return
}
