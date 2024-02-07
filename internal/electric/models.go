package electric

type Data struct {
	OriginalTime string  `json:"orig_time"`
	Time         string  `json:"time"`
	Price        float64 `json:"price"`
	VatFactor    float64 `json:"vat_factor"`
	IsToday      bool    `json:"isToday"`
}

type PriceSeries struct {
	Name string `json:"name"`
	Data []Data `json:"data"`
}

// Group represents 'hour', 'day', week', 'month', 'year'
type PriceData struct {
	Group  string        `json:"group"`
	Series []PriceSeries `json:"series"`
}

type PriceResponse struct {
	Data   PriceData `json:"data"`
	Status string    `json:"status"`
}

// vatIncluded will affect when it is different than 0
type PriceRequest struct {
	StartDate         string  `json:"starttime"`
	EndDate           string  `json:"endtime"`
	Marginal          float64 `json:"margin"`
	Group             string  `json:"group"`
	VatIncluded       int32   `json:"include_vat"`
	CompareToLastYear int32   `json:"compare_to_last_year"`
}
