package dao

import (
	"errors"

	"example.com/user/web-server/internal/db"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"gorm.io/gorm"
)

//go:generate mockgen --source=uam_dao.go --destination uam_dao_mocks/uam_dao.go --package uam_dao_mocks

type UamDAO interface {
	CreateUser(string, string) error
	GetUser(string) (models.User, error)
	DeleteUser(uint) error
}

type UamDAOImpl struct {
	dbConn *gorm.DB
}

func NewUamDAOImpl() (*UamDAOImpl, error) {
	conn, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}
	return &UamDAOImpl{dbConn: conn}, nil
}

func (d *UamDAOImpl) CreateUser(username string, password string) error {
	var count int64
	result := d.dbConn.Table("users").Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of users")
	} else if count > 0 {
		return myerr.NewClientError("A user with the same username exists")
	}

	user := models.User{
		Username: username,
		Password: password,
	}

	if result := d.dbConn.Create(&user); result.Error != nil {
		return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new user")
	}

	return nil
}

//this userID will be saved in the Token
func (d *UamDAOImpl) DeleteUser(userID uint) error {

	var count int64
	result := d.dbConn.Table("users").Where("id = ?", userID).Count(&count)

	if result.Error != nil {
		return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if user exists")
	} else if count == 0 {
		return myerr.NewItemNotFoundError("User with that id does not exist")
	}

	if result = d.dbConn.Unscoped().Delete(&models.User{}, userID); result.Error != nil {
		return myerr.NewServerErrorWrap(result.Error, "Problem with the deletion of the user")
	}
	return nil
}

func (d *UamDAOImpl) GetUser(username string) (models.User, error) {
	var user models.User

	result := d.dbConn.Table("users").
		Where("username = ?", username).
		Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, myerr.NewItemNotFoundError("User with those credentials does not exist")
	} else if result.Error != nil {
		return user, myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if user exists")
	}

	return user, nil
}
