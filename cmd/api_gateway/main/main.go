package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmrgreenteaa/user-management-service/config"
	_ "github.com/mmrgreenteaa/user-management-service/docs"
	api "github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			User Management service
// @version		1.0
// @description	API for user management and access
// @termsOfService	http://swagger.io/terms/
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	router := gin.Default()
	cfg, err := config.GetApiGatwayConfg()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.UserMengIP)

	mcs := api.NewMenrics()

	
	apgt, err := api.NewApiGatwaty(cfg)
	if err != nil {
		log.Fatal(err)
	}

	v1 := router.Group("/api/v1")
	{
		v1.Use(apgt.DoMetrics(mcs))

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
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	apgt.Srv = srv

	go func() {
		log.Println("Сервер запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска: %s", err)
		}
	}()
	<-quit
	GraceShoutDown(apgt)
}

func GraceShoutDown(apiGw *api.ApiGatway) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := apiGw.Srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	if apiGw.Rdb != nil {
		apiGw.Rdb.Close()
	}
	log.Println("Server SHOTDOWN")
	return nil
}
