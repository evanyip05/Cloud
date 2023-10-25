package config

import "go.mongodb.org/mongo-driver/bson"

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MongoID struct {
	ID string `json:"id"`
}

type MongoEntries struct {
	Entries []bson.M `json:"entries"`
}