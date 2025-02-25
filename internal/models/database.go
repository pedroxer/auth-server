package models

type User struct {
	ID       int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Team     string `json:"team"`
	Role     int    `json:"role"`
	Position int    `json:"position"`
}

type App struct {
	ID     int    `json:"app_id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}
