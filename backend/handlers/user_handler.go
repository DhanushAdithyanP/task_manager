package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"task_manager/backend/db"
	"task_manager/backend/models"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow("INSERT INTO users(username) VALUES($1) RETURNING id", user.Username).Scan(&user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUsersWithTasksAndSubtasks(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}

	userRows, err := db.DB.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer userRows.Close()

	for userRows.Next() {
		var user models.User
		if err := userRows.Scan(&user.ID, &user.Username); err != nil {
			log.Println(err)
			continue
		}

		taskRows, err := db.DB.Query("SELECT id, title, priority FROM tasks WHERE user_id = $1", user.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		var tasks []models.Task
		for taskRows.Next() {
			var task models.Task
			if err := taskRows.Scan(&task.ID, &task.Title, &task.Priority); err != nil {
				log.Println(err)
				continue
			}
			task.UserID = user.ID

			subtaskRows, err := db.DB.Query("SELECT id, title FROM subtasks WHERE task_id = $1", task.ID)
			if err != nil {
				log.Println(err)
				continue
			}

			var subtasks []models.Subtask
			for subtaskRows.Next() {
				var subtask models.Subtask
				if err := subtaskRows.Scan(&subtask.ID, &subtask.Title); err != nil {
					log.Println(err)
					continue
				}
				subtask.TaskID = task.ID
				subtasks = append(subtasks, subtask)
			}
			subtaskRows.Close()

			task.Subtasks = subtasks
			tasks = append(tasks, task)
		}
		taskRows.Close()

		user.Tasks = tasks
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserDetailsByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	err := db.DB.QueryRow("SELECT id, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	taskRows, err := db.DB.Query("SELECT id, title, priority FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	defer taskRows.Close()

	var tasks []models.Task
	for taskRows.Next() {
		var task models.Task
		err := taskRows.Scan(&task.ID, &task.Title, &task.Priority)
		if err != nil {
			log.Println(err)
			continue
		}
		task.UserID = user.ID

		subtaskRows, err := db.DB.Query("SELECT id, title FROM subtasks WHERE task_id = $1", task.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		defer subtaskRows.Close()

		var subtasks []models.Subtask
		for subtaskRows.Next() {
			var subtask models.Subtask
			err := subtaskRows.Scan(&subtask.ID, &subtask.Title)
			if err != nil {
				log.Println(err)
				continue
			}
			subtask.TaskID = task.ID
			subtasks = append(subtasks, subtask)
		}

		task.Subtasks = subtasks
		tasks = append(tasks, task)
	}

	user.Tasks = tasks

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
