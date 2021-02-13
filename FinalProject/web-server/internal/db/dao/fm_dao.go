package dao

import (
	"errors"
	"fmt"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/models"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
	"gorm.io/gorm"
)

//go:generate mockgen --source=fm_dao.go --destination dao_mocks/fm_dao.go --package dao_mocks

//FmDAO - interface, used for file management
type FmDAO interface {
	AddFileInfo(userID uint, fileName string, groupName string) (uint, error)
	GetFileInfo(userID uint, fileID uint, groupName string) (models.FileInfo, error)
	GetAllFilesInfo(userID uint, groupName string) ([]models.FileInfo, error)
	RemoveFileInfo(userID uint, fileID uint, groupName string) error
	Migrate() error
}

//FmDAOImpl - implementation of FmDAO
type FmDAOImpl struct {
	dbConn *gorm.DB
}

//NewFmDAOImpl - creates an instance of FmDAOImpl
func NewFmDAOImpl(dbConn *gorm.DB) *FmDAOImpl {
	return &FmDAOImpl{
		dbConn: dbConn,
	}
}

//Migrate - updates the models in the db
func (f *FmDAOImpl) Migrate() error {
	return f.dbConn.AutoMigrate(models.FileInfo{})
}

//AddFileInfo - saves metadate for a newly added file (just like in linux with inodes)
func (f *FmDAOImpl) AddFileInfo(userID uint, fileName string, groupName string) (uint, error) {
	var (
		fileID uint
		err    error
	)
	err = f.dbConn.Transaction(func(tx *gorm.DB) error {

		group, err := getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		}

		var count int64
		result := tx.Table("memberships").
			Where("user_id = ?", userID).
			Where("group_id = ?", group.ID).
			Count(&count)

		if result.Error != nil {
			return myerr.NewServerError(fmt.Sprintf("Couldnt check if membership exists. Reason: %v\n", result.Error))
		} else if count == 0 {
			return myerr.NewClientError("Cannot upload a file in a group you aren't part of")
		}

		fileInfo := models.FileInfo{
			Name:    fileName,
			OwnerID: userID,
			GroupID: group.ID,
		}

		if result = tx.Create(&fileInfo); result.Error != nil {
			return myerr.NewServerError(fmt.Sprintf("Cannot save file info in the db for group [%s]", groupName))
		}
		fileID = fileInfo.ID
		return nil
	})
	return fileID, err
}

//RemoveFileInfo - removes the file matadata from the db
func (f *FmDAOImpl) RemoveFileInfo(userID uint, fileID uint, groupName string) error {
	return f.dbConn.Transaction(func(tx *gorm.DB) error {
		group, err := getGroupWithConn(tx, groupName)
		if err != nil {
			return err
		}

		fileInfo, err := getFileInfoWithConn(tx, fileID)
		if err != nil {
			return err
		}

		if group.OwnerID != userID && fileInfo.OwnerID != userID {
			return myerr.NewClientError("Only the onwer of the file or the group owner can remove files from the group")
		}

		if result := tx.Delete(&fileInfo); result.Error != nil {
			return myerr.NewServerError(fmt.Sprintf("Cannot save file info in the db for group [%s]", groupName))
		} else if result.RowsAffected == 0 {
			return myerr.NewClientError("File info not found")
		}
		fileID = fileInfo.ID
		return nil
	})
}

//GetFileInfo - fetches metadata for a particular file
func (f *FmDAOImpl) GetFileInfo(userID uint, fileID uint, groupName string) (models.FileInfo, error) {

	var count int64
	result := f.dbConn.Table("memberships").Joins("inner join groups on memberships.group_id = groups.id").
		Where("groups.name = ?", groupName).
		Where("memberships.user_id = ?", userID).
		Count(&count)

	if result.Error != nil {
		return models.FileInfo{}, myerr.NewServerErrorWrap(result.Error, "Problem with checking if user is a member of the group.")
	} else if count == 0 {
		return models.FileInfo{}, myerr.NewClientError("You arent a member of the group.")
	}

	fileInfo, err := getFileInfoWithConn(f.dbConn, fileID)
	if err != nil {
		return models.FileInfo{}, err
	}
	return fileInfo, err

}

func (i *FmDAOImpl) GetAllFilesInfo(userID uint, groupName string) ([]models.FileInfo, error) {
	var count int64
	result := i.dbConn.Table("memberships").Joins("inner join groups on memberships.group_id = groups.id").
		Where("groups.name = ?", groupName).
		Where("memberships.user_id = ?", userID).
		Count(&count)

	if result.Error != nil {
		return nil, myerr.NewServerErrorWrap(result.Error, "Problem with checking if user is a member of the group.")
	} else if count == 0 {
		return nil, myerr.NewClientError("You arent a member of the group.")
	}

	var fileInfos []models.FileInfo
	result = i.dbConn.Table("file_infos").Joins("inner join groups on file_infos.group_id = groups.id").Take(&fileInfos)
	if result.Error != nil {
		return nil, myerr.NewServerErrorWrap(result.Error, "Problem with fetching all files from a specific group")
	}

	return fileInfos, nil
}

func getFileInfoWithConn(dbConn *gorm.DB, fileID uint) (models.FileInfo, error) {
	var fileInfo models.FileInfo

	result := dbConn.Table("file_infos").
		Where("id = ?", fileID).
		Take(&fileInfo)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fileInfo, myerr.NewItemNotFoundError("File does not exist")
	} else if result.Error != nil {
		return fileInfo, myerr.NewServerErrorWrap(result.Error, "Problem with the lookup if file exists")
	}

	return fileInfo, nil
}
