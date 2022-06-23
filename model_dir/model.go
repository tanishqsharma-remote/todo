package model_dir

import "time"

var JwtKey = []byte("MyKey")

type User struct {
	Id       int    `json:"id"`
	Username string `json:"userName"`
	Password string `json:"passWord"`
}

type TodoTask struct {
	Task string `json:"task"`
}

type Todolist struct {
	UserId    int    `json:"userId"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
	Archived  bool   `json:"archived"`
}
type Credentials struct {
	Username string `json:"userName"`
	Password string `json:"passWord"`
}
type Token struct {
	Username    string `json:"userName"`
	TokenString string `json:"token"`
}

var Sessions = map[string]Session{}

type Session struct {
	Username string    `json:"userName"`
	Expiry   time.Time `json:"expiry"`
}
