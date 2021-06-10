package io

import (
	consts "github.com/ezeeredisservices/constants"
	"github.com/ezeeredisservices/logger"
)

type Location struct {
	Key     string    `json:"key"`
	Address []Address `json:"address"`
}

func (loc Location) IsValid() bool {
	if loc.Key == "" || loc.Key == "NA" {
		logger.ErrorLogger.Println(consts.InvalidKey.Error())
		logger.ErrorLogger.Println(loc.Key)
		return false
	}
	if loc.Address == nil || len(loc.Address) == 0 {
		logger.ErrorLogger.Println(consts.InvalidKey.Error())
		logger.ErrorLogger.Println(loc.Address)
		return false
	}
	for _, address := range loc.Address {
		if !address.IsValid() {
			return false
		}
	}
	return true
}
