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
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	_, er := fmt.Fprintf(w, "Hello %s\n", id)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
}
func Refresh(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value("Id").(jwt.MapClaims)
	_, er := fmt.Fprintf(w, "Session Refreshed %s\n", id)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var item model_dir.User
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	_, er := database_dir.InsertUser(db, item)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var credentials model_dir.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	rows, er := database_dir.GetUser(db, credentials)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
	var authorized model_dir.User
	for rows.Next() {
		ScanErr := rows.Scan(&authorized.Id, &authorized.Username, &authorized.Password)
		if ScanErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(ScanErr)
			return
		}
	}
	if credentials.Password != authorized.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ExpiryTime := time.Now().Add(time.Minute * 10).Unix()
	Expires := time.Now().Add(time.Minute * 10)
	sessionToken := uuid.NewString()

	_, exErr := database_dir.InsertSession(db, sessionToken, authorized, Expires)
	if exErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(exErr)
		return
	}
	w.Header().Add("sessionToken", sessionToken)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["Id"] = authorized.Id
	claims["exp"] = ExpiryTime
	userTokenString, SignErr := token.SignedString(model_dir.JwtKey)
	if SignErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(SignErr)
		return
	}
	var userToken model_dir.Token
	userToken.Username = authorized.Username
	userToken.TokenString = userTokenString
	EncodeErr := json.NewEncoder(w).Encode(userToken)
	if EncodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(EncodeErr)
		return
	}

}
func Logout(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	sessionToken := r.Header.Get("sessionToken")
	_, execErr := database_dir.DelSession(db, sessionToken)
	if execErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(execErr)
		return
	}
	_, er := io.WriteString(w, "Successfully Logged out")
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
}
func CreateTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var todoTask model_dir.Todolist
	err := json.NewDecoder(r.Body).Decode(&todoTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	_, er := database_dir.InsertTask(db, todoTask)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}
}
func GetTask(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Context().Value("Id").(jwt.MapClaims)
	userid := fmt.Sprint(id["Id"])

	db := database_dir.DBconnect()
	rows, err := database_dir.GetTaskRows(db, userid, r.URL.Query().Get("pageNum"), r.URL.Query().Get("pageSize"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var items []model_dir.Todolist

	for rows.Next() {
		var item model_dir.Todolist
		err := rows.Scan(&item.UserId, &item.Task, &item.Completed, &item.Archived)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		items = append(items, item)
	}

	itemsBytes, _ := json.MarshalIndent(items, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	_, er := w.Write(itemsBytes)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}(rows)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var Task model_dir.TodoTask
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	_, er := database_dir.DoneTaskQuery(db, Task)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}

}

func ArchiveTask(w http.ResponseWriter, r *http.Request) {
	db := database_dir.DBconnect()
	var Task model_dir.TodoTask
	err := json.NewDecoder(r.Body).Decode(&Task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	_, er := database_dir.ArchiveTaskQuery(db, Task)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(er)
		return
	}

}
