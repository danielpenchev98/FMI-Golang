package rest

import (
	"net/http"
	"os"

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
	DeleteGroup(*gin.Context)
}

//UamEndpointImpl - implementation of UamEndpoint
type UamEndpointImpl struct {
	uamDAO     dao.UamDAO
	jwtCreator auth.JwtCreator
	validator  val.Validator
	groupsDir  string
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

//GroupPayload - request payload, containing the group name
type GroupPayload struct {
	GroupName string `json:"group_name"`
}

//GroupMembershipPayload - request payload, containing the group name and username
type GroupMembershipPayload struct {
	GroupPayload
	Username string `json:"username"`
}

//NewUamEndPointImpl - function for creation an instance of UamEndpointImpl
func NewUamEndPointImpl(uamDAO dao.UamDAO, creator auth.JwtCreator, validator val.Validator, groupsDir string) *UamEndpointImpl {
	return &UamEndpointImpl{
		uamDAO:     uamDAO,
		jwtCreator: creator,
		validator:  validator,
		groupsDir:  groupsDir,
	}
}

//Maybe the while validation procedure to be encapsualted in the registrationvalidator???
//CreateUser - handler for user creation request
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
		err = myerr.NewServerErrorWrap(err, "Problem encryption of password during the registration.")
		common.SendErrorResponse(c, err)
		return
	}

	err = e.uamDAO.CreateUser(rq.Username, string(hashedPassword))
	if _, ok := err.(*myerr.ClientError); ok {
		common.SendErrorResponse(c, err)
		return
	} else if err != nil {
		err = myerr.NewServerErrorWrap(err, "Problem crearing the user in the db.")
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

//DeleteUser - handler for user deletion request
func (e *UamEndpointImpl) DeleteUser(c *gin.Context) { //Not tested yet -> gin.Context cannot be mocked
	userID, err := common.GetIDFromContext(c)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	if err = e.uamDAO.DeleteUser(uint(userID)); err != nil {
		err = myerr.NewServerErrorWrap(err, "Problem with deletion of user.")
		common.SendErrorResponse(c, myerr.NewServerErrorWrap(err, "Problem with deleting user"))
		return
	}

	c.JSON(http.StatusOK, response.BasicResponse{
		Status: http.StatusOK,
	})
}

//Login - handler for user login request
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
			err = myerr.NewServerErrorWrap(err, "Problem with Login.")
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
		err = myerr.NewServerErrorWrap(err, "Problem with generating Jwt token in the login logic.")
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Status: http.StatusCreated,
		Token:  signedToken,
	})
}

//CreateGroup - handler for group creation request
func (e *UamEndpointImpl) CreateGroup(c *gin.Context) {
	userID, err := common.GetIDFromContext(c)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	var rq GroupPayload
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	//TODO should rename this function - maybe? or create specific function for the group name
	if err = e.validator.ValidateUsername(rq.GroupName); err != nil {
		err = myerr.NewClientErrorWrap(err, "Problem with the group name")
		common.SendErrorResponse(c, err)
		return
	}

	groupDir := e.groupsDir + "/" + rq.GroupName
	if _, err := os.Stat(groupDir); !os.IsNotExist(err) {
		common.SendErrorResponse(c, myerr.NewClientError("Problem with creation of group. Reason: Group already exists"))
		return
	} else if err = os.Mkdir(groupDir, 0755); err != nil {
		common.SendErrorResponse(c, myerr.NewServerError("Problem with creation of directory"))
		return
	}

	if err = e.uamDAO.CreateGroup(userID, rq.GroupName); err != nil {
		os.RemoveAll(groupDir)
		common.SendErrorResponse(c, myerr.NewServerErrorWrap(err, "Problem with creation of group."))
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

//AddMember - handler for membership creation request
func (e *UamEndpointImpl) AddMember(c *gin.Context) {
	userID, err := common.GetIDFromContext(c)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	var rq GroupMembershipPayload
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	if err = e.uamDAO.AddUserToGroup(userID, rq.Username, rq.GroupName); err != nil {
		err = myerr.NewServerErrorWrap(err, "Problem with creation of group.")
		common.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.BasicResponse{
		Status: http.StatusCreated,
	})
}

//RevokeMembership - handler for membership deletion request
func (e *UamEndpointImpl) RevokeMembership(c *gin.Context) {
	userID, err := common.GetIDFromContext(c)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	var rq GroupMembershipPayload
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	err = e.uamDAO.RemoveUserFromGroup(userID, rq.Username, rq.GroupName)

	if err == nil {
		c.JSON(http.StatusOK, response.BasicResponse{
			Status: http.StatusOK,
		})
	}

	if _, ok := err.(*myerr.ServerError); ok {
		err = myerr.NewServerErrorWrap(err, "Couldnt remove membership")
	}

	common.SendErrorResponse(c, err)
}

//DeleteGroup - handler for group deletion request
func (e *UamEndpointImpl) DeleteGroup(c *gin.Context) {
	userID, err := common.GetIDFromContext(c)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	var rq GroupPayload
	if err = c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	err = e.uamDAO.DeleteGroup(userID, rq.GroupName)
	if err == nil {
		c.JSON(http.StatusOK, response.BasicResponse{
			Status: http.StatusOK,
		})
	}

	if _, ok := err.(*myerr.ServerError); ok {
		err = myerr.NewServerErrorWrap(err, "Problem with deletion of group.")
	}

	common.SendErrorResponse(c, err)
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
