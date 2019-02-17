package models

import "gopkg.in/mgo.v2/bson"

type Planet struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Weather  	string        `bson:"weather" json:"weather"`
	Terrain		string        `bson:"terrain" json:"terrain"`
	Appearances int			  `bson:"appearances" json:"appearances"`
}
