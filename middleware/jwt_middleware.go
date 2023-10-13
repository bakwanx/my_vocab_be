package middleware

import (
	"encoding/json"
	"my_vocab/dto/out"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id int, fullname string) (generateToken string, refreshToken string, err error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["fullname"] = fullname
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	generateToken, err = token.SignedString([]byte("secret"))

	rtClaims := jwt.MapClaims{}
	claims["refresh_id"] = id
	claims["fullname"] = fullname
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, _ = rt.SignedString([]byte("secret"))

	return
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		var result out.Response

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				result.Code = http.StatusUnauthorized
				result.Status = "Failed"
				result.Message = "Unauthorized"
				json.NewEncoder(w).Encode(result)
			}
		}

		tokenString := c.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			result.Code = http.StatusUnauthorized
			result.Status = "Failed"
			result.Message = "Unauthorized"
			json.NewEncoder(w).Encode(result)
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			result.Code = http.StatusUnauthorized
			result.Status = "Failed"
			result.Message = "Unauthorized"
			json.NewEncoder(w).Encode(result)
		}
	})
}
