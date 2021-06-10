package io

import (
	"encoding/json"
	"time"
)

type Response struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
	Datetime  string `json:"datetime"`
	Data      string `json:"data"`
}

func Success(dataStruct interface{}) Response {
	var response Response
	var data []byte
	response.Status = 1
	response.Datetime = time.Now().Format("2006-01-02 15:04:05")
	if dataStruct != nil {
		data, _ = json.Marshal(dataStruct)
	}
	response.Data = string(data)
	return response
}

func Failure(errorCode int, errorDesc string) Response {
	var response Response
	response.Status = 0
	response.Datetime = time.Now().Format("2006-01-02 15:04:05")
	response.ErrorCode = errorCode
	response.ErrorDesc = errorDesc
	return response
}
