package model_dir

var JwtKey = []byte("MyKey")

type User struct {
	Id       int    `json:"Id"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type TodoTask struct {
	Task string `json:"Task"`
}

type Todolist struct {
	UserId    int    `json:"User_Id"`
	Task      string `json:"Task"`
	Completed bool   `json:"Completed"`
	Archived  bool   `json:"Archived"`
}
type Credentials struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type Token struct {
	Username    string `json:"Username"`
	TokenString string `json:"Token"`
}
