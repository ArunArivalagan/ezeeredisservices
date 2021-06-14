package io

type GeoResponse struct {
	Status    int        `json:"status"`
	Datetime  string     `json:"datetime"`
	Locations []Location `json:"data"`
}
