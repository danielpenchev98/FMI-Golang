package dao

import (
	"errors"
	"fmt"
	"log"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/models"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
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
	MemberExists(uint, uint) (bool, error)
	DeleteGroup(uint, string) error
	GetGroup(string) (models.Group, error)
}

//UamDAOImpl - implementation of UamDAO
type UamDAOImpl struct {
	dbConn *gorm.DB
}

//NewUamDAOImpl - function for creation an instance of UamDAOImpl
func NewUamDAOImpl(dbConn *gorm.DB) *UamDAOImpl {
	return &UamDAOImpl{dbConn: dbConn}
}

//Migrate - function which updates the models(table structure) in db
func (d *UamDAOImpl) Migrate() error {
	return d.dbConn.AutoMigrate(models.User{}, models.Group{}, models.Membership{})
}

//CreateUser - creates a new user in the database, given username and password (encrypted)
func (d *UamDAOImpl) CreateUser(username string, password string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
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

		log.Printf("Creating user with username [%s]", username)
		if result := tx.Create(&user); result.Error != nil {
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
		return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if user exists")
	} else if count == 0 {
		return myerr.NewItemNotFoundError("User with that id does not exist")
	}

	log.Printf("Deleting user with id [%d]\n", userID)
	if result = d.dbConn.Delete(&models.User{}, userID); result.Error != nil {
		return myerr.NewServerErrorWrap(result.Error, "Problem with the deletion of the user from db")
	}
	log.Printf("User with id [%d] is deleted\n", userID)

	return nil

}

//GetUser - fetches information about an existing user
func (d *UamDAOImpl) GetUser(username string) (models.User, error) {
	return getUserWithConn(d.dbConn, username)
}

//CreateGroup - creates a new group for sharing files
func (d *UamDAOImpl) CreateGroup(userID uint, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
		var count int64

		result := tx.Table("groups").Where("name = ?", groupName).Count(&count)
		if result.Error != nil {
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
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of group [%s] in db")
		}
		log.Printf("Group with name [%s] and owner [%d] created\n", groupName, userID)

		membership := models.Membership{
			UserID:  userID,
			GroupID: group.ID,
		}

		//its usedless to check if the membership already exists, because basically the group is created in this transaction
		log.Printf("Creating membership of user [%d] for group [%d]\n", userID, group.ID)
		if result := tx.Create(&membership); result.Error != nil {
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of membership in db")
		}
		log.Printf("Membership of user [%d] for group [%d] is created\n", userID, group.ID)

		return nil
	})
}

//GetGroup - gets information about the group
func (d *UamDAOImpl) GetGroup(groupName string) (models.Group, error) {
	return getGroupWithConn(d.dbConn, groupName)
}

//AddUserToGroup - adds a new member to a specified group
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
		} else if !group.Active {
			return myerr.NewClientError("The group is currently being deleted")
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
			return myerr.NewServerErrorWrap(result.Error, "Problem with the lookup of membership in db")
		} else if count != 0 {
			return myerr.NewClientError("The user is already a member of the group")
		}

		membership := models.Membership{
			GroupID: group.ID,
			UserID:  user.ID,
		}

		log.Printf("Creating membership for user with id [%d] in group with id [%d]", membership.UserID, membership.GroupID)
		if result := tx.Create(&membership); result.Error != nil {
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new membership in db")
		}
		log.Printf("Membership for user with id [%d] in group id [%d] created", membership.UserID, membership.GroupID)

		return nil
	})
}

//DeleteGroup - deletes all memberships and changes the status of the group to non active
func (d *UamDAOImpl) DeleteGroup(currUserID uint, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
		group, err := getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		} else if group.OwnerID != currUserID {
			return myerr.NewClientError("Only the group owner can delete the group")
		} else if !group.Active {
			return myerr.NewClientError("The group is currently being deleted")
		}

		log.Printf("Revolking membership for users in group [%s]", groupName)
		result := tx.Table("memberships").
			Where("group_id = ?", group.ID).Delete(&models.Membership{})
		if result.Error != nil {
			return myerr.NewServerErrorWrap(result.Error, "Problem with deletion of memberships in db")
		}
		log.Printf("Revolked membership for users in group [%s]", groupName)

		log.Printf("Change status of group [%s] to non active\n", groupName)
		if result = tx.Model(&group).Update("active", false); result.Error != nil {
			return myerr.NewServerErrorWrap(result.Error, "Problem with deletion of the group in db")
		}
		log.Printf("Status of group [%s] is set to non active\n", groupName)
		return nil
	})
}

//RemoveUserFromGroup - removes a membership of a user to a specific group
func (d *UamDAOImpl) RemoveUserFromGroup(currUserID uint, username string, groupName string) error {
	return d.dbConn.Transaction(func(tx *gorm.DB) error {
		var (
			group models.Group
			user  models.User
			err   error
		)

		/*
			SHOULD THIS BISNESS LOGIC BE HERE AT ALL?????????????????/
		*/

		group, err = getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		} else if !group.Active {
			return myerr.NewClientError("The group is currently being deleted")
		}

		user, err = getUserWithConn(tx, username)
		if err != nil {
			return err
		}

		if group.OwnerID != currUserID && user.ID != currUserID {
			return myerr.NewClientError("Only the owner of the group can revoke membership of other members")
		} else if group.OwnerID == currUserID && user.ID == currUserID {
			return myerr.NewClientError("The owner cannot remove its own membership. Yet to be added this functionality")
		}

		log.Printf("Revolking membership for user with id [%d] in group with id [%d]", user.ID, group.ID)
		result := tx.Where("user_id = ?", user.ID).
			Where("group_id = ?", group.ID).
			Delete(&models.Membership{})

		if result.Error != nil {
			return myerr.NewServerErrorWrap(result.Error, "Problem with the creation of new membership in db")
		} else if result.RowsAffected == 0 {
			return myerr.NewClientError("Membership not found")
		}
		log.Printf("Membership for user with id [%d] in group id [%d] is revoked", user.ID, group.ID)

		return nil
	})
}

//MemberExists - check if membership exists for a particular group
func (d *UamDAOImpl) MemberExists(userID uint, groupID uint) (bool, error) {
	var count int64
	result := d.dbConn.Table("memberships").
		Where("user_id = ?", userID).
		Where("group_id = ?", groupID).
		Count(&count)

	if result.Error != nil {
		return false, myerr.NewServerErrorWrap(result.Error, "Problem with check existance of membership")
	}
	return count != 0, nil
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
