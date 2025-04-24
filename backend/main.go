package main

import (
	"fmt"
	"log"
	"net/http"

	"task_manager/backend/config"
	"task_manager/backend/db"
	"task_manager/backend/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	config.LoadConfig("config.json")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.SSLMode,
	)

	db.InitDB(connStr)

	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{user_id}", handlers.GetTasksByUserID).Methods("GET")
	r.HandleFunc("/subtasks", handlers.CreateSubtask).Methods("POST")
	r.HandleFunc("/tasks/{id}/subtasks", handlers.GetSubtasksByTaskID).Methods("GET")
	r.HandleFunc("/users/details", handlers.GetUsersWithTasksAndSubtasks).Methods("GET")
	r.HandleFunc("/users/details/{id}", handlers.GetUserDetailsByID).Methods("GET")
	r.HandleFunc("/users/{id}/top-tasks", handlers.GetTopPriorityTasksFromMemory).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	log.Println("Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
