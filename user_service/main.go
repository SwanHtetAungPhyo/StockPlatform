package main

import (
	"os"
	"time"

	"github.com/SwanHtetAungPhyo/closure/closure"
	"github.com/SwanHtetAungPhyo/user-service/database"
	"github.com/SwanHtetAungPhyo/user-service/handler"
	"github.com/SwanHtetAungPhyo/user-service/models"
	"github.com/SwanHtetAungPhyo/user-service/repository"
	service "github.com/SwanHtetAungPhyo/user-service/services"
)


func main() {
	database.DB_INIT()
	database.Migration(&models.User{}, &models.Balance{})
	app := closure.New()

	repo := repository.NewRepository(database.DB)

	service.SetRepository(repo, os.Getenv("SECRET_KEY"),24 *time.Hour)
	app.Cluster("/api/v1", func(apiCluster *closure.Cluster) {
		apiCluster.
			Post("/signup",handler.SignUpHandler).
			Post("/signin", handler.SignInHandler)
	})

	app.Cluster("/api/balance/",func(balanceCluster  *closure.Cluster) {
		balanceCluster.Post("", handler.Deposit)
		balanceCluster.Get("", handler.Read)
	})
	app.Start(":3001")
}