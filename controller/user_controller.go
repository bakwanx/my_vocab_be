package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"my_vocab/config"
	"my_vocab/dto/out"
	"my_vocab/middleware"
	"my_vocab/models"
	"my_vocab/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(response http.ResponseWriter, request *http.Request) {
	var (
		result out.Response
		user   models.User
	)
	response.Header().Set("Content-Type", "application/json")
	email := request.FormValue("email")
	password := request.FormValue("password")
	fullname := request.FormValue("fullname")
	file, fileHeader, err := request.FormFile("profile")
	timeNow := time.Now()

	// check user
	checkUser := models.User{}
	config.DB.Where("email = ?", email).First(&checkUser)
	if checkUser.Fullname != "" {
		response.WriteHeader(http.StatusConflict)
		result.Code = http.StatusConflict
		result.Status = "OK"
		result.Message = "Pengguna sudah terdaftar"
		json.NewEncoder(response).Encode(result)
		return
	}

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status internal server error"
		json.NewEncoder(response).Encode(result)
		return
	}

	defer file.Close()

	err = os.MkdirAll("./public", os.ModePerm)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status internal server error"
		json.NewEncoder(response).Encode(result)
		return
	}

	// Create a new file in the uploads directory
	f, err := os.Create(fmt.Sprintf("./public/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status internal server error"
		json.NewEncoder(response).Encode(result)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status internal server error"
		json.NewEncoder(response).Encode(result)
		return
	}

	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	user = models.User{
		Email:     email,
		Password:  string(hashedPassword),
		Fullname:  fullname,
		Profile:   utils.ImageUrlProvider(f.Name(), request),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err = config.DB.Save(&user).Error

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result.Code = http.StatusInternalServerError
		result.Status = "Failed"
		result.Message = "Status Internal Server Error"
		json.NewEncoder(response).Encode(result)
		return
	}

	response.WriteHeader(http.StatusOK)
	result.Code = 200
	result.Status = "Success"
	result.Data = user
	result.Message = "Berhasil mendaftar"
	json.NewEncoder(response).Encode(result)
	return
}

func RefreshToken(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var result out.Response
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	err := json.NewDecoder(request.Body).Decode(&tokenReq)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		result.Code = http.StatusBadRequest
		result.Status = "Failed"
		result.Message = "Bad request"
		json.NewEncoder(response).Encode(result)
		return
	}

	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return middleware.JWT_SIGNATURE_KEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		id, err := strconv.Atoi(fmt.Sprint(claims["id"]))
		fullname := fmt.Sprint(claims["fullname"])
		token, refreshToken, err := middleware.CreateToken(id, fullname)

		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			result.Code = http.StatusInternalServerError
			result.Status = "Failed"
			result.Message = "Error"
			json.NewEncoder(response).Encode(result)
			return
		}
		tokenResponse := map[string]interface{}{
			"token":         token,
			"refresh_token": refreshToken,
		}

		response.WriteHeader(http.StatusOK)
		result.Code = http.StatusOK
		result.Status = "Success"
		result.Message = "Success generate token"
		result.Data = tokenResponse
		json.NewEncoder(response).Encode(result)
		return
	}

	response.WriteHeader(http.StatusInternalServerError)
	result.Code = http.StatusInternalServerError
	result.Status = "Failed"
	result.Message = "Error"
	json.NewEncoder(response).Encode(result)
	return
}

func Login(response http.ResponseWriter, request *http.Request) {
	var result out.Response
	response.Header().Set("Content-Type", "application/json")
	user := models.User{}
	dbUser := models.User{}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		result.Code = http.StatusBadRequest
		result.Status = "Failed"
		result.Message = "Bad request"
		json.NewEncoder(response).Encode(result)
		return
	}

	err = config.DB.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		result.Code = http.StatusNotFound
		result.Status = "Failed"
		result.Message = "User not found"
		json.NewEncoder(response).Encode(result)
		return
	}

	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		response.WriteHeader(http.StatusUnauthorized)
		result.Code = http.StatusUnauthorized
		result.Status = "Failed"
		result.Message = "User unauthorized"
		json.NewEncoder(response).Encode(result)
		return
	}

	token, refreshToken, err := middleware.CreateToken(user.IdUser, user.Fullname)

	userResponse := map[string]interface{}{
		"idUser":        dbUser.IdUser,
		"email":         dbUser.Email,
		"fullname":      dbUser.Fullname,
		"profile":       dbUser.Profile,
		"token":         token,
		"refresh_token": refreshToken,
	}

	response.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Status = "Success"
	result.Message = "Success login"
	result.Data = userResponse
	json.NewEncoder(response).Encode(result)
	return
}
