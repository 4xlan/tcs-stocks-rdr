package tcsrdrserver

type TCSPriceData interface {
	getLastPrice(figi *string)
}

type TCSPriceResponse struct {
	TrackingId string `json:"trackingId"`
	Status     string `json:"status"`
	Payload    OrderbookPayload
}

type OrderbookPayload struct {
	Figi              string  `json:"figi"`
	Depth             float64 `json:"depth"`
	TradeStatus       string  `json:"tradeStatus"`
	MinPriceIncrement float64 `json:"minPriceIncrement"`
	FaceValue         float64 `json:"faceValue"`
	LastPrice         float64 `json:"lastPrice"`
	ClosePrice        float64 `json:"closePrice"`
	LimitUp           float64 `json:"limitUp"`
	LimitDown         float64 `json:"limitDown"`
	Bids              []AsksBids
	Asks              []AsksBids
}

type AsksBids struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
