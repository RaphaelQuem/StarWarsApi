
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "github.com/RaphaelQuem/StarWarsApi/config"
	. "github.com/RaphaelQuem/StarWarsApi/dao"
	. "github.com/RaphaelQuem/StarWarsApi/models"
)

var config = Config{}
var dao = PlanetsDAO{}

// GET list of planets
func AllPlanetsEndPoint(w http.ResponseWriter, r *http.Request) {
	planets, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, planets)
}

// GET a planet by its ID
func FindPlanetEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	planet, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Planet ID")
		return
	}
	respondWithJson(w, http.StatusOK, planet)
}

// GET a planet by its Name
func FindPlanetByNameEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["name"])
	planet, err := dao.FindByName(params["name"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Planet Name")
		return
	}
	respondWithJson(w, http.StatusOK, planet)
}

// POST a new planet
func CreatePlanetEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var planet Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	planet.ID = bson.NewObjectId()
	if err := dao.Insert(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, planet)
}

// PUT update an existing planet
func UpdatePlanetEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var planet Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing planet
func DeletePlanetEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var planet Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/planets", AllPlanetsEndPoint).Methods("GET")
	r.HandleFunc("/planets/name/{name}", FindPlanetByNameEndPoint).Methods("GET")
	r.HandleFunc("/planets", CreatePlanetEndPoint).Methods("POST")
	r.HandleFunc("/planets", UpdatePlanetEndPoint).Methods("PUT")
	r.HandleFunc("/planets", DeletePlanetEndPoint).Methods("DELETE")
	r.HandleFunc("/planets/{id}", FindPlanetEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

