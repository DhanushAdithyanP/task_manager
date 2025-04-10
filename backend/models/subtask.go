package models

type Subtask struct {
	ID     int    `json:"id"`
	TaskID int    `json:"task_id"`
	Title  string `json:"title"`
}
