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
	"golang.org/x/crypto/bcrypt"
)

var router = gin.Default()
var uamDAO dao.UamDAO

func init() {
	var err error
	uamDAO, err = dao.NewUamDAOImpl()
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot create DAO object. Reason: %v", err))
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

//Set requirement for only json request bodies
//Maybe the while validation procedure to be encapsualted in the registrationvalidator???
func createUser(c *gin.Context) {
	var rq request.RegistrationRequest

	//Decide what exactly to return as response -> custom message + 400 or?
	if err := c.ShouldBindJSON(&rq); err != nil {
		sendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	if err := validateRegistration(rq); err != nil {
		sendErrorResponse(c, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rq.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	var userID uint
	userID, err = uamDAO.CreateUser(rq.Username, string(hashedPassword))
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.RegistrationResponse{
		StatusCode: http.StatusCreated,
		JWTToken:   fmt.Sprintf("%d", userID),
	})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		sendErrorResponse(c, myerr.NewClientError("Invalid type of id"))
		return
	}

	err = uamDAO.DeleteUser(uint(id))
	if err != nil {
		sendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Status: http.StatusOK,
	})
}

func sendErrorResponse(c *gin.Context, err error) {
	errorCode, errorMsg := getErrorResponseArguments(err)
	c.JSON(errorCode, response.ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
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
