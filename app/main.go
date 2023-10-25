package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./Cloud/config"
)

// uuid generated by mongo, gets returned on any write
//type Payload struct {
//	Name string `json:"name"`
//	Data string `json:"data"`
//}
//
//type MongoID struct {
//	ID string `json:"id"`
//}
//
//type MongoEntries struct {
//	Entries []bson.M `json:"entries"`
//}

func main() {
	write()
	read()
}

func read() {
	r, err := http.Get("http://localhost:8080/get?filter=name:someData")

	if err != nil {
		fmt.Println("Error!:", err)
		return
	}

	var mongoentries config.MongoEntries
	err = json.NewDecoder(r.Body).Decode(&mongoentries)

	if err != nil {
		log.Fatal(err)
	}

	// Access the decoded data
	fmt.Println("resp:", mongoentries)
}

func write() {
	data, err := json.Marshal(config.Payload{
		Name: "someData",
		Data: "lmao",
	})

	if err != nil {
		log.Println("Error!:", err)
	}

	r, err := http.Post("http://localhost:8080/put", "application/json", bytes.NewBuffer(data))

	// Decode the response body
	var mongoid config.MongoID
	err = json.NewDecoder(r.Body).Decode(&mongoid)

	if err != nil {
		log.Fatal(err)
	}

	// Access the decoded data
	fmt.Println("resp:", mongoid)
}