package middleware_dir

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"todo/model_dir"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
