package main

import (
	"fmt"
	"log"
	"net/http"

	"task_manager/db"
	"task_manager/handlers"

	"github.com/gorilla/mux"
)

func main() {
	connStr := "user=postgres password=KA10m1355@car22 dbname=task_manager sslmode=disable"
	db.InitDB(connStr)

	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	fmt.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
