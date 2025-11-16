package main

import (
	_ "avito-tech-go-task/docs"
	"avito-tech-go-task/internal/application/service"
	"avito-tech-go-task/internal/clients/postgres"
	"avito-tech-go-task/internal/infrastructure/http/controller"
	"avito-tech-go-task/internal/infrastructure/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var (
	dsn string
)

func init() {
	dsn = os.Getenv("DSN")
}

func main() {
	r := gin.Default()

	db, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	prRepo := storage.NewPRRepo(db)
	teamRepo := storage.NewTeamRepo(db)
	userRepo := storage.NewUserRepo(db)

	prService := service.NewPRService(prRepo, userRepo, teamRepo)
	c := controller.NewApiService(prService)

	teams := r.Group("/teams")
	{
		teams.POST("add", c.AddTeamHandler)
		teams.GET("get", c.GetTeamHandler)
	}
	users := r.Group("/users")
	{
		users.POST("setIsActive", c.SetIsActiveUserHandler)
		users.GET("getReview", c.GetReviewerUserHandler)
		users.GET("getUserReviewStats", c.GetStatsHandler)
	}
	pullRequests := r.Group("/pullRequests")
	{
		pullRequests.POST("create", c.CreatePullRequestHandler)
		pullRequests.POST("merge", c.MergePullRequestHandler)
		pullRequests.POST("reassign", c.ReassignPullRequestHandler)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
