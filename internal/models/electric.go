// AnhCao 2024
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
	TimeUTC      string  `json:"time_utc" example:"2024-12-08 22:00:00"`  // timestamp in UTC format
	OriginalTime string  `json:"orig_time" example:"2024-12-09 00:00:00"` // the current time where server is located
	Time         string  `json:"time" example:"2024-12-09 00:00:00"`      // the current time.
	Price        float64 `json:"price" example:"2.47"`                    // the price of specified time range
	VatFactor    float64 `json:"vat_factor" example:"1.255"`              // amount of VAT that applies to electric price.
	IsToday      bool    `json:"isToday" example:"false"`                 // IsToday indicates whether the current time is today or not
	IncludeVat   string  `json:"includeVat" example:"1" enums:"0,1"`      // IncludeVat is legacy property that return string value and value "0" means no VAT included and string "1" is included
}

// Represents a series of electric data with the name of unit (ex: c/kwh)
type PriceSeries struct {
	Name string `json:"name" example:"c/kWh"` // unit of electric price
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
	StartDate         string `json:"starttime" example:"2024-12-11"` // StartDate has to be in this format "YYYY-MM-DD"
	EndDate           string `json:"endtime" example:"2024-12-31"`   // EndDate has to be in this format "YYYY-MM-DD"
	Group             string `json:"group" example:"hour" enums:"hour,day,week,month,year"`
	CompareToLastYear int32  `json:"compare_to_last_year" example:"0" enums:"0,1"` // CompareToLastYear is allowed to equal to "0" and "1"
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

// PriceSettings represents the schema for the PriceSettings collection
type PriceSettings struct {
	UserID      string  `bson:"user_id" json:"user_id" example:"123456789"`      // id of the user. When sends as request, the clients (web, mobile) does not need to provide `user_id` because the service will read through `access_token`.
	VatIncluded bool    `bson:"vat_included" json:"vat_included" example:"true"` // indicates whether tax is included to price stats or not
	Marginal    float64 `bson:"margin" json:"margin" example:"0.59"`             // amount of margin applied to price stats
}
