package cron

import (
	"log"
	"os"
	"path"
	"sync"

	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/db/dao"
)

//GroupEraserJob - interface for group erase job
type GroupEraserJob interface {
	DeleteGroups()
}

//GroupEraserJobImpl - implementation of GroupEraserJob
type GroupEraserJobImpl struct {
	uamDAO    dao.UamDAO
	groupsDir string
}

//NewGroupEraserJobImpl - creates an instance of GroupEraserJobImpl
func NewGroupEraserJobImpl(uamDAO dao.UamDAO, groupsDir string) *GroupEraserJobImpl {
	return &GroupEraserJobImpl{
		uamDAO:    uamDAO,
		groupsDir: groupsDir,
	}
}

//DeleteGroups - deletes all deactivated job resources
func (i *GroupEraserJobImpl) DeleteGroups() {
	groupNames, err := i.uamDAO.GetDeactivatedGroupNames()
	if err != nil {
		log.Printf("Couldnt delete the resources of the groups in deleted state. Reason: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	for _, name := range groupNames {
		groupDir := path.Join(i.groupsDir, name)
		wg.Add(1)
		go func() {
			defer wg.Done()
			os.RemoveAll(groupDir)
		}()
	}

	wg.Wait()
	err = i.uamDAO.EraseDeactivatedGroups(groupNames)
	if err != nil {
		log.Printf("Couldnt erase inactive groups. Reason: %v\n", err)
	}
}
