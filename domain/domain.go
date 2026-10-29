// Package domain chứa hợp đồng sự kiện dùng chung giữa CDK stack, Lambda và ECS executor.
package domain

// EventBridge sources & detail-types.
const (
	SourceForecast = "solar.forecast"
	SourcePosition = "solar.position"
	SourceExchange = "solar.exchange"

	ForecastGenerated = "ForecastGenerated"
	DeficitDetected   = "DeficitDetected"
	SurplusDetected   = "SurplusDetected"
	OrderFilled       = "OrderFilled"
	OrderRejected     = "OrderRejected"
)

// Parameter Store của project trading. Mỗi project × env là một account riêng nên path không chứa env.
//
//	/solar-trading/zones/<zone>/<param>   cấu hình nghiệp vụ theo khu vực (zone là dữ liệu, không phải hạ tầng)
//	/solar-trading/exchange/<param>       cấu hình kết nối sàn
const (
	ParamRoot         = "/solar-trading"
	ExchangeParamPath = ParamRoot + "/exchange"

	ParamDeficitThreshold = "deficit-threshold-kwh"
	ParamMaxBuyPrice      = "max-buy-price"
	ParamMinSellPrice     = "min-sell-price"
	ParamPlantCapacity    = "own-plant-capacity-kwh"
	ParamExchangeEndpoint = "endpoint"
	ParamExchangeMode     = "mode" // mock | live
)

func ZoneParamPath(zone string) string { return ParamRoot + "/zones/" + zone }

// Forecast do AI forecast engine (EC2) phát ra cho từng khách hàng, từng khung 15 phút.
type Forecast struct {
	CustomerID     string  `json:"customerId"`
	Zone           string  `json:"zone"`
	Interval       string  `json:"interval"` // RFC3339, đầu khung 15 phút
	ProductionKwh  float64 `json:"forecastProductionKwh"`
	ConsumptionKwh float64 `json:"forecastConsumptionKwh"`
	Model          string  `json:"model"`
}

// TradeIntent là phần cần mua/bán trên sàn sau khi đã trừ phần nhà máy của mình cấp.
type TradeIntent struct {
	IntentID    string  `json:"intentId"` // zone#customer#interval - dùng làm idempotency key
	CustomerID  string  `json:"customerId"`
	Zone        string  `json:"zone"`
	Interval    string  `json:"interval"`
	Side        string  `json:"side"` // BUY | SELL
	QuantityKwh float64 `json:"quantityKwh"`
	LimitPrice  float64 `json:"limitPrice"` // VND/kWh
}

type OrderResult struct {
	IntentID  string  `json:"intentId"`
	OrderID   string  `json:"orderId"`
	Status    string  `json:"status"` // FILLED | REJECTED
	FilledKwh float64 `json:"filledKwh"`
	Price     float64 `json:"price"`
}
