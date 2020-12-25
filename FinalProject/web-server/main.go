package main

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

type User struct {
	username string
	password string
}

var router = gin.Default()

func main() {
	//router.POST("/user/create")
}

func ValidateUsername(username string) error {
	if length := len(username); length < 8 || length > 20 {
		return errors.New("Username should be between 8 and 20 symbols")
	}

	matched, _ := regexp.Match(`^[^a-zA-Z].+`, []byte(username))
	if matched {
		return errors.New("Username should always begin only with a letter")
	}

	matched, _ = regexp.Match(`^[-_0-9a-zA-Z]+$`, []byte(username))
	if !matched {
		return errors.New("Username cannot contain special symbols except \"-\" and \"_\"")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 10 {
		return errors.New("Password should be greater than 9 symbols")
	}

	matched, _ := regexp.Match(`^[^0-9]+$`, []byte(password))
	if matched {
		return errors.New("Password should contain atleast one number")
	}

	matched, _ = regexp.Match(`^[0-9a-zA-Z]+$`, []byte(password))
	if matched {
		return errors.New("Password should contain atleast one special char")
	}

	return nil
}

/*func CreateUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "BALABAL")
	}
}*/
