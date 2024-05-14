package electric

// Represents single electric data at specific time
type Data struct {
	OriginalTime string  `json:"orig_time"`
	Time         string  `json:"time"`
	Price        float64 `json:"price"`
	VatFactor    float64 `json:"vat_factor"`
	IsToday      bool    `json:"isToday"`
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
