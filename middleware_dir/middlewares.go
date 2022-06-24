package middleware_dir

import (
	"context"
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

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := database_dir.DBconnect()

		sessionToken := r.Header.Get("sessionToken")
		rows, err := db.Query("select username,expiry from sessions where sessiontoken=$1", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Query error")
			return
		}
		var validSession model_dir.Session
		for rows.Next() {
			ScanErr := rows.Scan(&validSession.Username, &validSession.Expiry)
			if ScanErr != nil {
				log.Fatal(ScanErr)
			}
		}
		if validSession.Expiry.Before(time.Now()) {
			query := "delete from sessions where sessiontoken=$1"
			_, execErr := db.Exec(query, sessionToken)
			if execErr != nil {
				log.Fatal(execErr)
			}
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "exp error")
			return
		}
		var userToken model_dir.Token
		DecodeErr := json.NewDecoder(r.Body).Decode(&userToken)
		if DecodeErr != nil {
			log.Fatal(DecodeErr)
		}

		checkToken, err := jwt.Parse(userToken.TokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token. ")
			}
			return model_dir.JwtKey, nil
		})
		if err != nil {
			log.Fatalln(err)
		}
		claims, ok := checkToken.Claims.(jwt.MapClaims)

		if ok && checkToken.Valid {
			ctx := context.WithValue(r.Context(), "Id", claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func RefreshMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := database_dir.DBconnect()

		sessionToken := r.Header.Get("sessionToken")
		rows, err := db.Query("select username,expiry from sessions where sessiontoken=$1", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var validSession model_dir.Session
		for rows.Next() {
			ScanErr := rows.Scan(&validSession.Username, &validSession.Expiry)
			if ScanErr != nil {
				log.Fatal(ScanErr)
			}
		}
		if validSession.Expiry.Before(time.Now()) {
			query := "delete from sessions where sessiontoken=$1"
			_, execErr := db.Exec(query, sessionToken)
			if execErr != nil {
				log.Fatal(execErr)
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newSessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)
		_, qErr := db.Exec("insert into sessions(sessiontoken, username, expiry) VALUES ($1,$2,$3)", newSessionToken, validSession.Username, expiresAt)
		if qErr != nil {
			log.Fatal(qErr)
		}

		query := "delete from sessions where sessiontoken=$1"
		_, execErr := db.Exec(query, sessionToken)
		if execErr != nil {
			log.Fatal(execErr)
		}
		w.Header().Add("sessionToken", newSessionToken)

		var userToken model_dir.Token
		DecodeErr := json.NewDecoder(r.Body).Decode(&userToken)
		if DecodeErr != nil {
			log.Fatal(DecodeErr)
		}

		checkToken, err := jwt.Parse(userToken.TokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token. ")
			}
			return model_dir.JwtKey, nil
		})
		if err != nil {
			log.Fatalln(err)
		}
		claims, ok := checkToken.Claims.(jwt.MapClaims)

		if ok && checkToken.Valid {
			ctx := context.WithValue(r.Context(), "Id", claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
