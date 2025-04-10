package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"task_manager/backend/db"
	"task_manager/backend/models"

	"github.com/gorilla/mux"
)

func CreateSubtask(w http.ResponseWriter, r *http.Request) {
	var subtask models.Subtask
	err := json.NewDecoder(r.Body).Decode(&subtask)
	if err != nil {
		http.Error(w, "Invalid subtask data", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO subtasks (task_id, title) VALUES ($1, $2) RETURNING id`
	err = db.DB.QueryRow(query, subtask.TaskID, subtask.Title).Scan(&subtask.ID)
	if err != nil {
		http.Error(w, "Failed to insert subtask", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subtask)
}

func GetSubtasksByTaskID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	rows, err := db.DB.Query("SELECT id, task_id, title FROM subtasks WHERE task_id = $1", taskID)
	if err != nil {
		http.Error(w, "Failed to fetch subtasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subtasks []models.Subtask
	for rows.Next() {
		var st models.Subtask
		err := rows.Scan(&st.ID, &st.TaskID, &st.Title)
		if err != nil {
			log.Println(err)
			continue
		}
		subtasks = append(subtasks, st)
	}

	json.NewEncoder(w).Encode(subtasks)
}
