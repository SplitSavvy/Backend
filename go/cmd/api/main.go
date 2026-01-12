package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"splitsavvy/internal/database"
)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error while loading env vars")
	}

	DB_URL := os.Getenv("DB_URL")
	pool, err := database.Connect(DB_URL)
	if err != nil{
		log.Fatal("Error while connecting to DB", err)
	}
	defer pool.Close()
	fmt.Println("Database connection successful!")

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HEMLO"))
	})
	fmt.Println("server starting on port :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil{
		log.Fatal("Server failed to start: ", err)
	}
}