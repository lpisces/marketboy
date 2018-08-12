package boot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math"
	"reflect"
	//"strings"
)

var (
	orderBook10 map[string]OrderBook10
	position    map[string]Position
	order       []Order
	operate     chan Operate
)

type (
	Operate struct {
		Action string
		Params map[string]interface{}
	}

	OrderBook10Msg struct {
		Table  string
		Action string
		Data   []OrderBook10
	}

	OrderBook10 struct {
		Symbol    string `json:"symbol"`
		Bids      []Bid  `json:"bids"`
		Asks      []Ask  `json:"asks"`
		Timestamp string `json:"timestamp"`
	}
	Bid []float64
	Ask []float64

	PositionMsg struct {
		Table  string     `json:"table"`
		Action string     `json:"Action"`
		Data   []Position `json:"data"`
	}

	Position struct {
		Account              float64 `json:"account"`
		Symbol               string  `json:"symbol"`
		Currency             string  `json:"currency"`
		InitMarginReq        float64 `json:"initMarginReq"`        // 初始保证金
		MaintMarginReq       float64 `json:"maintMarginReq"`       // 维持保证金
		Leverage             float64 `json:"leverage"`             // 杠杆率
		RiskLimit            float64 `json:"riskLimit"`            //	风险限额
		CrossMargin          bool    `json:"crossMargin"`          // 全仓保证金(false)/逐仓保证金(true)
		DeleveragePercentile float64 `json:"deleveragePercentile"` // 自动减仓百分比 越大越先减仓
		RealisedPnl          float64 `json:"realisedPnl"`          // 已实现盈亏
		UnrealisedPnl        float64 `json:"unrealisedPnl"`        // 未实现盈亏
		HomeNotional         float64 `json:"homeNotional"`         // 头寸价值 以标的物计价
		ForeignNotional      float64 `json:"foreignNotional"`      // 头寸价值 以货币计价
		LiquidationPrice     float64 `json:"liquidationPrice"`     // 强平价格
		BankruptPrice        float64 `json:"bankruptPrice"`        // 破产价格 即头寸无价值
		MarkPrice            float64 `json:"markPrice"`            // 标记价格 用于计算平仓价格等
		MarkValue            float64 `json:"markValue"`            // 标记指 ForeignNotional * 10000 * 10000
		CurrentQty           float64 `json:"currentQty"`           // 持仓量 <0 (做空) >0(做多)
		Timestamp            string  `json:"timestamp"`            // 时间
		LastPrice            float64 `json:"lastPrice"`            // 最新价格
	}

	OrderMsg struct {
		Table  string  `json:"table"`
		Action string  `json:"Action"`
		Data   []Order `json:"data"`
	}
	Order struct {
		Account               float64 `json:"account"`
		OrderID               string  `json:"orderID"`
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
		MultiLegReportingType string  `json:"multiLegReportingType"`
		Text                  string  `json:"text"`
		TransactTime          string  `json:"transactTime"`
		Timestamp             string  `json:"timestamp"`
	}

	ExecutionMsg struct {
		Table  string      `json:"table"`
		Action string      `json:"Action"`
		Data   []Execution `json:"data"`
	}

	Execution struct {
		ExecID                string  `json:"execID"`
		OrderID               string  `json:"orderID"`
		ClOrdID               string  `json:"clOrdID"`
		ClOrdLinkID           string  `json:"clOrdLinkID"`
		Account               float64 `json:"account"`
		Symbol                string  `json:"symbol"`
		Side                  string  `json:"side"`
		LastQty               float64 `json:"lastQty"`
		LastPx                float64 `json:"lastPx"`
		UnderlyingLastPx      float64 `json:"underlyingLastPx"`
		LastMkt               string  `json:"lastMkt"`
		LastLiquidityInd      string  `json:"lastLiquidityInd"`
		SimpleOrderQty        float64 `json:"simpleOrderQty"`
		OrderQty              float64 `json:"orderQty"`
		Price                 float64 `json:"price"`
		DisplayQty            float64 `json:"displayQty"`
		StopPx                float64 `json:"stopPx"`
		PegOffsetValue        float64 `json:"pegOffsetValue"`
		PegPriceType          string  `json:"pegPriceType"`
		Currency              string  `json:"currency"`
		SettlCurrency         string  `json:"settlCurrency"`
		ExecType              string  `json:"execType"`
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
		Commission            float64 `json:"commission"`
		TradePublishIndicator string  `json:"tradePublishIndicator"`
		MultiLegReportingType string  `json:"multiLegReportingType"`
		Text                  string  `json:"text"`
		TrdMatchID            string  `json:"trdMatchID"`
		ExecCost              float64 `json:"execCost"`
		ExecComm              float64 `json:"execComm"`
		HomeNotional          float64 `json:"homeNotional"`
		ForeignNotional       float64 `json:"foreignNotional"`
		TransactTime          string  `json:"transactTime"`
		Timestamp             string  `json:"timestamp"`
	}
)

func init() {
	operate = make(chan Operate, 1)

	go func() {
		for {
			op := <-operate
			switch op.Action {
			case "create":
				if err := createOrder(op.Params); err != nil {
					log.Info(err)
				}
			case "amend":
				if err := amendOrder(op.Params); err != nil {
					log.Info(err)
				}
			case "cancel":
				if err := cancelOrder(op.Params); err != nil {
					log.Info(err)
				}
			default:
				log.Info("not supported action")
			}
		}
	}()
}

func dispatch(msg []byte) (err error) {

	//log.Debug(string(msg))
	message := string(msg)

	if message == "pong" {
		return handlePing(message)
	}

	topic := gjson.GetBytes(msg, "table")

	switch topic.String() {
	case "orderBook10":
		return handleOrderBook10(msg)
	case "execution":
		return handleExecution(msg)
	case "position":
		return handlePosition(msg)
	case "order":
		return handleOrder(msg)
	default:
		return
	}

	return
}

// ping
func handlePing(msg string) (err error) {
	log.Debug(msg)
	// 移仓
	for _, v := range order {

		if v.OrdStatus != "New" {
			continue
		}

		if v.Side == "Buy" && v.Price >= orderBook10[v.Symbol].Bids[Conf.Trading.Range-1][0] {
			continue
		}

		if v.Side == "Sell" && v.Price <= orderBook10[v.Symbol].Asks[Conf.Trading.Range-1][0] {
			continue
		}
		go func(v Order) {
			action := "cancel"
			params := make(map[string]interface{})
			params["orderID"] = v.OrderID
			operate <- Operate{
				action,
				params,
			}
			log.Infof("%s order %s to be canceled, price is %v, qty is %v", v.Side, v.OrderID, v.Price, v.OrderQty)
		}(v)
	}

	// 超出最大持仓量
	for k, v := range position {

		toBuy := true
		toSell := true
		for _, vv := range order {
			if k != vv.Symbol {
				continue
			}
			if vv.Side == "Buy" && vv.Price == orderBook10[vv.Symbol].Bids[0][0] && vv.OrdStatus == "New" {
				log.Debug(vv)
				toBuy = false
				break
			}
		}

		for _, vv := range order {
			if k != vv.Symbol {
				continue
			}
			if vv.Side == "Sell" && vv.Price == orderBook10[vv.Symbol].Asks[0][0] && vv.OrdStatus == "New" {
				log.Debug(vv)
				toSell = false
				break
			}
		}

		log.Infof("CurrentQty: %v", v.CurrentQty)
		if math.Abs(v.CurrentQty) > Conf.Trading.MaxHoldQty {
			go func() {
				action := "create"
				params := make(map[string]interface{})
				params["symbol"] = k
				params["orderQty"] = Conf.Trading.UnitQty * 2
				params["side"] = "Buy"
				params["price"] = orderBook10[k].Bids[0][0]
				if v.CurrentQty > 0 {
					params["side"] = "Sell"
					params["price"] = orderBook10[k].Asks[0][0]
				}
				if (v.CurrentQty > 0 && toSell) || (v.CurrentQty < 0 && toBuy) {
					operate <- Operate{
						action,
						params,
					}
					log.Infof("%s order to be created at %v, qty is %v", k, params["price"], params["orderQty"])
				}
			}()
		}
	}

	// 填价
	for _, s := range Conf.Trading.Symbol {
		toBuy := true
		toSell := true
		for _, v := range order {
			if s != v.Symbol {
				continue
			}
			if v.Side == "Buy" && v.Price == orderBook10[v.Symbol].Bids[0][0] && v.OrdStatus == "New" {
				log.Debug(v)
				toBuy = false
				break
			}
		}

		for _, v := range order {
			if s != v.Symbol {
				continue
			}
			if v.Side == "Sell" && v.Price == orderBook10[v.Symbol].Asks[0][0] && v.OrdStatus == "New" {
				log.Debug(v)
				toSell = false
				break
			}
		}

		if math.Abs(position[s].CurrentQty) > Conf.Trading.MaxHoldQty {
			log.Info("reach max hold Qty, stop create order")
			break
		}

		log.Infof("toBuy: %v, toSell: %v", toBuy, toSell)
		go func(toBuy, toSell bool) {
			action := "create"
			if toBuy {
				params := make(map[string]interface{})
				params["symbol"] = s
				params["orderQty"] = Conf.Trading.UnitQty
				params["side"] = "Buy"
				params["price"] = orderBook10[s].Bids[0][0]
				operate <- Operate{
					action,
					params,
				}
				log.Infof("%s order to be create at %v, qty is %v", params["side"], params["price"], params["orderQty"])
			}
			if toSell {
				params2 := make(map[string]interface{})
				params2["symbol"] = s
				params2["orderQty"] = Conf.Trading.UnitQty
				params2["side"] = "Sell"
				params2["price"] = orderBook10[s].Asks[0][0]
				operate <- Operate{
					action,
					params2,
				}
				log.Infof("%s order to be create at %v, qty is %v", params2["side"], params2["price"], params2["orderQty"])
			}
		}(toBuy, toSell)
	}

	return
}

// 10档报价
func handleOrderBook10(msg []byte) (err error) {
	obm := &OrderBook10Msg{}
	if err = json.Unmarshal(msg, obm); err != nil {
		return
	}

	if obm.Action == "partial" || obm.Action == "update" {
		orderBook10 = make(map[string]OrderBook10)
		for _, order := range obm.Data {
			orderBook10[order.Symbol] = order
			log.Info("---")
			log.Infof("range: %v ~ %v", order.Asks[Conf.Trading.Range-1][0], order.Bids[Conf.Trading.Range-1][0])
			log.Infof("Asks: %v\t%v\t%v\t%v\t%v", order.Asks[0], order.Asks[1], order.Asks[2], order.Asks[3], order.Asks[4])
			log.Infof("Bids: %v\t%v\t%v\t%v\t%v", order.Bids[0], order.Bids[1], order.Bids[2], order.Bids[3], order.Bids[4])
			log.Info("---")
		}
		return
	}
	return
}

// 订单成交
func handleExecution(msg []byte) (err error) {
	em := &ExecutionMsg{}
	if err = json.Unmarshal(msg, em); err != nil {
		log.Info(err)
		return
	}

	if em.Action == "insert" {
		for _, v := range em.Data {
			if v.OrdStatus == "Filled" {
				log.Infof("%s order %s filled at %v, qty is %v", v.Side, v.OrderID, v.Price, v.CumQty)
				params := make(map[string]interface{})
				params["symbol"] = v.Symbol

				params["side"] = "Buy"
				if v.Side == "Buy" {
					params["side"] = "Sell"
				}

				params["orderQty"] = v.CumQty

				spread := Conf.Trading.Spread
				if v.Side == "Sell" {
					spread *= -1
				}
				params["price"] = v.Price + spread*Conf.Trading.PriceUint

				if params["side"] == "Buy" && params["price"].(float64) > orderBook10[v.Symbol].Bids[0][0] {
					params["price"] = orderBook10[v.Symbol].Asks[0][0]
				}

				if params["side"] == "Sell" && params["price"].(float64) < orderBook10[v.Symbol].Asks[0][0] {
					params["price"] = orderBook10[v.Symbol].Bids[0][0]
				}

				operate <- Operate{
					"create",
					params,
				}
				log.Infof("%s order to be created at %v, qty is %v", params["side"], params["price"], params["orderQty"])
			}
		}
	}
	return
}

// 头寸
func handlePosition(msg []byte) (err error) {
	pm := &PositionMsg{}
	if err = json.Unmarshal(msg, pm); err != nil {
		log.Info(err)
		return
	}

	position = make(map[string]Position)
	if pm.Action == "partial" {
		for _, p := range pm.Data {
			position[p.Symbol] = p
		}
		log.Debug("partial position", position["XBTUSD"])
		return
	}

	update := func(source, patch Position) Position {
		s := source
		v := reflect.ValueOf(s)
		vv := reflect.ValueOf(patch)
		v_elem := reflect.ValueOf(&s).Elem()

		for i := 0; i < v.NumField(); i++ {
			//f := v.Field(i)
			ff := vv.Field(i)
			if reflect.DeepEqual(ff.Interface(), reflect.Zero(ff.Type()).Interface()) {
				continue
			}
			v_elem.Field(i).Set(ff)
		}
		return s
	}

	if pm.Action == "update" {
		for _, p := range pm.Data {
			position[p.Symbol] = update(position[p.Symbol], p)
		}
		log.Debug("update position", position["XBTUSD"])
		return
	}

	return
}

// 未成交订单
func handleOrder(msg []byte) (err error) {
	log.Debug(string(msg))
	om := &OrderMsg{}
	if err = json.Unmarshal(msg, om); err != nil {
		return
	}

	if om.Action == "partial" {
		order = om.Data
		log.Debug(order)
		log.Debug(len(order))
		return
	}

	if om.Action == "insert" {
		order = append(order, om.Data...)
		log.Debug(order)
		log.Debug(len(order))
		return
	}

	update := func(source, patch Order) Order {
		s := source
		v := reflect.ValueOf(s)
		vv := reflect.ValueOf(patch)
		v_elem := reflect.ValueOf(&s).Elem()

		for i := 0; i < v.NumField(); i++ {
			//f := v.Field(i)
			ff := vv.Field(i)
			if reflect.DeepEqual(ff.Interface(), reflect.Zero(ff.Type()).Interface()) {
				continue
			}
			v_elem.Field(i).Set(ff)
		}
		return s
	}

	if om.Action == "update" {
		for _, v := range om.Data {
			for kk, vv := range order {
				if v.OrderID == vv.OrderID {
					order[kk] = update(vv, v)
				}
			}
		}
		for k, v := range order {
			if v.OrdStatus == "Canceled" {
				order = append(order[:k], order[k+1:]...)
			}
		}
		return
	}

	return
}

// 下单
func createOrder(params map[string]interface{}) error {
	or := OrderResponse{}
	ep := Endpoint{
		"POST",
		"/order",
		Conf.RestConfig,
		Conf.AuthConfig,
		params,
		&or,
		nil,
	}
	return ep.Do()
}

// 改单
func amendOrder(params map[string]interface{}) error {
	or := OrderResponse{}
	ep := Endpoint{
		"PUT",
		"/order",
		Conf.RestConfig,
		Conf.AuthConfig,
		params,
		&or,
		nil,
	}
	return ep.Do()
}

// 取消订单
func cancelOrder(params map[string]interface{}) error {
	ep := Endpoint{
		"DELETE",
		"/order",
		Conf.RestConfig,
		Conf.AuthConfig,
		params,
		nil,
		nil,
	}
	return ep.Do()
}

// 设置杠杆率
func setLeverage(params map[string]interface{}) error {
	ep := Endpoint{
		"POST",
		"/position/leverage",
		Conf.RestConfig,
		Conf.AuthConfig,
		params,
		nil,
		nil,
	}
	return ep.Do()
}
