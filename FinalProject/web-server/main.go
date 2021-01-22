package main

import (
	"log"

	"example.com/user/web-server/api/rest"
	"example.com/user/web-server/internal/auth"
	"example.com/user/web-server/internal/db/dao"
	myerr "example.com/user/web-server/internal/error"
	"example.com/user/web-server/internal/middleware"
	val "example.com/user/web-server/internal/validator"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	var (
		uamDAO     dao.UamDAO
		jwtCreator auth.JwtCreator
		filter     middleware.AuthzFilter
	)

	uamDAO, err := dao.NewUamDAOImpl()
	uamDAO.Migrate()
	if err != nil {
		log.Fatal(myerr.NewServerErrorWrap(err, "Couldnt create a new User Access Management DAO"))
	}

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
			protected.POST("/group/membership/revocation", endpoint.RovokeMembership)
			protected.POST("/group/creation", endpoint.CreateGroup)
			protected.POST("/group/invitation", endpoint.AddMember)
			protected.DELETE("/user/deletion", endpoint.DeleteUser)
		}
	}
	log.Fatal(router.Run(":8080"))
}
