package dao

import (
	"log"

	//"strings"
	. "github.com/RaphaelQuem/StarWarsApi/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/peterhellberg/swapi"

)

type PlanetsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "planets"
)

// Conecta com o MONGO
func (m *PlanetsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Lista de planetas
func (m *PlanetsDAO) FindAll() ([]Planet, error) {
	var planets []Planet
	err := db.C(COLLECTION).Find(bson.M{}).All(&planets)
	for i := 0; i < len(planets); i++ {
		planets[i].Appearances = GetPlanetApperances(planets[i].Name)
	}
	return planets, err
}

// Planeta por ID
func (m *PlanetsDAO) FindById(id string) (Planet, error) {
	var planet Planet
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&planet)
	planet.Appearances = GetPlanetApperances(planet.Name)
	return planet, err
}

//Planeta por nome
func (m *PlanetsDAO) FindByName(name string) (Planet, error) {
	var planet Planet
	err := db.C(COLLECTION).Find(name).Select(bson.M{"name": name}).One(&planet)
	return planet, err
}

// Insere planeta
func (m *PlanetsDAO) Insert(planet Planet) error {
	err := db.C(COLLECTION).Insert(&planet)
	return err
}

// Deleta planeta
func (m *PlanetsDAO) Delete(planet Planet) error {
	planet.Appearances = 0 
	err := db.C(COLLECTION).Remove(&planet)
	return err
}

// Update Planeta
func (m *PlanetsDAO) Update(planet Planet) error {
	err := db.C(COLLECTION).UpdateId(planet.ID, &planet)
	return err
}

// Busca o número de aparições na api publica
func  GetPlanetApperances(name string) int{
	c := swapi.DefaultClient
	var cname = "name"
	var page = 1
	for cname != ""{
		if atst, err := c.Planets(page); err == nil {
			
			for _, element := range atst.Results {
				cname = element.Name
				if(cname == name){
					return len(element.FilmURLs)
				}
			}
			
		}
		page++
	}	
	return 0
}
