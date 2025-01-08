package main

import (
	"api-tarefas/configs"
	"api-tarefas/routes"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := configs.ConnectionDB()

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	api := gin.Default()

	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Métodos HTTP permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Headers permitidos
		ExposeHeaders:    []string{"Content-Length"},                          // Headers visíveis para o cliente
		AllowCredentials: true,                                                // Permitir cookies e credenciais
		MaxAge:           12 * time.Hour,
	}))

	// registrar rotas
	routes.TaskRoute(api, db)

	if err := api.Run(":3000"); err != nil {
		fmt.Println("Erro ao iniciar servidor")
	}
}
