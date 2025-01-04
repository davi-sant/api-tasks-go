package routes

import (
	"api-tarefas/controllers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func TaskRoute(r *gin.Engine, db *sql.DB) {
	tc := controllers.TaskController{DB: db}
	tr := r.Group("v1/tarefas", gin.BasicAuth(gin.Accounts{
		"admin": "36132294",
	}))
	{
		tr.GET("/", tc.GetTasks)
		tr.GET("/:id", tc.GetTask)
		tr.GET("/busca", tc.FindTasks)
		tr.POST("/", tc.CreateTask)
		tr.DELETE("/:id", tc.DeleteTask)
		tr.PUT("/:id", tc.UpdateTask)
	}

}
