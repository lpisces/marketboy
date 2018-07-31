package boot

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
	path := ep.Path
	endpoint := fmt.Sprintf("%s://%s%s%s", ep.Scheme, ep.Host, ep.Prefix, ep.Path)
	if ep.Port != "" {
		endpoint := fmt.Sprintf("%s://%s%s%s%s", ep.Scheme, ep.Host, ":"+ep.Port, ep.Prefix, ep.Path)
	}

	// verb
	verb := ep.Verb

	// sign
	expires := time.Now().Unix() + 5
	sign := getSign(ep.Secret, ep.Verb, ep.Prefix+ep.Path, expires, string(mustMarshal(ep.Params)))

	// header
	header := make(map[string]string)
	header["api-expires"] = expires
	header["api-key"] = ep.Key
	header["api-signature"] = sign

	// request instance
	request := resty.R()

	// body
	request.SetBody(string(mustMarshal(ep.Params)))

	// header
	request.SetHeader("Accept", "application/json")
	for k, v := range header {
		request.SetHeader(k, v)
	}

	// response
	request.SetResult(ep.Result)

	// do request
	resp, err := request.Post(endpoint)
	if err != nil {
		return
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code: %v, msg: %s", resp.StatusCode, string(resp.Body))
		return
	}

	ep.Response = resp
	return
}
