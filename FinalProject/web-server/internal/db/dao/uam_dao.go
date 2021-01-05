package dao

import (
	"example.com/user/web-server/internal/db"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"gorm.io/gorm"
)

//go:generate mockgen --source=uam_dao.go --destination uam_dao_mocks/uam_dao.go --package uam_dao_mocks

type UamDAO interface {
	CreateUser(string, string) (uint, error)
	CheckIfUserExists(string,string) (bool,error)
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

func (d *UamDAOImpl) CreateUser(username string, password string) (uint, error) {
	var count int64
	result := d.dbConn.Table("users").Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return 0, myerr.NewServerError("Problem with the lookup of users", result.Error)
	} else if count > 0 {
		return 0, myerr.NewClientError("A user with the same username exists")
	}

	user := models.User{
		Username: username,
		Password: password,
	}

	if result := d.dbConn.Create(&user); result.Error != nil {
		return 0, myerr.NewServerError("Problem with the creation of new user", result.Error)
	}

	return user.ID, nil
}

//this userID will be saved in the Token
func (d *UamDAOImpl) DeleteUser(userID uint) error {

	var count int64
	result := d.dbConn.Table("users").Where("id = ?", userID).Count(&count)

	if result.Error != nil {
		return myerr.NewServerError("Problem with the lookup if user exists", result.Error)
	} else if count == 0 {
		return myerr.NewItemNotFoundError("User with that id does not exist")
	}

	if result = d.dbConn.Unscoped().Delete(&models.User{}, userID); result.Error != nil {
		return myerr.NewServerError("Problem with the deletion of the user", result.Error)
	}
	return nil
}

func (d *UamDAOImpl) CheckIfUserExists(username string, password string) (bool,error) {
	var count int64
	result := d.dbConn.Table("users").
		Where("username = ?", username).
		Where("password = ?", password).
		Count(&count)
	
	if result.Error != nil {
		return false, myerr.NewServerError("Problem with the lookup if user exists", result.Error)
	} else if count == 0 {
		return false, nil
	}
	return true, nil
}