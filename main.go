package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezeeredisservices/communicator"
	cnf "github.com/ezeeredisservices/config"
	consts "github.com/ezeeredisservices/constants"
	entity "github.com/ezeeredisservices/io"
	logger "github.com/ezeeredisservices/logger"
	util "github.com/ezeeredisservices/utils"
	"github.com/gorilla/mux"
)

var cachePrefix = "LOC_"

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(verifyAccessToken)
	router.HandleFunc("/ezeeredisservices", welcomePage)

	router.HandleFunc("/ezeeredisservices/locations/add", addLocations).Methods("POST")
	router.HandleFunc("/ezeeredisservices/locations/nearest/{key}", geoNearestLocations).Methods("GET")

	/**It's not for geo util*/
	router.HandleFunc("/ezeeredisservices/locations/{key}", getLocations).Methods("GET")
	router.HandleFunc("/ezeeredisservices/locations", getAllLocations).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func addLocations(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("fromDate")
	toDate := r.URL.Query().Get("toDate")
	filterType := r.URL.Query().Get("filterType")

	var err error

	locations := communicator.GetGeoLocations(fromDate, toDate, filterType)

	for _, location := range locations {
		err = cnf.GeoAdd(cachePrefix+location.Key, location.Address)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UpdateFailed), consts.UpdateFailed.Error()))
			return
		}
	}
	json.NewEncoder(w).Encode(entity.Success(nil))
}

func geoNearestLocations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	radius := r.URL.Query().Get("radius")
	count := r.URL.Query().Get("count")

	var err error
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.InvalidRequest), consts.InvalidRequest.Error()))
		return
	}
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.InvalidRequest), consts.InvalidRequest.Error()))
		return
	}
	rad, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.InvalidRequest), consts.InvalidRequest.Error()))
		return
	}
	cnt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.InvalidRequest), consts.InvalidRequest.Error()))
		return
	}

	cacheData, err := cnf.GeoNearestLocations(cachePrefix+key, lat, lon, rad, int(cnt))
	logger.InfoLogger.Println(cacheData)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}

	json.NewEncoder(w).Encode(entity.Success(cacheData))
}

func getLocations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" || key == "NA" {
		json.NewEncoder(w).Encode(entity.Failure(int(consts.InvalidKey), consts.InvalidKey.Error()))
		return
	}
	logger.InfoLogger.Println(key)

	var err error
	cacheData, err := cnf.GetLocations(cachePrefix + key)
	logger.InfoLogger.Println(cacheData)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	address, err := util.UnMarshalBinaryArray([]byte(cacheData))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	json.NewEncoder(w).Encode(address)
}

func getAllLocations(w http.ResponseWriter, r *http.Request) {
	var locations []entity.Location
	var err error

	iter, err := cnf.GetAllKeys()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	for iter.Next() {
		var location entity.Location
		location.Key = iter.Val()
		cacheData, err := cnf.GetLocations(location.Key)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
			return
		}

		location.Address, err = util.UnMarshalBinaryArray([]byte(cacheData))
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
			return
		}
		locations = append(locations, location)
	}
	if err := iter.Err(); err != nil {
		logger.ErrorLogger.Println(err.Error())
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	json.NewEncoder(w).Encode(entity.Success(locations))
}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func verifyAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token")
		header = strings.TrimSpace(header)

		if header != "hxxjfehp79q69nzp" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(entity.Failure(int(consts.Unauthorized), consts.Unauthorized.Error()))
			return
		}
		next.ServeHTTP(w, r)
	})
}
