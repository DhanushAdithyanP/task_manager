package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash,omitempty"` // Don't expose this in responses
	Role         string `json:"role"`
	Tasks        []Task `json:"tasks,omitempty"`
}
