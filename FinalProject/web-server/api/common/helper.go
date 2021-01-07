package common

import (
	"fmt"
	"net/http"

	"example.com/user/web-server/api/common/response"
	myerr "example.com/user/web-server/internal/error"
	"github.com/gin-gonic/gin"
)

//SendErrorResponse - generic method for sending error response to the user
func SendErrorResponse(c *gin.Context, err error) {
	errorCode, errorMsg := getErrorResponseArguments(err)
	c.JSON(errorCode, response.ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	})
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
