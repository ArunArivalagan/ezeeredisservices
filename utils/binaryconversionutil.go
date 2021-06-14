package utils

import (
	"encoding/json"

	entity "github.com/ezeeredisservices/io"
)

func MarshalBinaryArray(address []entity.Address) ([]byte, error) {
	bytes, err := json.Marshal(address)
	return bytes, err
}

func UnMarshalBinaryArray(data []byte) ([]entity.Address, error) {
	var address []entity.Address
	err := json.Unmarshal(data, &address)
	return address, err
}

func UnMarshalBinaryArrayLocation(data []byte) ([]entity.Location, error) {
	var locations []entity.Location
	err := json.Unmarshal(data, &locations)
	return locations, err
}

func UnMarshalBinaryArrayGeoResponse(data []byte) (entity.GeoResponse, error) {
	var georespo entity.GeoResponse
	err := json.Unmarshal(data, &georespo)
	return georespo, err
}
