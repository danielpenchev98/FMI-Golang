package rest

import (
	"net/http"

	"example.com/user/web-server/api/common"
	"example.com/user/web-server/api/common/request"
	"example.com/user/web-server/api/common/response"
	"example.com/user/web-server/internal/auth"
	"example.com/user/web-server/internal/db/dao"
	myerr "example.com/user/web-server/internal/error"
	val "example.com/user/web-server/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
}

//LoginRequest - request representation for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//LoginResponse - when the login is succesfull a JWT is sent to the user
type LoginResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

//NewUamEndPointImpl - function for creation an instance of UamEndpointImpl
func NewUamEndPointImpl(uamDAO dao.UamDAO, creator auth.JwtCreator) *UamEndpointImpl {
	return &UamEndpointImpl{
		uamDAO:     uamDAO,
		jwtCreator: creator,
	}
}

//Set requirement for only json request bodies
//Maybe the while validation procedure to be encapsualted in the registrationvalidator???
//CreateUser - handler of request for creation of new user
func (e *UamEndpointImpl) CreateUser(c *gin.Context) {
	var rq request.RegistrationRequest

	//Decide what exactly to return as response -> custom message + 400 or?
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	if err := validateRegistration(rq); err != nil {
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

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Status: http.StatusCreated,
	})
}

//DeleteUser - handler for request of deletetion of user
func (e *UamEndpointImpl) DeleteUser(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		common.SendErrorResponse(c, myerr.NewServerError("", errors.New("Cannot retrieve the user id")))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	if err := e.uamDAO.DeleteUser(uint(userID)); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Status: http.StatusOK,
	})
}

//Login - handler for request of login of user
func (e *UamEndpointImpl) Login(c *gin.Context) {
	var payload LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	user, err := e.uamDAO.GetUser(payload.Username)
	if err != nil {
		if _, ok := err.(*myerr.ItemNotFoundError); ok {
			//log the error
			common.SendErrorResponse(c, myerr.NewClientError("Invalid credentials"))
		} else {
			common.SendErrorResponse(c, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid credentials"))
		return
	}

	signedToken, err := e.jwtCreator.GenerateToken(user.ID)
	if err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Status: http.StatusOK,
		Token:  signedToken,
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
