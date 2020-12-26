package dao

import (
	"example.com/user/web-server/pkg/db"
	"example.com/user/web-server/pkg/db/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//To write interface

type UamDAO struct {
	dbConn *gorm.DB
}

func NewUamDAO() (*UamDAO, error) {
	conn, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}
	return &UamDAO{dbConn: conn}, nil
}

func (d *UamDAO) CreateUser(username string, password string) (uint, error) {
	d.dbConn.AutoMigrate(&models.User{})

	var count int64
	result := d.dbConn.Table("users").Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return 0, errors.Wrapf(result.Error, "Cannot check if the user already exists")
	} else if count > 0 {
		return 0, errors.New("A user with the same username exists")
	}

	user := models.User{
		Username: username,
		Password: password,
	}

	if result := d.dbConn.Create(&user); result.Error != nil {
		return 0, errors.Wrapf(result.Error, "Problem with the creation of new user")
	}

	return user.ID, nil
}

//this userID will be saved in the Token
func (d *UamDAO) DeleteUser(userID uint) error {

	var count int64
	result := d.dbConn.Table("users").Where("id = ?", userID).Count(&count)

	if result.Error != nil {
		return errors.Wrapf(result.Error, "Cannot check if the user exists")
	} else if count == 0 {
		return errors.New("User not found")
	}

	if result = d.dbConn.Unscoped().Delete(&models.User{}, userID); result.Error != nil {
		return errors.Wrapf(result.Error, "Problem with the deletion of the user")
	}
	return nil
}
