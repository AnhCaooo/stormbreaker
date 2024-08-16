package models

const (
	BASE_URL     string = "https://oomi.fi/wp-json"
	SPOT_PRICE   string = "spot-price"
	GET_V1       string = "v1/get"
	CLIENT_ERROR string = "client"
	SERVER_ERROR string = "server"
)

// Represents single electric data at specific time
type Data struct {
	TimeUTC      string      `json:"time_utc"`
	OriginalTime string      `json:"orig_time"`
	Time         string      `json:"time"`
	Price        float64     `json:"price"`
	VatFactor    float64     `json:"vat_factor"`
	IsToday      bool        `json:"isToday"`
	IncludeVat   interface{} `json:"includeVat"` // IncludeVat is legacy interface which "false" means no VAT included and string "1" is included
}

// Represents a series of electric data with the name of unit (ex: c/kwh)
type PriceSeries struct {
	Name string `json:"name"`
	Data []Data `json:"data"`
}

// Represent a series data of electric price in targeting group
type PriceData struct {
	Group  string        `json:"group"` // Group represents 'hour', 'day', 'week', 'month', 'year'
	Series []PriceSeries `json:"series"`
}

// Represent as a struct of response data when fetching electric price data
type PriceResponse struct {
	Data   PriceData `json:"data"`
	Status string    `json:"status"`
}

// Represents as request body when client (web, mobile, backend service) call to get market price in specific time range
type PriceRequest struct {
	StartDate         string  `json:"starttime"` // StartDate has to be in this format "YYYY-MM-DD"
	EndDate           string  `json:"endtime"`   // EndDate has to be in this format "YYYY-MM-DD"
	Marginal          float64 `json:"margin"`    // Marginal is not allowed to be empty, it is ok to equal to "0"
	Group             string  `json:"group"`
	VatIncluded       int32   `json:"include_vat"`          // VatIncluded is allowed to equal to "0" and "1"
	CompareToLastYear int32   `json:"compare_to_last_year"` // CompareToLastYear is allowed to equal to "0" and "1"
}

// Represents a struct of today and tomorrow exchange price
type TodayTomorrowPrice struct {
	Today    DailyPrice `json:"today"`
	Tomorrow DailyPrice `json:"tomorrow"`
}

// Represents a struct of daily price and bool flag to indicate does tomorrow's price available or not
type DailyPrice struct {
	Available bool        `json:"available"`
	Prices    PriceSeries `json:"prices"`
}
