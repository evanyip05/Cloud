package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// close me somehow
var mongoClient *mongo.Client = nil
var dbName = "myData"
var colName = "col1"

func main() {
    InitMongo() // blocks until mongo client connects
    InitHTTP()
}

func InitHTTP() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {MarshalAndSend(w, "this is the backend.")})
    
    http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
        // Check the HTTP method
	    if r.Method != http.MethodPost {w.WriteHeader(http.StatusMethodNotAllowed); return}

	    // Parse the JSON payload
	    var payload json.RawMessage
        if json.NewDecoder(r.Body).Decode(&payload) != nil {w.WriteHeader(http.StatusBadRequest); return}

        log.Println(payload)

        // convert decoded json struct to bson, write to mongo
        // find a better storage method to write anything to db, or find final datastruct
        objID := WriteData(dbName, colName, StructToBson(payload))

        // return id
        MarshalAndSend(w, strings.Split(objID.String(), "\"")[1])
    })

    // localhost:8080/get?filter=key:value key2:value2
    http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        requests := r.URL.Query()
        filters := strings.Split(requests.Get("filter"), " ")
        filter := bson.M{}

        for _, f := range filters {
            kvp := strings.Split(f, ":")
            filter[kvp[0]] = kvp[1]
        }

        MarshalAndSend(w, {Entries: ReadData(dbName, colName, filter)})
    })

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println("got request from", fmt.Sprint(r.Host, r.URL))
        MarshalAndSend(w, "pong")
	})

    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        log.Println("sending OK")
        w.WriteHeader(http.StatusOK)
    })

    log.Println("listening on 0.0.0.0:8080")

    err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Println("error listening and serving")
	}
}

func InitMongo() {
    // while composed
	//clientOptions := options.Client().ApplyURI("mongodb://database:27017")
	// while standalone
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the MongoDB server to verify the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal("Failed to ping the MongoDB server:", err)
    }

    mongoClient = client

    log.Println("Connected to MongoDB!")
}

func WriteData(dbName, colName string, data primitive.M) primitive.ObjectID {
    collection := mongoClient.Database(dbName).Collection(colName)

    insertResult, err := collection.InsertOne(nil, data)
    
	if err != nil {
		log.Println("Failed to insert document:", err)
		return primitive.NilObjectID
	}

    insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)

    if !ok {
        log.Println("Failed to get id")
		return primitive.NilObjectID
    }

    return insertedID
}

func ReadData(dbName, colName string, selector primitive.M) []primitive.M {
    collection := mongoClient.Database(dbName).Collection(colName)

    cursor, err := collection.Find(context.TODO(), selector)

    if err != nil {
        log.Println("Failed to execute find query:", err)
        return []primitive.M{}
    }

    defer cursor.Close(context.TODO())

    var results []primitive.M

    for cursor.Next(context.TODO()) {
        var document bson.M
        if err := cursor.Decode(&document); err != nil {
            log.Println("Failed to decode document:", err)
            return []primitive.M{}
        }
        results = append(results, document)
    }

    return results
}

func MarshalAndSend(w http.ResponseWriter, data any) {
    json, err := json.Marshal(data)
    if err != nil {w.WriteHeader(http.StatusInternalServerError); return}
    w.Write(json)
}

func StructToBson(structData any) bson.M {
    payloadData := reflect.ValueOf(structData)
    bsonned := bson.M{}

    for i := 0; i < payloadData.NumField(); i++ {
        fieldName := payloadData.Type().Field(i).Name
        fieldValue := payloadData.Field(i).Interface()
        bsonned[fieldName] = fieldValue
    }

    return bsonned
}

/*
volumes:
    - ../config:/app/config
*/