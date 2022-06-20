package handler_dir

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
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
	//defer db.Close()
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
