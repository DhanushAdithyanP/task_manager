package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"task_manager/backend/db"
	"task_manager/backend/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user struct {
		ID       int
		Username string
		Password string
		Role     string
	}

	err := db.DB.QueryRow(`SELECT username, password_hash, role FROM users WHERE username=$1`, req.Username).
		Scan(&user.Username, &user.Password, &user.Role)

	if err != nil {
		fmt.Println("Query error:", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if req.Password != user.Password {
		fmt.Println("Password mismatch")
		fmt.Println("Entered:", req.Password)
		fmt.Println("Expected:", user.Password)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"role":  user.Role,
	})
}
