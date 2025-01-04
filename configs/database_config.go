package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectionDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco: %v", err)
		return nil, err
	}

	// Testar a conexão
	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
		return nil, err
	}

	log.Println("Conexão com o banco de dados bem-sucedida!")
	return db, nil
}
