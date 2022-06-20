package model_dir

var JwtKey = []byte("MyKey")

type User struct {
	Id       int    `json:"Id"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type Todolist struct {
	TaskId    int    `json:"Task_Id"`
	UserId    int    `json:"User_Id"`
	Task      string `json:"Task"`
	Completed bool   `json:"Completed"`
}
type Credentials struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type Token struct {
	Username    string `json:"Username"`
	TokenString string `json:"Token"`
}
