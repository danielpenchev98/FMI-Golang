package dao

import (
	"errors"
	"fmt"
	"log"

	"example.com/user/web-server/internal/db"
	"example.com/user/web-server/internal/db/models"
	myerr "example.com/user/web-server/internal/error"
	"gorm.io/gorm"
)

//go:generate mockgen --source=uam_dao.go --destination uam_dao_mocks/uam_dao.go --package uam_dao_mocks

//UamDAO - interface for working with the Database in regards to the User Access Management
type UamDAO interface {
	Migrate() error
	CreateUser(string, string) error
	GetUser(string) (models.User, error)
	DeleteUser(uint) error

	CreateGroup(uint, string) error
	AddUserToGroup(uint, string, string) error
	RemoveUserFromGroup(uint, string, string) error
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

func (d *UamDAOImpl) Migrate() error {
	return d.dbConn.AutoMigrate(models.User{}, models.Group{}, models.Membership{})
}

//CreateUser - creates a new user in the database, given username and password (encrypted)
func (d *UamDAOImpl) CreateUser(username string, password string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
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
	if result = d.dbConn.Delete(&models.User{}, userID); result.Error != nil {
		log.Printf("Problem with the request to delete user with [%d] in the database. Reason: %v\n", userID, result.Error)
		return myerr.NewServerErrorWrap(result.Error, "Problem with the deletion of the user")
	}
	log.Printf("User with id [%d] is deleted\n", userID)

	return nil

}

//GetUser - fetches information about an existing user
func (d *UamDAOImpl) GetUser(username string) (models.User, error) {
	return getUserWithConn(d.dbConn, username)
}

func (d *UamDAOImpl) CreateGroup(userID uint, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
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
}

func (d *UamDAOImpl) GetGroup(groupName string) (models.Group, error) {
	return getGroupWithConn(d.dbConn, groupName)
}

func (d *UamDAOImpl) AddUserToGroup(ownerID uint, username string, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
		var (
			count int64
			group models.Group
			user  models.User
			err   error
		)

		group, err = getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		} else if group.OwnerID != ownerID {
			return myerr.NewClientError("Only the group owner can add members to the group")
		}

		user, err = getUserWithConn(tx, username)
		if err != nil {
			return err
		}

		result := tx.Table("memberships").
			Where("group_id = ?", group.ID).
			Where("user_id = ?", user.ID).
			Count(&count)

		if result.Error != nil {
			log.Printf("Problem with the request to check if membership for user with id [%d] in group with id [%d] exists in db. Reason: %v\n",
				user.ID, group.ID, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of membership")
		} else if count != 0 {
			return myerr.NewClientError("The user is already a member of the group")
		}

		membership := models.Membership{
			GroupID: group.ID,
			UserID:  user.ID,
		}

		log.Printf("Creating membership for user with id [%d] in group with id [%d]", membership.UserID, membership.GroupID)
		if result := tx.Create(&membership); result.Error != nil {
			log.Printf("Problem with the request to create membership for user with id [%d] in group with id [%d] in the database. Reason: %v\n",
				membership.GroupID, membership.UserID, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new membership")
		}
		log.Printf("Membership for user with id [%d] in group id [%d] created", membership.UserID, membership.GroupID)

		return nil
	})
}

func (d *UamDAOImpl) RemoveUserFromGroup(currUserID uint, username string, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
		var (
			group models.Group
			user  models.User
			err   error
		)

		group, err = getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		}

		user, err = getUserWithConn(tx, username)
		if err != nil {
			return err
		}

		if group.OwnerID != currUserID && user.ID != currUserID {
			return myerr.NewClientError("Only the owner of the group can revoke membership of other members")
		}

		membership := models.Membership{
			GroupID: group.ID,
			UserID:  user.ID,
		}

		fmt.Println(group.ID)
		fmt.Println(user.ID)

		log.Printf("Revolking membership for user with id [%d] in group with id [%d]", membership.UserID, membership.GroupID)
		result := tx.Where("user_id = ?", user.ID).
			Where("group_id = ?", group.ID).
			Delete(&membership)

		if result.Error != nil {
			log.Printf("Problem with the request to revoke membership for user with id [%d] in group with id [%d] in the database. Reason: %v\n",
				membership.GroupID, membership.UserID, result.Error)
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new membership")
		} else if result.RowsAffected == 0 {
			return myerr.NewClientError("Membership not found")
		}
		log.Printf("Membership for user with id [%d] in group id [%d] is revoked", membership.UserID, membership.GroupID)

		return nil
	})
}

func getUserWithConn(dbConn *gorm.DB, username string) (models.User, error) {
	var user models.User

	result := dbConn.Table("users").
		Where("username = ?", username).
		Take(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, myerr.NewItemNotFoundError("User does not exist")
	} else if result.Error != nil {
		return user, myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if user exists")
	}

	return user, nil
}

func getGroupWithConn(dbConn *gorm.DB, groupName string) (models.Group, error) {
	var group models.Group

	result := dbConn.Table("groups").
		Where("name = ?", groupName).
		Take(&group)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return group, myerr.NewItemNotFoundError(fmt.Sprintf("Group [%s] does not exist", groupName))
	} else if result.Error != nil {
		return group, myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if group exists")
	}

	return group, nil
}
