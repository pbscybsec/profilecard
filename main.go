package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

type User struct {
	FirstName   string `bson:"firstName" json:"firstName"`
	LastName    string `bson:"lastName" json:"lastName"`
	PicturePath string `bson:"picturePath" json:"picturePath"`
	Impressions int    `bson:"impressions" json:"impressions"`
	Status      string `bson:"status" json:"status"`
}

// SetCORS sets the necessary headers for CORS
func SetCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("pbscybsec").Collection("users")

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		SetCORS(&w)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		var user User
		if err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error fetching user", http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		SetCORS(&w)
		if r.Method == "OPTIONS" {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status Ok!"))
	})

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
