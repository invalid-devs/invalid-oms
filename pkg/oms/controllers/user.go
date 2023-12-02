package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/invalid-devs/invalid-oms/pkg/oms/models"
	"github.com/invalid-devs/invalid-oms/pkg/oms/response"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GetUserResponseDTO struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		msg := response.Error{
			Message: err.Error(),
		}
		response.JSON(w, http.StatusBadRequest, msg)
		return
	}

	var user models.User
	db.First(&user, id)

	responseDTO := GetUserResponseDTO{
		Id:        user.ID,
		Username:  user.Username,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.JSON(w, http.StatusOK, responseDTO)
}

type CreateUserRequestDTO struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type CreateUserResponseDTO struct {
	Id uint `json:"id"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	body, _ := io.ReadAll(r.Body)

	var requestDTO CreateUserRequestDTO
	if err := json.Unmarshal(body, &requestDTO); err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username: requestDTO.Username,
		Password: requestDTO.Password,
		IsActive: true,
	}
	db.Create(&user)

	responseDTO := CreateUserResponseDTO{
		Id: user.ID,
	}
	response.JSON(w, http.StatusCreated, responseDTO)
}
