package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mmrgreenteaa/user-management-service/config"
	api "github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
)

func main() {
	router := gin.Default()
	cfg, err := config.GetApiGatwayConfg()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.UserMengIP)
	apgt, err := api.NewApiGatwaty(cfg)
	if err != nil {
		log.Fatal(err)
	}
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", apgt.LogIn)
			auth.GET("/logout", apgt.LogOut)
			auth.GET("/refresh", apgt.RefreshToken)
		}
		userM := v1.Group("/user-manager")
		{
			userM.POST("/users", apgt.RegistrateUser)
			userM.PATCH("/users", apgt.AuthMiddleware(), apgt.EditUserLogin)
			userM.DELETE("/users", apgt.AuthMiddleware(), apgt.DeleteAccount)
		}
	}
	router.Run(":8080") 
	log.Println(router.Routes())
}
