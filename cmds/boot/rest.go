package boot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"time"
)

type (
	Endpoint struct {
		Verb string
		Path string
		*RestConfig
		*AuthConfig
		Params   map[string]interface{}
		Result   interface{}
		Response *resty.Response
	}

	OrderResponse struct {
		OrderID               string  `json:"orderID"`
		ClOrdID               string  `json:"clOrdID"`
		ClOrdLinkID           string  `json:"clOrdLinkID"`
		Account               float64 `json:"account"`
		Symbol                string  `json:"symbol"`
		Side                  string  `json:"side"`
		SimpleOrderQty        float64 `json:"simpleOrderQty"`
		OrderQty              float64 `json:"orderQty"`
		Price                 float64 `json:"price"`
		DisplayQty            float64 `json:"displayQty"`
		StopPx                float64 `json:"stopPx"`
		PegOffsetValue        float64 `json:"pegOffsetValue"`
		PegPriceType          string  `json:"pegPriceType"`
		Currency              string  `json:"currency"`
		SettlCurrency         string  `json:"settlCurrency"`
		OrdType               string  `json:"ordType"`
		TimeInForce           string  `json:"timeInForce"`
		ExecInst              string  `json:"execInst"`
		ContingencyType       string  `json:"contingencyType"`
		ExDestination         string  `json:"exDestination"`
		OrdStatus             string  `json:"ordStatus"`
		Triggered             string  `json:"triggered"`
		WorkingIndicator      bool    `json:"workingIndicator"`
		OrdRejReason          string  `json:"ordRejReason"`
		SimpleLeavesQty       float64 `json:"simpleLeavesQty"`
		LeavesQty             float64 `json:"leavesQty"`
		SimpleCumQty          float64 `json:"simpleCumQty"`
		CumQty                float64 `json:"cumQty"`
		AvgPx                 float64 `json:"avgPx"`
		Text                  string  `json:"text"`
		TransactTime          string  `json:"transactTime"`
		Timestamp             string  `json:"timestamp"`
		MultiLegReportingType string  `json:"multiLegReportingType"`
	}
)

func (ep *Endpoint) Do() (err error) {
	// endpoint
	endpoint := fmt.Sprintf("%s://%s%s%s", ep.Scheme, ep.Host, ep.Prefix, ep.Path)
	if ep.Port != "" {
		endpoint = fmt.Sprintf("%s://%s%s%s%s", ep.Scheme, ep.Host, ":"+ep.Port, ep.Prefix, ep.Path)
	}

	// sign
	expires := time.Now().Unix() + 5
	sign := getSign(ep.Secret, ep.Verb, ep.Prefix+ep.Path, expires, string(mustMarshal(ep.Params)))

	// header
	header := make(map[string]interface{})
	header["api-expires"] = expires
	header["api-key"] = ep.Key
	header["api-signature"] = sign

	// request instance
	request := resty.R()

	// body
	request.SetBody(string(mustMarshal(ep.Params)))

	// header
	request.SetHeader("Accept", "application/json")
	request.SetHeader("Content-Type", "application/json")
	for k, v := range header {
		request.SetHeader(k, fmt.Sprintf("%v", v))
	}

	// response
	if ep.Result != nil {
		request.SetResult(ep.Result)
	}

	// do request
	switch ep.Verb {
	case "POST":
		ep.Response, err = request.Post(endpoint)
	case "GET":
		ep.Response, err = request.Get(endpoint)
	case "DELETE":
		ep.Response, err = request.Delete(endpoint)
	case "PUT":
		ep.Response, err = request.Put(endpoint)
	default:
		err = fmt.Errorf("verb not supported: %s", ep.Verb)
		return
	}
	if err != nil {
		return
	}

	respHeader := ep.Response.Header()
	log.Info(ep.Params)
	log.Infof("x-ratelimit-remaining: %v", respHeader["X-Ratelimit-Remaining"])

	limit, _ := strconv.ParseInt(respHeader["X-Ratelimit-Remaining"][0], 10, 64)
	if limit < 30 {
		return fmt.Errorf("op limit")
	}

	// check status code
	if ep.Response.StatusCode() != http.StatusOK {
		err = fmt.Errorf("status code: %v, msg: %s", ep.Response.StatusCode(), string(ep.Response.Body()))
		return
	}

	return
}
