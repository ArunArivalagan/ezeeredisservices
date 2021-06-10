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
