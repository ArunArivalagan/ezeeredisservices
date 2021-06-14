package communicator

import (
	"io/ioutil"
	"net/http"

	entity "github.com/ezeeredisservices/io"
	"github.com/ezeeredisservices/logger"
	"github.com/ezeeredisservices/utils"
)

func GetGeoLocations(fromDate, toDate, filterType string) []entity.Location {
	response, err := http.Get("http://localhost:8010/geoservices/geo/geo/location/address?fromDate=" + fromDate + "&toDate=" + toDate + "&filterType=" + filterType)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	georesponse, err := utils.UnMarshalBinaryArrayGeoResponse(responseData)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	var locations []entity.Location
	locations = append(locations, georesponse.Locations...)

	return locations
}
