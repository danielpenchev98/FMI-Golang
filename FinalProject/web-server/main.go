package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/user/web-server/pkg/db/dao"
	val "example.com/user/web-server/pkg/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegistrationRequest struct {
	Username string
	Password string
}

type UserDTO struct {
	gorm.Model
	Username string
	Password string
}

type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMsg  string `json:"message"`
}

type RegistrationResponse struct {
	Status int  `json:"status"`
	ID     uint `json:"id"`
}

type SuccessResponse struct {
	Status int `json:"status"`
}

var router = gin.Default()
var uamDAO *dao.UamDAO

func init() {
	var err error
	uamDAO, err = dao.NewUamDAO()
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot create DAO object. Reason: %v", err))
		panic("Problem with DAO")
	}

	router = gin.Default()
}

func main() {

	v1 := router.Group("/v1")
	{
		v1.POST("/user/register", createUser)
		v1.POST("/user/delete/:id", deleteUser)
	}
	log.Fatal(router.Run(":8080"))
}

//Maybe the while validation procedure to be encapsualted in the registrationvalidator???
func createUser(c *gin.Context) {
	var rq RegistrationRequest

	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	validator := val.NewValidator()
	if err := validator.ValidateUsername(rq.Username); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	if err := validator.ValidatePassword(rq.Password); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	userID, err := uamDAO.CreateUser(rq.Username, rq.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Problem with the server, please try again later",
		})
		return
	}

	c.JSON(http.StatusCreated, RegistrationResponse{
		Status: http.StatusCreated,
		ID:     userID,
	})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Invalid type of id",
		})
		return
	}

	err = uamDAO.DeleteUser(uint(id))
	if err != nil {
		//if the error type is pariculary -> user not found
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Problem with the deletion of the user",
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Status: http.StatusOK,
	})
}
