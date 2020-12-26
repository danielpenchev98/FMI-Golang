package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/user/web-server/pkg/db/dao"
	myerr "example.com/user/web-server/pkg/errors"
	val "example.com/user/web-server/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	if err := validateRegistration(rq); err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	userID, err := uamDAO.CreateUser(rq.Username, rq.Password)

	if err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, RegistrationResponse{
		Status: http.StatusCreated,
		ID:     userID,
	})
}

func validateRegistration(rq RegistrationRequest) error {
	validator := val.NewValidator()
	if err := validator.ValidateUsername(rq.Username); err != nil {
		return errors.Wrapf(err, "Problem with the username")
	}

	if err := validator.ValidatePassword(rq.Password); err != nil {
		return errors.Wrapf(err, "Problem with the password")
	}
	return nil
}

func getErrorResponseArguments(err error) (errorCode int, errorMsg string) {
	switch err.(type) {
	case *myerr.ClientError:
		errorCode = http.StatusBadRequest
		errorMsg = fmt.Sprintf("Invalid request. Reason :%s", err.Error())
	case *myerr.ItemNotFoundError:
		errorCode = http.StatusNotFound
		errorMsg = err.Error()
	default:
		errorCode = http.StatusInternalServerError
		errorMsg = fmt.Sprintf("Problem with the server, please try again later")
		fmt.Println(err)
	}
	return
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorCode, errorMsg := getErrorResponseArguments(myerr.NewClientError("Invalid type of id"))
		c.JSON(errorCode, ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	err = uamDAO.DeleteUser(uint(id))
	if err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Status: http.StatusOK,
	})
}
