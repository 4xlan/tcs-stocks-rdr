package tcsrdrserver

type TCSPortfolioData interface {
	getData(response *TCSPortfolioResponse)
}

type TCSPortfolioResponse struct {
	TrackingId       string `json:"trackingID"`
	PortfolioPayload `json:"payload"`
	Status           string `json:"status"`
}

type PortfolioPayload struct {
	Positions []Position
}

type Position struct {
	Figi                 string  `json:"figi"`
	Ticker               string  `json:"ticker"`
	Isin                 string  `json:"isin"`
	InstrumentType       string  `json:"instrumentType"`
	Balance              float64 `json:"balance"`
	Lots                 float64 `json:"lots"`
	ExpectedYield        CurVal
	AveragePositionPrice CurVal
	Name                 string `json:"name"`
}

type CurVal struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}
