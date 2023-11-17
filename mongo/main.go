package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/evanyip05/Cloud/config"
)

var mongoComposed = true
var mongoClient *mongo.Client = nil // close me maybe
var dbName = "myData"
var colName = "col1"

func main() {
    InitMongo() // blocks until mongo client connects
    InitHTTP()
}

// http handler
func InitHTTP() {    
    // make a post request with encoded json
    // data = json.marshal(whatever)
    // http.Post("http://localhost:8080/put", "application/json", bytes.NewBuffer(data))
    http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
	    if r.Method != http.MethodPost {w.WriteHeader(http.StatusMethodNotAllowed); return}

	    var payload config.Payload
        if json.NewDecoder(r.Body).Decode(&payload) != nil {w.WriteHeader(http.StatusBadRequest); return}

        // convert decoded json struct to bson, write to mongo
        objID := WriteData(dbName, colName, StructToBson(payload))

        MarshalAndSend(w, config.MongoID{ID: strings.Split(objID.String(), "\"")[1]}) // send back the id of the mongo entry
    })

    // localhost:8080/get?filter=key:value key2:value2     parameter filtering
    // OR
    // localhost:8080/get?_id=6537236861552323c9b4c264     mongo id
    http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        requests := r.URL.Query()

        // parse url _id param, prioritize this if it exists
        if requests.Get("_id") != "" {
            id, err := primitive.ObjectIDFromHex(requests.Get("_id"))
	        if err != nil {w.WriteHeader(http.StatusBadRequest); MarshalAndSend(w, "malformed \"_id\""); return}     
            MarshalAndSend(w, config.MongoEntries{Entries: ReadData(dbName, colName, bson.M{"_id": id})}) // send the single entry
            return
        }

        // check for filter based search 
        if requests.Get("filter") == "" {w.WriteHeader(http.StatusBadRequest); MarshalAndSend(w, "missing filter or _id"); return}

        // parse url filter param
        filter := bson.M{}
        for _, f := range strings.Split(requests.Get("filter"), " ") {
            kvp := strings.Split(f, ":")
            filter[kvp[0]] = kvp[1]
        }

        MarshalAndSend(w, config.MongoEntries{Entries: ReadData(dbName, colName, filter)})
    })

    // server conf message
    log.Println("listening on 0.0.0.0:8080")
    err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Println("error listening and serving")
	}
}


// mongo handler 
func InitMongo() {
    var clientOptions *options.ClientOptions
    if mongoComposed {
        clientOptions = options.Client().ApplyURI("mongodb://database:27017") // equivalent of localhost while its on a docker container (for the go api)
    } else {
        clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
    }

    // connect go api to mongo
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the MongoDB server to verify the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal("Failed to ping the MongoDB server:", err)
    }

    mongoClient = client // global

    log.Println("Connected to MongoDB!")
}

// mongo write boilerplate
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

// mongo read boilerplate
func ReadData(dbName, colName string, selector primitive.M) []primitive.M {
    collection := mongoClient.Database(dbName).Collection(colName)

    log.Println(selector)

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

// send basic messages 
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