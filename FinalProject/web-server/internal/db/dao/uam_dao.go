package dao

import (
	"example.com/user/web-server/internal/db"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"gorm.io/gorm"
)

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
	//d.dbConn.AutoMigrate(&models.User{})

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
func (d *UamDAO) DeleteUser(userID uint) error {

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
