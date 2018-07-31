package boot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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

	Position struct {
		Table  string `json:"table"`
		Action string `json:"Action"`
		Data   []struct {
			Account              float64 `json:"account"`
			Symbol               string  `json:"symbol"`
			InitMarginReq        float64 `json:"initMarginReq"`        // 初始保证金
			MaintMarginReq       float64 `json:"maintMarginReq"`       // 维持保证金
			Leverage             float64 `json:"leverage"`             // 杠杆率
			RiskLimit            float64 `json:"riskLimit"`            //	风险限额
			crossMargin          bool    `json:"crossMargin"`          // 全仓保证金(false)/逐仓保证金(true)
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
			LastPrice            string  `json:"lastPrice"`            // 最新价格
		} `json:"data"`
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
		return handleExecution(msg)
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

func handleExecution(msg []byte) (err error) {
	log.Debug(string(msg))
	return
}

func handlePosition(msg []byte) (err error) {
	log.Debug(string(msg))
	return
}

func handleOrder(msg []byte) (err error) {
	log.Debug(string(msg))
	return
}
