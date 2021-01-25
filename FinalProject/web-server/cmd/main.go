package main

import (
	"log"

	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/api/rest"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/auth"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/dao"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/db/dbconn"
	myerr "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/error"
	"github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/middleware"
	val "github.com/danielpenchev98/FMI-Golang/FinalProject/web-server/internal/validator"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	var (
		uamDAO     dao.UamDAO
		jwtCreator auth.JwtCreator
		filter     middleware.AuthzFilter
	)

	dbConn, err := dbconn.GetDBConn(dbconn.PostgresDialectorCreator)
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a connection to the database"))
	}

	uamDAO = dao.NewUamDAOImpl(dbConn)
	if err = uamDAO.Migrate(); err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt migrate the database schemas"))
	}

	jwtCreator, err = auth.NewJwtCreatorImpl()
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a new Jwt Creator"))
	}

	filter = middleware.NewAuthzFilterImpl(jwtCreator)

	endpoint := rest.NewUamEndPointImpl(uamDAO, jwtCreator, val.NewBasicValidator())

	v1 := router.Group("/v1")
	{
		public := v1.Group("/public")
		{
			public.POST("/user/registration", endpoint.CreateUser)
			public.POST("/user/login", endpoint.Login)
		}

		protected := v1.Group("/protected").Use(filter.Authz)
		{
			protected.POST("/group/membership/revocation", endpoint.RevokeMembership)
			protected.POST("/group/creation", endpoint.CreateGroup)
			protected.POST("/group/invitation", endpoint.AddMember)
			protected.DELETE("/user/deletion", endpoint.DeleteUser)
		}
	}
	log.Fatal(router.Run(":8080"))
}
