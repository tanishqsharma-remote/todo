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
	"todo/model_dir"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, CookieErr := r.Cookie("session_token")
		if CookieErr != nil {
			log.Fatal(CookieErr)
		}
		sessionToken := c.Value

		userSession, exists := model_dir.Sessions[sessionToken]
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userSession.Expiry.Before(time.Now()) {
			delete(model_dir.Sessions, sessionToken)
			w.WriteHeader(http.StatusUnauthorized)
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
	}
}

func RefreshMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, CookieErr := r.Cookie("session_token")
		if CookieErr != nil {
			log.Fatal(CookieErr)
		}
		sessionToken := c.Value

		userSession, exists := model_dir.Sessions[sessionToken]
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userSession.Expiry.Before(time.Now()) {
			delete(model_dir.Sessions, sessionToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newSessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		model_dir.Sessions[newSessionToken] = model_dir.Session{
			Username: userSession.Username,
			Expiry:   expiresAt,
		}

		delete(model_dir.Sessions, sessionToken)

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   newSessionToken,
			Expires: time.Now().Add(120 * time.Second),
		})
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
	}
}
