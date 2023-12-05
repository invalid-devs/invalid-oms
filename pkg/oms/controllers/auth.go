package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/invalid-devs/invalid-oms/pkg/oms/middleware"
	"github.com/invalid-devs/invalid-oms/pkg/oms/models"
	"github.com/invalid-devs/invalid-oms/pkg/oms/response"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type SignInRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInResponseDTO struct {
	Token string `json:"token"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	var requestDTO SignInRequestDTO
	err := response.Validate(r, &requestDTO)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	var user models.User

	username := strings.Trim(requestDTO.Username, " ")
	username = strings.ToLower(username)

	result := db.First(&user, "username = ?", username)

	if result.Error != nil {
		msg := response.Error{
			Message: "Invalid credentials!",
		}
		response.JSON(w, http.StatusBadRequest, msg)
		return
	}

	password := strings.Trim(requestDTO.Password, " ")

	if user.Password != password {
		msg := response.Error{
			Message: "Invalid credentials!",
		}
		response.JSON(w, http.StatusBadRequest, msg)
		return
	}

	expire := time.Now().AddDate(0, 0, 1).UTC().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"expire":  expire,
	})
	tokenString, _ := token.SignedString(middleware.Secret)

	responseDTO := SignInResponseDTO{
		Token: tokenString,
	}
	response.JSON(w, http.StatusCreated, responseDTO)
}
