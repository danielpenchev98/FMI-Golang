package rest

import (
	"log"
	"net/http"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/api/common"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/api/common/response"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/auth"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/dao"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
	val "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//UamEndpoint - rest endpoint for configuration of the user access management
type UamEndpoint interface {
	CreateUser(*gin.Context)
	DeleteUser(*gin.Context)
	Login(*gin.Context)

	CreateGroup(*gin.Context)
	AddMember(*gin.Context)
	RevokeMembership(*gin.Context)
}

//UamEndpointImpl - implementation of UamEndpoint
type UamEndpointImpl struct {
	uamDAO     dao.UamDAO
	jwtCreator auth.JwtCreator
	validator  val.Validator
}

//RegistrationResponse is returned to the client when the registration was successfull
//it contains the statius of hist request and the jwt token
type RegistrationResponse struct {
	Status   int    `json:"status"`
	JWTToken string `json:"jwt_token"`
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

type GroupCreationRequest struct {
	GroupName string `json:"group_name"`
}

type GroupMembershipRequest struct {
	GroupName string `json:"group_name"`
	Username  string `json:"username"`
}

//NewUamEndPointImpl - function for creation an instance of UamEndpointImpl
func NewUamEndPointImpl(uamDAO dao.UamDAO, creator auth.JwtCreator, validator val.Validator) *UamEndpointImpl {
	return &UamEndpointImpl{
		uamDAO:     uamDAO,
		jwtCreator: creator,
		validator:  validator,
	}
}

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
		log.Printf("Problem encryption of use password during the registration. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	if err = e.uamDAO.CreateUser(rq.Username, string(hashedPassword)); err != nil {
		log.Printf("Problem crearing the user in the db. Reason: %v\n", err)
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
		log.Println("Problem retieval of userID from context.")
		common.SendErrorResponse(c, myerr.NewServerError("Cannot retrieve the user id"))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	if err := e.uamDAO.DeleteUser(uint(userID)); err != nil {
		log.Printf("Problem with deletion of user. Reason: %v\n", err)
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
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	user, err := e.uamDAO.GetUser(request.Username)
	if err != nil {
		if _, ok := err.(*myerr.ItemNotFoundError); ok {
			common.SendErrorResponse(c, myerr.NewClientError("Invalid credentials"))
		} else {
			log.Printf("Problem with Login. Reason: %v\n", err)
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
		log.Printf("Problem with generating Jwt token in the login logic. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Status: http.StatusCreated,
		Token:  signedToken,
	})
}

func (e *UamEndpointImpl) CreateGroup(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		log.Println("Problem retieval of userID from context.")
		common.SendErrorResponse(c, myerr.NewServerError("Cannot retrieve the user id"))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	var rq GroupCreationRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	//TODO should rename this function - maybe? or create specific function for the group name
	if err := e.validator.ValidateUsername(rq.GroupName); err != nil {
		err = myerr.NewClientErrorWrap(err, "Problem with the group name")
		log.Printf("Problem with validation of the user registration input. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	if err := e.uamDAO.CreateGroup(userID, rq.GroupName); err != nil {
		log.Printf("Problem with creation of group. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

func (e *UamEndpointImpl) AddMember(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		log.Println("Problem retieval of userID from context.")
		common.SendErrorResponse(c, myerr.NewServerError("Cannot retrieve the user id"))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	var rq GroupMembershipRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	//TODO should rename this function - maybe? or create specific function for the group name
	if err := e.validator.ValidateUsername(rq.GroupName); err != nil {
		err = myerr.NewClientErrorWrap(err, "Problem with the group name")
		common.SendErrorResponse(c, err)
		return
	} else if err = e.validator.ValidateUsername(rq.Username); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	if err := e.uamDAO.AddUserToGroup(userID, rq.Username, rq.GroupName); err != nil {
		log.Printf("Problem with creation of group. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

func (e *UamEndpointImpl) RevokeMembership(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		log.Println("Problem retieval of userID from context.")
		common.SendErrorResponse(c, myerr.NewServerError("Cannot retrieve the user id"))
		return
	}

	var userID uint
	if userID, ok = id.(uint); !ok {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user ID"))
	}

	var rq GroupMembershipRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	//TODO should rename this function - maybe? or create specific function for the group name
	if err := e.validator.ValidateUsername(rq.GroupName); err != nil {
		err = myerr.NewClientErrorWrap(err, "Problem with the group name")
		common.SendErrorResponse(c, err)
		return
	} else if err = e.validator.ValidateUsername(rq.Username); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	if err := e.uamDAO.RemoveUserFromGroup(userID, rq.Username, rq.GroupName); err != nil {
		log.Printf("Problem with membership revoke. Reason: %v\n", err)
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, response.BasicResponse{
		Status: http.StatusOK,
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
