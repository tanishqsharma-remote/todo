package middleware_dir

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
		rows, err := database_dir.GetSession(db, sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var validSession model_dir.Session
		for rows.Next() {
			ScanErr := rows.Scan(&validSession.Username, &validSession.Expiry)
			if ScanErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(ScanErr)
				return
			}
		}
		if validSession.Expiry.Before(time.Now()) {
			_, execErr := database_dir.DelSession(db, sessionToken)
			if execErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(execErr)
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var userToken model_dir.Token
		DecodeErr := json.NewDecoder(r.Body).Decode(&userToken)
		if DecodeErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(DecodeErr)
			return
		}

		checkToken, err := jwt.Parse(userToken.TokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token. ")
			}
			return model_dir.JwtKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
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
		rows, err := database_dir.GetSession(db, sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var validSession model_dir.Session
		for rows.Next() {
			ScanErr := rows.Scan(&validSession.Username, &validSession.Expiry)
			if ScanErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(ScanErr)
				return
			}
		}
		if validSession.Expiry.Before(time.Now()) {
			_, execErr := database_dir.DelSession(db, sessionToken)
			if execErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(execErr)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newSessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)
		_, qErr := database_dir.InsertRefreshedSession(db, newSessionToken, validSession, expiresAt)
		if qErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(qErr)
			return
		}

		_, execErr := database_dir.DelSession(db, sessionToken)
		if execErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(execErr)
			return
		}
		w.Header().Add("sessionToken", newSessionToken)

		var userToken model_dir.Token
		DecodeErr := json.NewDecoder(r.Body).Decode(&userToken)
		if DecodeErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(DecodeErr)
			return
		}

		checkToken, err := jwt.Parse(userToken.TokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token. ")
			}
			return model_dir.JwtKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		claims, ok := checkToken.Claims.(jwt.MapClaims)

		if ok && checkToken.Valid {
			ctx := context.WithValue(r.Context(), "Id", claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
