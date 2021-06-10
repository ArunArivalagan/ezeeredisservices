package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	router.HandleFunc("/ezeeredisservices/locations/add", addAddress).Methods("POST")
	router.HandleFunc("/ezeeredisservices/locations/{key}", getLocations).Methods("GET")
	router.HandleFunc("/ezeeredisservices/locations", getAllLocations).Methods("GET")

	http.ListenAndServe(":8080", router)
}

func addAddress(w http.ResponseWriter, r *http.Request) {
	var locations []entity.Location
	json.NewDecoder(r.Body).Decode(&locations)
	logger.InfoLogger.Println(locations)

	for _, location := range locations {
		if !location.IsValid() {
			json.NewEncoder(w).Encode(entity.Failure(int(consts.Unauthorized), consts.Unauthorized.Error()))
			return
		}
	}

	for _, location := range locations {
		var err error
		data, err := util.MarshalBinaryArray(location.Address)
		if err != nil {
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UpdateFailed), consts.UpdateFailed.Error()))
			return
		}
		err = cnf.PutLocation(cachePrefix+location.Key, data)
		if err != nil {
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UpdateFailed), consts.UpdateFailed.Error()))
			return
		}
	}

	json.NewEncoder(w).Encode(entity.Success(nil))
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
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	address, err := util.UnMarshalBinaryArray([]byte(cacheData))
	if err != nil {
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
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	for iter.Next() {
		var location entity.Location
		location.Key = iter.Val()
		cacheData, err := cnf.GetLocations(location.Key)
		if err != nil {
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
			return
		}

		location.Address, err = util.UnMarshalBinaryArray([]byte(cacheData))
		if err != nil {
			json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
			return
		}
		locations = append(locations, location)
	}
	if err := iter.Err(); err != nil {
		json.NewEncoder(w).Encode(entity.Failure(int(consts.UnableToProvideData), consts.UnableToProvideData.Error()))
		return
	}
	json.NewEncoder(w).Encode(locations)
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
