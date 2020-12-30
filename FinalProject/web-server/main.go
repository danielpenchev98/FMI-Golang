package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/user/web-server/api/type/request"
	"example.com/user/web-server/api/type/response"
	"example.com/user/web-server/internal/db/dao"
	myerr "example.com/user/web-server/internal/error"
	val "example.com/user/web-server/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

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
	var rq request.RegistrationRequest

	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorCode: http.StatusUnprocessableEntity,
			ErrorMsg:  fmt.Sprintf("Cannot register user. Reason %v", err),
		})
		return
	}

	if err := validateRegistration(rq); err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, response.ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	userID, err := uamDAO.CreateUser(rq.Username, rq.Password)

	if err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, response.ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, response.RegistrationResponse{
		StatusCode: http.StatusCreated,
		JWTToken:   string(userID),
	})
}

func validateRegistration(rq request.RegistrationRequest) error {
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
		c.JSON(errorCode, response.ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	err = uamDAO.DeleteUser(uint(id))
	if err != nil {
		errorCode, errorMsg := getErrorResponseArguments(err)
		c.JSON(errorCode, response.ErrorResponse{
			ErrorCode: errorCode,
			ErrorMsg:  errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Status: http.StatusOK,
	})
}
