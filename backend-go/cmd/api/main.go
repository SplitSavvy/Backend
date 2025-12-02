package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting SplitSavvy Go Server on port 8080...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the SplitSavvy Go API!"))
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}