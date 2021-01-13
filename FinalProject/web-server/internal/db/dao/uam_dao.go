package dao

import (
	"errors"
	"log"

	"example.com/user/web-server/internal/db"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"gorm.io/gorm"
)

//go:generate mockgen --source=uam_dao.go --destination uam_dao_mocks/uam_dao.go --package uam_dao_mocks

//UamDAO - interface for working with the Database in regards to the User Access Management
type UamDAO interface {
	CreateUser(string, string) error
	GetUser(string) (models.User, error)
	DeleteUser(uint) error

	CreateGroup(uint, string) error
}

//UamDAOImpl - implementation of UamDAO
type UamDAOImpl struct {
	dbConn *gorm.DB
}

//NewUamDAOImpl - function for creation an instance of UamDAOImpl
func NewUamDAOImpl() (*UamDAOImpl, error) {
	conn, err := db.GetDBConn()
	if err != nil {
		return nil, err
	}
	return &UamDAOImpl{dbConn: conn}, nil
}

//CreateUser - creates a new user in the database, given username and password (encrypted)
func (d *UamDAOImpl) CreateUser(username string, password string) error {
	err := d.dbConn.Transaction(func(tx *gorm.DB) error {
		var count int64
		result := tx.Table("users").Where("username = ?", username).Count(&count)

		if result.Error != nil {
			log.Printf("Problem with request to check if user with username [%s] exists in the database. Reason: %v\n", username, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of users")
		} else if count > 0 {
			return myerr.NewClientError("A user with the same username exists")
		}

		user := models.User{
			Username: username,
			Password: password,
		}

		log.Printf("Creating user with username [%s]", username)
		if result := tx.Create(&user); result.Error != nil {
			log.Printf("Problem with the request to create user with username [%s] in the database. Reason: %v\n", username, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new user")
		}
		log.Printf("User with username [%s] created", username)

		return nil
	})
	return err
}

//DeleteUser - deletes user given an id of the user
func (d *UamDAOImpl) DeleteUser(userID uint) error {
	var count int64
	result := d.dbConn.Table("users").Where("id = ?", userID).Count(&count)

	if result.Error != nil {
		log.Printf("Couldnt search for a user with id [%d] in the database. Reason: %v\n", userID, result.Error)
		return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if user exists")
	} else if count == 0 {
		return myerr.NewItemNotFoundError("User with that id does not exist")
	}

	log.Printf("Deleting user with id [%d]\n", userID)
	if result = d.dbConn.Unscoped().Delete(&models.User{}, userID); result.Error != nil {
		log.Printf("Problem with the request to delete user with [%d] in the database. Reason: %v\n", userID, result.Error)
		return myerr.NewServerErrorWrap(result.Error, "Problem with the deletion of the user")
	}
	log.Printf("User with id [%d] is deleted\n", userID)

	return nil

}

//GetUser - fetches information about an existing user
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

func (d *UamDAOImpl) CreateGroup(userID uint, groupName string) error {
	err := d.dbConn.Transaction(func(tx *gorm.DB) error {
		var count int64
		result := tx.Table("groups").Where("name = ?", groupName).Count(&count)

		if result.Error != nil {
			log.Printf("Problem with request to check if group [%s] exists in the database. Reason: %v\n", groupName, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of groups")
		} else if count > 0 {
			return myerr.NewClientError("A group with the same name exists")
		}

		group := models.Group{
			Name:    groupName,
			OwnerID: userID,
		}

		log.Printf("Creating group [%s] with owner [%d]\n", groupName, userID)
		if result := tx.Create(&group); result.Error != nil {
			log.Printf("Problem with the request to create group [%s] with owner [%d] in the database. Reason: %v\n", groupName, userID, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new group")
		}
		log.Printf("Group with name [%s] and owner [%d] created\n", groupName, userID)

		membership := models.Membership{
			UserID:  userID,
			GroupID: group.ID,
		}

		//its usedless to check if the membership already exists, because basically the group is created in this transaction
		log.Printf("Creating membership of user [%d] for group [%d]\n", userID, group.ID)
		if result := tx.Create(&membership); result.Error != nil {
			log.Printf("Problem with the request to create membership of user [%d] for group [%d] in the database. Reason: %v\n", userID, group.ID, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of membership")
		}
		log.Printf("Membership of user [%d] for group [%d] is created\n", userID, group.ID)

		return nil
	})
	return err
}

/*
func (d *UamDAOImpl) addUserToGroup(ownerID uint, userID uint, groupName string) error {
	var group models.Group
	result := d.dbConn.Table("groups").
		Where("name = ? ", groupName).
		Where("ownerID = ?", ownerID)
		.Count(&count)

	if result.Error != nil {
		log.Printf("Problem with request to check if group [%s] exists in the database. Reason: %v\n", username, result.Error)
		return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of groups")
	} else if count == 0 {
		return myerr.NewClientError("Invalid group")
	}

	group := models.Membership{
		Name:    groupName,
		OwnerID: userID,
	}

	log.Printf("Creating group [%s] with owner [%d]", groupName, userID)
	if result := d.dbConn.Create(&group); result.Error != nil {
		log.Printf("Problem with the request to create group [%s] with owner [%d] in the database. Reason: %v\n", groupName, ownerID, result.Error)
		return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new group")
	}
	log.Printf("Group with name [%s] and owner [%d] created", groupName, userID)

	return nil
}*/
