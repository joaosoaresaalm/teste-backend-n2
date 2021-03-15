package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joaosoaresa/fullstack/api/models"
)

type Servidor struct {
	DB     *gorm.DB
	Router *mux.Router
}
// Inicializa o db utilizado
func (server *Servidor) Inicializar(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var erro error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, erro = gorm.Open(Dbdriver, DBURL)
		if erro != nil {
			fmt.Printf("Não é possível conectar %s db", Dbdriver)
			log.Fatal("Este é o erro:", erro)
		} else {
			fmt.Printf("Estamos conectados ao %s db", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.Usuario{}) //migração db

	server.Router = mux.NewRouter()

	server.inicializarRotas()
}

func (server *Servidor) Run(addr string) {
	fmt.Println("Ouvindo a porta 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}