package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/api/common"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/api/common/response"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/dao"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/models"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
	"github.com/gin-gonic/gin"
)

type FileManagementEndpoint interface {
	UploadFile(*gin.Context)
	DownloadFile(*gin.Context)
	DeleteFile(*gin.Context)
}

type FileManagementEndpointImpl struct {
	UamDAO    dao.UamDAO
	groupsDir string
	FmDAO     dao.FmDAO
}

func NewFileManagementEndpointImpl(uam dao.UamDAO, fm dao.FmDAO, groupsDir string) *FileManagementEndpointImpl {
	return &FileManagementEndpointImpl{
		UamDAO:    uam,
		FmDAO:     fm,
		groupsDir: groupsDir,
	}
}

func (f *FileManagementEndpointImpl) UploadFile(c *gin.Context) {
	var (
		userID uint
		err    error
	)

	if userID, err = common.GetIDFromContext(c); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Problem with the file"))
		return
	}

	var group models.Group
	groupName := c.Query("groupname")
	if group, err = f.UamDAO.GetGroup(groupName); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	// could be replaced with getMembership???
	if exists, err := f.UamDAO.MemberExists(userID, group.ID); err != nil {
		common.SendErrorResponse(c, err)
		return
	} else if exists != true {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user input"))
		return
	}

	fileID, err := f.FmDAO.AddFileInfo(userID, file.Filename, groupName)
	if err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	dst := fmt.Sprintf("%s/%s/%d", f.groupsDir, groupName, fileID)
	if err = c.SaveUploadedFile(file, dst); err != nil {
		f.FmDAO.RemoveFileInfo(userID, fileID, groupName)
		common.SendErrorResponse(c, myerr.NewServerError(fmt.Sprintf("Couldnt save the file in the group dir [%s]", groupName)))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"file_id": fileID,
	})
}

type DownloadFileRequest struct {
	GroupName string `json:"group_name"`
	FileID    uint   `json:"file_id"`
}

//DownloadFile - downloads a file given group
func (f *FileManagementEndpointImpl) DownloadFile(c *gin.Context) {
	var (
		userID uint
		err    error
	)

	if userID, err = common.GetIDFromContext(c); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	var rq DownloadFileRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	group, err := f.UamDAO.GetGroup(rq.GroupName)
	if exists, err := f.UamDAO.MemberExists(userID, group.ID); err != nil {
		common.SendErrorResponse(c, err)
		return
	} else if exists != true {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user input"))
		return
	}

	fileInfo, err := f.FmDAO.GetFileInfo(userID, rq.FileID, rq.GroupName)
	if err != nil {
		common.SendErrorResponse(c, err)
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name))
	filePath := fmt.Sprintf("%s/%s/%d", f.groupsDir, rq.GroupName, fileInfo.ID)
	c.File(filePath)
}

type FileRequest struct {
	GroupName string `json:"group_name"`
	FileID    uint   `json:"file_id"`
}

func (f *FileManagementEndpointImpl) DeleteFile(c *gin.Context) {
	var (
		userID uint
		err    error
	)

	if userID, err = common.GetIDFromContext(c); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	var rq FileRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid json body"))
		return
	}

	group, err := f.UamDAO.GetGroup(rq.GroupName)
	if exists, err := f.UamDAO.MemberExists(userID, group.ID); err != nil {
		common.SendErrorResponse(c, err)
		return
	} else if exists != true {
		common.SendErrorResponse(c, myerr.NewClientError("Invalid user input"))
		return
	}

	if err := f.FmDAO.RemoveFileInfo(userID, rq.FileID, rq.GroupName); err != nil {
		common.SendErrorResponse(c, err)
		return
	}

	path := fmt.Sprintf("%s/%s/%d", f.groupsDir, rq.GroupName, rq.FileID)
	os.Remove(path)

	c.JSON(http.StatusOK, response.BasicResponse{
		Status: http.StatusOK,
	})
}
