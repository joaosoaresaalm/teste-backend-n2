package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/joaosoaresa/fullstack/api/controllers"
	"github.com/joaosoaresa/fullstack/api/seed"
)

var server = controllers.Servidor{}

func Run() {

	var erro error
	erro = godotenv.Load()
	if erro != nil {
		log.Fatalf("Erro ao obter env,  %v", erro)
	} else {
		fmt.Println("Obtendo os valores do .env")
	}

	server.Inicializar(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")

}