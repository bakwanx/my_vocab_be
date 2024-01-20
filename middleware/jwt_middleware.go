package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("secret")

func CreateToken(id int, fullname string) (generateToken string, generateRefreshToken string, err error) {
	var EXPIRATION_DURATION = time.Now().Add(time.Hour * 24).Unix()
	var REFRESH_EXPIRATION_DURATION = time.Now().Add(time.Hour * 720).Unix()

	// Create token
	token := jwt.New(JWT_SIGNING_METHOD)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["fullname"] = fullname
	claims["exp"] = EXPIRATION_DURATION

	generateToken, _ = token.SignedString(JWT_SIGNATURE_KEY)

	refreshToken := jwt.New(JWT_SIGNING_METHOD)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["fullname"] = fullname
	rtClaims["exp"] = REFRESH_EXPIRATION_DURATION

	generateRefreshToken, _ = refreshToken.SignedString(JWT_SIGNATURE_KEY)

	return
}

func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return http.StatusUnauthorized, fmt.Errorf("Signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return http.StatusUnauthorized, fmt.Errorf("Signing method invalid")
			}

			return JWT_SIGNATURE_KEY, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
