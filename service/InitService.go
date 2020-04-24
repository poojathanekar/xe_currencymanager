package service

import (
	"log"
	"xe_currencymanager/config"
	"xe_currencymanager/db"
	"xe_currencymanager/model"

	"encoding/json"
	"io/ioutil"

	logger "github.com/sirupsen/logrus"
	// "log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq" // here
	// "github.com/robfig/cron"
)

// Call for currencyConverter
func InitService() {
	var wg sync.WaitGroup
	startTime := time.Now()
	currencies := config.GetStringSlice("currencies")

	for _, fromCurrency := range currencies {
		wg.Add(1)
		go callService(fromCurrency, &wg)
	}

	endTime := time.Since(startTime)
	wg.Wait()

	logger.WithField(" Total time taken for execution: ", endTime).Info("Execution Time")
}

// callService : get data from XE API and updates to db
func callService(from string, wg *sync.WaitGroup) {
	responseData, apiErr := getXeAPIData(from, wg)
	if apiErr != nil {
		logger.WithField("get XE API data error:", apiErr.Error()).Info("Get API data failed")
		log.Panic(apiErr)
	}
	// prepared query values with parametrs
	queryValues := make([]string, 0, len(responseData.To))
	queryParams := make([]interface{}, 0, len(responseData.To)*5)
	for _, toValue := range responseData.To {
		queryValues = append(queryValues, "(?, ?, ?, ?, ?)")
		queryParams = append(queryParams, responseData.Amount, toValue.Mid, responseData.From, toValue.Quotecurrency, responseData.Timestamp)
	}
	// update responsedata to db
	dbErr := db.UpdateResponseData(responseData, queryValues, queryParams)
	if dbErr != nil {
		logger.WithField("update data to database error:", apiErr.Error()).Info("Updating to database failed")
		log.Panic(apiErr)
	}
	wg.Done()

}

// getXeApiData : convert to other exchange format which hit XE API
func getXeAPIData(from string, wg *sync.WaitGroup) (model.XEResponse, error) {
	url := config.GoDotEnvVariable("URL")
	username := config.GoDotEnvVariable("USERNAME")
	password := config.GoDotEnvVariable("PASSWORD")
	xeResponseData := model.XEResponse{}

	client := &http.Client{}
	req, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return xeResponseData, requestErr
	}

	q := req.URL.Query()
	q.Add("to", "*")
	q.Add("from", from)

	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(username, password)

	client = &http.Client{
		Timeout: time.Second * time.Duration(1500),
	}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		logger.WithField("error from api", clientErr.Error()).Error("Get Request Failed")
		return xeResponseData, clientErr
	}
	if resp.StatusCode != 200 {
		logger.WithField("error from api", resp).Error("Get Request Failed")
		return xeResponseData, clientErr
	}

	bodyText, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logger.WithField("read http response error", readErr.Error()).Error("read response failed")
		return xeResponseData, readErr
	}

	unmarshalErr := json.Unmarshal(bodyText, &xeResponseData)
	if unmarshalErr != nil {
		logger.WithField("unmarshal error", unmarshalErr.Error()).Error("Unmarshal Failed")
		return xeResponseData, unmarshalErr
	}

	defer resp.Body.Close()
	// logger.WithField("Total time taken for get api execution:", requestElapsed).Info("Exec Time")
	// call for update currency rate to all exchange rates
	return xeResponseData, nil
}

//TestXEapiCallService : test for XE api call function
