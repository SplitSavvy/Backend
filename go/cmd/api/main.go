package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"splitsavvy/internal/database"
	"splitsavvy/internal/groups"
	"splitsavvy/internal/users"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system env vars")
	}

	// Grab your specific DB_URL variable
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in the environment")
	}

	// Initialize the database connection
	pool, err := database.Connect(dbURL)
	if err != nil {
		log.Fatal("Error while connecting to DB: ", err)
	}
	// Defers closing the pool until main() exits
	defer pool.Close()
	fmt.Println("Database connection successful!")

	// Initialize Chi router
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HEMLO FROM SPLITSAVVY"))
	})

	// === AUTH TEMPORARILY PARKED ===
	// We will mount the users/expenses handlers here next!

	usersHandler := users.NewHandler(pool)
	r.Post("/users", usersHandler.HandleCreateUser)

	groupsHandler := groups.NewHandler(pool)
	r.Post("/groups", groupsHandler.HandleGroupRequest)

	// Start server
	addr := ":8080"
	fmt.Println("Server starting on port " + addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
