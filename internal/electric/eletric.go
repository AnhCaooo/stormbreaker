package electric

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnhCaooo/stormbreaker/internal/logger"
	"go.uber.org/zap"
)

const (
	BASE_URL   = "https://oomi.fi/wp-json"
	SPOT_PRICE = "spot-price"
	GET_V1     = "v1/get"
)

// basic logic: fetch data from oomi.fi  once per day and store data to database. Then return this value to client
// advanced logic: before fetch data from oomi.fi, check from database if the query exists or not. If exists, get from db, otherwise call to Oomi.fi
// Goal: this will prevent someone tries to use this service spam Oomi.fi
// TODO: bring the commented code back in next implementation and remove current one
// func FetchSpotPrice(requestParameters PriceRequest) (responseData PriceResponse, err error) {
func FetchSpotPrice() (responseData *PriceResponse, err error) {
	// TODO: bring the commented code back in next implementation and remove current one
	externalUrl, err := formatRequestParameters()
	// externalUrl, err := formatRequestParameters(requestParameters)
	if err != nil {
		logger.Logger.Error("failed to format Url", zap.Error(err))
		return nil, err
	}

	// Make HTTP request to the external source
	resp, err := http.Get(externalUrl)
	if err != nil {
		logger.Logger.Error("failed to fetch data from external source (Oomi)", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil { // Parse []byte to the go struct pointer
		logger.Logger.Error("can not unmarshal JSON", zap.Error(err))
		return nil, err
	}

	return

}

// TODO: bring the commented code back in next implementation and remove current one
func formatRequestParameters() (endPoint string, err error) {
	// func formatRequestParameters(requestParameters PriceRequest) (endPoint string, err error) {
	url := fmt.Sprintf("%s/%s/%s", BASE_URL, SPOT_PRICE, GET_V1)

	logger.Logger.Info("request url", zap.String("url", url)) // todo: remove this log
	return url, nil
}
