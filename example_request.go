package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (api *Backend) ExampleGetRequest(id string, loggerCallBack func(string)) (responseData string, err error) {

	responseRaw, _, err := api.call(http.MethodGet, api.config.BaseUrl+"/api/example/"+id, bytes.NewBuffer([]byte{}))
	if err != nil {
		printLog(loggerCallBack, err.Error())
		return "", err
	}

	return string(responseRaw), nil
}

func (api *Backend) ExamplePostRequest(id string, loggerCallBack func(string)) (responseData string, err error) {

	jsonRequestData, err := json.Marshal(struct {
		ID string `json:"id"`
	}{
		ID: id,
	})

	if err != nil {
		printLog(loggerCallBack, err.Error())
		return "", err
	}

	responseRaw, _, err := api.call(http.MethodPost, api.config.BaseUrl+"/api/example", bytes.NewBuffer(jsonRequestData))
	if err != nil {
		printLog(loggerCallBack, err.Error())
		return "", err
	}

	return string(responseRaw), nil
}
