package routes

import (
	"api-tarefas/controllers"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TaskRoute(api *gin.Engine, db *sql.DB) {
	controller := controllers.TaskController{DB: db}
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		return
	}
	admin := os.Getenv("ADMIN")
	password := os.Getenv("PASSWORD")
	tr := api.Group("v1/tarefas", gin.BasicAuth(gin.Accounts{admin: password}))
	{
		tr.GET("/", controller.GetTasks)
		tr.GET("/:id", controller.GetTask)
		tr.GET("/busca", controller.FindTasks)
		tr.POST("/", controller.CreateTask)
		tr.DELETE("/:id", controller.DeleteTask)
		tr.PUT("/:id", controller.UpdateTask)
	}

}
