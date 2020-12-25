package main

import (
	"fmt"
	"log"
	"net/http"

	val "example.com/user/web-server/pkg/validator"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMsg  string `json:"message"`
}

type SuccessResponse struct {
	Status int `json:"status"`
}

var router = gin.Default()

func main() {
	router.POST("/user/create", CreateUser)
	log.Fatal(router.Run(":8080"))
}

func CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
	}

	validator := val.NewValidator()
	if err := validator.ValidateUsername(u.Username); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	if err := validator.ValidatePassword(u.Password); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Status: http.StatusCreated,
	})
}
