package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func enableCros(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Method", "POST, GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func saveIPHandler(w http.ResponseWriter, r *http.Request) {
	enableCros(&w)
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed!", http.StatusMethodNotAllowed)
		return
	}
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing the JSON request!", http.StatusBadRequest)
		return
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	data["time"] = currentTime

	collection := mongoClient.Database("userData").Collection("VisitorLogs")
	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		http.Error(w, "Error saving data to MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "IP saved successfully")
}
func main() {
	var err error
	mongoURI := os.Getenv("MONGO_URI")
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Erros Connecting to MongoDB\n")
		log.Fatal(err)
	}
	fmt.Printf("Connected to MongoDB\n")
	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)
	http.HandleFunc("/api/saveIP", saveIPHandler)
	fmt.Printf("Starting server on port 80...\n")
	log.Fatal(http.ListenAndServe(":80", nil))
}
