package electric

import "time"

type TimelyPrice struct {
	OriginalTime time.Time `json:"orig_time"`
	Time         time.Time `json:"time"`
	Price        float64   `json:"price"`
	VatFactor    float64   `json:"vat_factor"`
	IsToday      bool      `json:"isToday"`
}

type PriceGroup struct {
	Name string        `json:"name"`
	Data []TimelyPrice `json:"data"`
}

// Group represents 'hour', 'day', week', 'month', 'year'
type PriceData struct {
	Group  string       `json:"group"`
	Series []PriceGroup `json:"series"`
}

type PriceResponse struct {
	Status string    `json:"status"`
	Data   PriceData `json:"data"`
}

// vatIncluded will affect when it is different than 0
type PriceRequestParameters struct {
	StartDate         time.Time `json:"starttime"`
	EndDate           time.Time `json:"endtime"`
	Marginal          float64   `json:"margin"`
	Group             string    `json:"group"`
	VatIncluded       int32     `json:"include_vat"`
	CompareToLastYear int32     `json:"compare_to_last_year"`
}
