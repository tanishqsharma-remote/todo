package handler_dir

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
	"todo/database_dir"
	"todo/model_dir"
)

func Home(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value("Id").(jwt.MapClaims)
	_, err := io.WriteString(w, "You are authorized!")
	if err != nil {
		log.Fatal(err)
	}
	_, er := fmt.Fprintf(w, "Hello %s\n", id)
	if er != nil {
		log.Fatal(er)
	}
}
func Refresh(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value("Id").(jwt.MapClaims)
	_, er := fmt.Fprintf(w, "Session Refreshed %s\n", id)
	if er != nil {
		log.Fatal(er)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var item model_dir.User
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Fatal(err)
	}
	query := "Insert into users(username,password) values($1,$2)"

	_, er := db.Exec(query, item.Username, item.Password)
	if er != nil {
		log.Fatal(er)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var credentials model_dir.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Fatal(err)
	}
	rows, er := db.Query("Select * from users where username=$1", credentials.Username)
	if er != nil {
		log.Fatal(er)
	}
	var authorized model_dir.User
	for rows.Next() {
		ScanErr := rows.Scan(&authorized.Id, &authorized.Username, &authorized.Password)
		if ScanErr != nil {
			log.Fatal(ScanErr)
		}
	}
	if credentials.Password != authorized.Password {
		var fail = []byte("Failed Authentication")
		_, WriteErr := w.Write(fail)
		if WriteErr != nil {
			log.Fatal(WriteErr)
		}
		return
	}
	ExpiryTime := time.Now().Add(time.Minute * 10).Unix()
	Expires := time.Now().Add(time.Minute * 10)
	sessionToken := uuid.NewString()

	query := "insert into sessions(sessiontoken, username, expiry) VALUES ($1,$2,$3)"
	_, exErr := db.Exec(query, sessionToken, authorized.Username, Expires)
	if exErr != nil {
		log.Fatal(exErr)
	}
	w.Header().Add("sessionToken", sessionToken)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["Id"] = authorized.Id
	claims["exp"] = ExpiryTime
	userTokenString, SignErr := token.SignedString(model_dir.JwtKey)
	if SignErr != nil {
		log.Fatalln(err)
	}
	var userToken model_dir.Token
	userToken.Username = authorized.Username
	userToken.TokenString = userTokenString
	EncodeErr := json.NewEncoder(w).Encode(userToken)
	if EncodeErr != nil {
		log.Fatal(EncodeErr)
	}

}
func Logout(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	sessionToken := r.Header.Get("sessionToken")
	query := "delete from sessions where sessiontoken=$1"
	_, execErr := db.Exec(query, sessionToken)
	if execErr != nil {
		log.Fatal(execErr)
	}
	_, er := io.WriteString(w, "Successfully Logged out")
	if er != nil {
		log.Fatal(er)
	}
}
func CreateTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var todoTask model_dir.Todolist
	err := json.NewDecoder(r.Body).Decode(&todoTask)
	if err != nil {
		log.Fatal(err)
	}
	query := "Insert into todolist(user_id, task, completed,archived) values($1,$2,$3,$4)"

	_, er := db.Exec(query, todoTask.UserId, todoTask.Task, todoTask.Completed, todoTask.Archived)
	if er != nil {
		log.Fatal(er)
	}
}
func GetTask(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value("Id").(jwt.MapClaims)
	userid := fmt.Sprint(id["Id"])

	db := database_dir.DBconnect()
	rows, err := db.Query("with pagingCTE as(SELECT user_id,task,completed,archived, row_number() over (order by task) as rowNumber FROM todolist)select user_id,task,completed,archived from pagingCTE where user_id=$1 and rowNumber between ($2-1)*$3+1 and $2*$3", userid, r.URL.Query().Get("pageNum"), r.URL.Query().Get("pageSize"))
	if err != nil {
		log.Fatal(err)
	}

	var items []model_dir.Todolist

	for rows.Next() {
		var item model_dir.Todolist
		err := rows.Scan(&item.UserId, &item.Task, &item.Completed, &item.Archived)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	itemsBytes, _ := json.MarshalIndent(items, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	_, er := w.Write(itemsBytes)
	if er != nil {
		log.Fatal(er)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var Task model_dir.TodoTask
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		log.Fatal(err)
	}

	query := "update todolist set completed=true where task=$1"
	_, er := db.Exec(query, Task.Task)
	if er != nil {
		log.Fatal(er)
	}

}

func ArchiveTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var Task model_dir.TodoTask
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		log.Fatal(err)
	}
	query := "update todolist set archived=true where task=$1"
	_, er := db.Exec(query, Task.Task)
	if er != nil {
		log.Fatal(er)
	}

}
