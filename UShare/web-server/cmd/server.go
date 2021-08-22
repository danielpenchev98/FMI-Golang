package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"net/http"
	"os/signal"
	"syscall"
	"context"
	"time"

	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/api/rest"
	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/auth"
	asyncJob "github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/cron"
	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/db/dao"
	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/db/dbconn"
	myerr "github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/error"
	"github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/middleware"
	val "github.com/danielpenchev98/FMI-Golang/UShare/web-server/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/pkg/errors"
)

var router = gin.Default()
var groupDirPath string

func main() {
	host := os.Getenv("HOST")
	if host == ""{
		log.Fatal("Please set HOST")
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("Please set PORT env variable")
	}

	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("The env variable PORT has illegal port number")
	}

	if err = createGroupsDir(); err != nil {
		log.Fatal(err)
	}

	jwtCreator, err := auth.NewJwtCreatorImpl()
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a new Jwt Creator"))
	}

	filter := middleware.NewAuthzFilterImpl(jwtCreator)

	uamEndpoint := rest.NewUamEndPointImpl(createUamDAO(), jwtCreator, val.NewBasicValidator(), groupDirPath)
	fmEndpoint := rest.NewFileManagementEndpointImpl(createUamDAO(), createFmDAO(), groupDirPath)

	groupDeleter := asyncJob.NewGroupEraserJobImpl(createUamDAO(), groupDirPath)
	asyncJob := cron.New()
	asyncJob.AddFunc("@every 1m", groupDeleter.DeleteGroups)
	asyncJob.Start()
	defer asyncJob.Stop()

	v1 := router.Group("/v1")
	{
		public := v1.Group("/public")
		{
			public.POST("/user/registration", uamEndpoint.CreateUser)
			public.POST("/user/login", uamEndpoint.Login)
		}

		protected := v1.Group("/protected").Use(filter.Authz)
		{
			protected.DELETE("/group/membership/revocation", uamEndpoint.RevokeMembership)
			protected.POST("/group/creation", uamEndpoint.CreateGroup)
			protected.POST("/group/invitation", uamEndpoint.AddMember)
			protected.DELETE("/user/deletion", uamEndpoint.DeleteUser)
			protected.DELETE("/group/deletion", uamEndpoint.DeleteGroup)
			protected.POST("/group/file/upload", fmEndpoint.UploadFile)
			protected.GET("/group/file/download", fmEndpoint.DownloadFile)
			protected.DELETE("/group/file/deletion", fmEndpoint.DeleteFile)
			protected.GET("/group/files", fmEndpoint.RetrieveAllFilesInfo)
			protected.GET("/groups", uamEndpoint.GetAllGroupsInfo)
			protected.GET("/users", uamEndpoint.GetAllUsersInfo)
			protected.GET("/group/users", uamEndpoint.GetAllUsersInGroup)
		}
	}

	httpServer := &http.Server{
        Addr:    fmt.Sprintf("%s:%d", host, portNum),
        Handler: router,
    }
 
    go func() {
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            panic(errors.Wrapf(err, "server listen-and-serve failed"))
        }
    }()
 
    done := make(chan os.Signal, 1)
    signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
    <-done
 
    //log.Info("shutting down http server...")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
 
    if err := httpServer.Shutdown(ctx); err != nil {
        panic(errors.Wrapf(err, "failed to shutdown server"))
    }
    // nolint:gosimple
    select {
    case <-ctx.Done():
    }
}

func createGroupsDir() error {
	currDir := os.Getenv("GROUP_DIR")
	if currDir == "" {
		return errors.New("Please set GROUP_DIR is not set")
	}

	groupDirPath = currDir + "/groups"
	if _, err := os.Stat(groupDirPath); err == nil {
		return nil
	}

	if err := os.Mkdir(groupDirPath, 0755); err != nil {
		return myerr.NewServerErrorWrap(err, "Couldnt create directory for the groups")
	}

	return nil

}

func createUamDAO() dao.UamDAO {
	dbConn, err := dbconn.GetDBConn(dbconn.PostgresDialectorCreator)
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a connection to the database"))
	}

	uamDAO := dao.NewUamDAOImpl(dbConn)
	if err = uamDAO.Migrate(); err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt migrate the database schemas"))
	}

	return uamDAO
}

func createFmDAO() dao.FmDAO {
	dbConn, err := dbconn.GetDBConn(dbconn.PostgresDialectorCreator)
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a connection to the database"))
	}

	fmDAO := dao.NewFmDAOImpl(dbConn)
	if err = fmDAO.Migrate(); err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt migrate the database schemas"))
	}

	return fmDAO
}
