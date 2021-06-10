package io

import (
	"github.com/ezeeredisservices/logger"
)

type Address struct {
	Code      string `json:"code"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (address Address) IsValid() bool {
	if address.Code == "" || address.Code == "NA" {
		logger.ErrorLogger.Println("Invalid Address Code " + address.Code)
		return false
	}
	if address.Latitude == "" || address.Latitude == "NA" {
		logger.ErrorLogger.Println("Invalid Latitude " + address.Latitude)
		return false
	}
	if address.Longitude == "" || address.Longitude == "NA" {
		logger.ErrorLogger.Println("Invalid Longitude " + address.Longitude)
		return false
	}
	return true
}
