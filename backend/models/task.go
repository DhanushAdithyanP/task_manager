package models

type Task struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Title    string    `json:"title"`
	Priority int       `json:"priority"`
	Subtasks []Subtask `json:"subtasks,omitempty"`
}
