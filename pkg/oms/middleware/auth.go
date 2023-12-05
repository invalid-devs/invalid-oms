package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/invalid-devs/invalid-oms/pkg/oms/response"
	"net/http"
)

var Secret = []byte("L6IuFWc8LWBrvsyyITcHdZl7vIK1ujysM9WVPa4jAGg")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		tokenString := r.Header.Get("Authorization")
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return Secret, nil
		})

		if err != nil {
			msg := response.Error{
				Message: "Unauthorized!",
			}
			response.JSON(w, http.StatusUnauthorized, msg)
			return
		}

		next.ServeHTTP(w, r)

	})
}
