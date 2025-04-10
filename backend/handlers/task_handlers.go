package handlers

import (
	"container/heap"
	"encoding/json"
	"log"
	"net/http"

	"task_manager/backend/db"
	"task_manager/backend/models"

	"github.com/gorilla/mux"

	"strconv"
)

type TaskHeap []models.Task

func (th TaskHeap) Len() int { return len(th) }

func (th TaskHeap) Less(i, j int) bool {
	return th[i].Priority > th[j].Priority
}

func (th TaskHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

func (th *TaskHeap) Push(x interface{}) {
	*th = append(*th, x.(models.Task))
}

func (th *TaskHeap) Pop() interface{} {
	old := *th
	n := len(old)
	item := old[n-1]
	*th = old[0 : n-1]
	return item
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, user_id, title, priority FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Priority)
		if err != nil {
			log.Println(err)
			continue
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks(user_id, title, priority) VALUES($1, $2, $3) RETURNING id`
	err = db.DB.QueryRow(query, task.UserID, task.Title, task.Priority).Scan(&task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/tasks/"):]

	rows, err := db.DB.Query("SELECT id, user_id, title, priority FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Priority)
		if err != nil {
			log.Println(err)
			continue
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func GetTopPriorityTasksFromMemory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, title, priority FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	taskHeap := &TaskHeap{}
	heap.Init(taskHeap)

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Priority); err != nil {
			log.Println(err)
			continue
		}

		heap.Push(taskHeap, task)

		if taskHeap.Len() > 5 {
			heap.Pop(taskHeap)
		}
	}

	var topTasks []models.Task
	for taskHeap.Len() > 0 {
		topTasks = append([]models.Task{heap.Pop(taskHeap).(models.Task)}, topTasks...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topTasks)
}
