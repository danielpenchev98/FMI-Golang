package rest

import (
	"fmt"
	"net/http"

	"example.com/user/web-server/api/common"
	"example.com/user/web-server/api/common/response"
	"example.com/user/web-server/internal/auth"
	"example.com/user/web-server/internal/db/dao"
	myerr "example.com/user/web-server/internal/error"
	val "example.com/user/web-server/internal/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//UamEndpoint - rest endpoint for configuration of the user access management
type UamEndpoint interface {
	CreateUser(*gin.Context)
	DeleteUser(*gin.Context)
	Login(*gin.Context)
}

//UamEndpointImpl - implementation of UamEndpoint
type UamEndpointImpl struct {
	uamDAO     dao.UamDAO
	jwtCreator auth.JwtCreator
	validator  val.Validator
}

//RegistrationResponse is returned to the client when the registration was successfull
//it contains the statius of hist request and the jwt token
type RegistrationResponse struct{
	Status int    `json:"status"`
	JWTToken   string `json:"jwt_token"`
}

//RequestWithCredentials - request representation for login
type RequestWithCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//LoginResponse - when the login is succesfull a JWT is sent to the user
type LoginResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

//NewUamEndPointImpl - function for creation an instance of UamEndpointImpl
func NewUamEndPointImpl(uamDAO dao.UamDAO, creator auth.JwtCreator, validator val.Validator) *UamEndpointImpl {
	return &UamEndpointImpl{
		uamDAO:     uamDAO,
		jwtCreator: creator,
		validator:  validator,
	}
}

//Set requirement for only json request bodies
//Maybe the while validation procedure to be encapsualted in the registrationvalidator???
//CreateUser - handler of request for creation of new user
func (e *UamEndpointImpl) CreateUser(c *gin.Context) {
	var rq RequestWithCredentials

	//Decide what exactly to return as response -> custom message + 400 or?
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	if err := validateRegistration(e.validator, rq); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rq.Password), bcrypt.DefaultCost)
	if err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	if err = e.uamDAO.CreateUser(rq.Username, string(hashedPassword)); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

//DeleteUser - handler for request of deletetion of user
func (e *UamEndpointImpl) DeleteUser(c *gin.Context) { //Not tested yet -> gin.Context cannot be mocked
	id, ok := c.Get("userID")
	if !ok {
		common.SendErrorResponse(c, myerr.NewServerError("Cannot retrieve the user id"))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	if err := e.uamDAO.DeleteUser(uint(userID)); err != nil {
		common.SendErrorResponse(c, myerr.NewServerErrorWrap(err, "Problem with deleting user"))
		return
	}

	c.JSON(http.StatusOK, response.BasicResponse{
		Status: http.StatusOK,
	})
}

//Login - handler for request of login of user
func (e *UamEndpointImpl) Login(c *gin.Context) {
	var request RequestWithCredentials
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	user, err := e.uamDAO.GetUser(request.Username)
	if err != nil {
		if _, ok := err.(*myerr.ItemNotFoundError); ok {
			//log the error
			common.SendErrorResponse(c, myerr.NewClientError("Invalid credentials"))
		} else {
			common.SendErrorResponse(c, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid credentials"))
		return
	}

	signedToken, err := e.jwtCreator.GenerateToken(user.ID)
	if err != nil {
		fmt.Println(err)
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Status: http.StatusCreated,
		Token:  signedToken,
	})
}

func validateRegistration(validator val.Validator, rq RequestWithCredentials) error {
	if err := validator.ValidateUsername(rq.Username); err != nil {
		return myerr.NewClientErrorWrap(err, "Problem with the username")
	}

	if err := validator.ValidatePassword(rq.Password); err != nil {
		return myerr.NewClientErrorWrap(err, "Problem with the password")
	}
	return nil
}
