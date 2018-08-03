package boot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"reflect"
	//"strings"
)

var (
	orderBook10 map[string]OrderBook10
	position    map[string]Position
	order       []Order
)

type (
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
)

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

func handlePing(msg string) (err error) {
	log.Debug(msg)
	return
}

func handleOrderBook10(msg []byte) (err error) {
	obm := &OrderBook10Msg{}
	if err = json.Unmarshal(msg, obm); err != nil {
		return
	}

	if obm.Action == "partial" || obm.Action == "update" {
		orderBook10 = make(map[string]OrderBook10)
		for _, order := range obm.Data {
			orderBook10[order.Symbol] = order
		}
		return
	}
	return
}

func handleExecution(msg []byte) (err error) {
	log.Debug(string(msg))
	return
}

func handlePosition(msg []byte) (err error) {
	pm := &PositionMsg{}
	if err = json.Unmarshal(msg, pm); err != nil {
		return
	}

	if pm.Action == "partial" {
		position = make(map[string]Position)
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
			//log.Debug(p)
			/*
				pp := position[p.Symbol]
				v := reflect.ValueOf(pp)
				vv := reflect.ValueOf(p)
				v_elem := reflect.ValueOf(&pp).Elem()

				for i := 0; i < v.NumField(); i++ {
					//f := v.Field(i)
					ff := vv.Field(i)
					if reflect.DeepEqual(ff.Interface(), reflect.Zero(ff.Type()).Interface()) {
						continue
					}
					v_elem.Field(i).Set(ff)
				}
				position[p.Symbol] = pp
			*/
			position[p.Symbol] = update(position[p.Symbol], p)

		}
		log.Debug("update position", position["XBTUSD"])
		return
	}

	return
}

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
		log.Debug(order)
		log.Debug(len(order))
		return
	}

	return
}
