package models

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type App struct {
	ID     int    `json:"app_id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}
