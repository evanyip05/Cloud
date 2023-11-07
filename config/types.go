package config

import "go.mongodb.org/mongo-driver/bson"

// through the url rn
type MongoGetRequest struct {}

// type to send to mongo container (should be the only type in the db, or seperate different types by mongo dbs or collections)
// stringify json data
type MongoPutRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MongoPutResponse struct {
	ID string `json:"id"`
}

type MongoGetResponse struct {
	Entries []bson.M `json:"entries"`
}

/*

mongo server response
mongo server request

*/