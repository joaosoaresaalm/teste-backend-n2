
package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/joaosoaresa/fullstack/api/controllers"
	"github.com/joaosoaresa/fullstack/api/models"
)

var server = controllers.Servidor{}
var userInstance = models.Usuario{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	log.Printf("--------Antes da chamada api.Run()--------")
	ret := m.Run()
	log.Printf("--------Depois da chamadaapi.Run()---------")
	//os.Exit(m.Run())
	os.Exit(ret)
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")
	
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Não é possível conectar %s db\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Conexão com  %s db\n", TestDbDriver)
		}
	}

}

func refreshUserTable() error {
	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.Usuario{}).Error
	if err != nil {
		return err
	}
	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.Usuario{}).Error
	if err != nil {
		return err
	}
	log.Printf("Tabela atualizada com suceso!")
	log.Printf("Rotina de método ok.")
	return nil
}

func seedUser() (models.Usuario, error) {

	_ = refreshUserTable()

	user := models.Usuario{
		Nome: "Pet",
		Email:    "pet@gmail.com",
		Senha: "password",
	}

	err := server.DB.Debug().Model(&models.Usuario{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Não é possível popular a Tabela de Usuários: %v", err)
	}

	log.Printf("Rotina de método ok.")
	return user, nil
}

func seedUsers() error {

	users := []models.Usuario{
		models.Usuario{
			Nome: "João Soares",
			Email:    "joaosoaresa.alm@gmail.com",
			Senha: "password",
		},
		models.Usuario{
			Nome: "Tarcisio do Acordeon",
			Email:    "tarcisio.acordeon@gmail.com",
			Senha: "password",
		},
	}

	for i := range users {
		err := server.DB.Debug().Model(&models.Usuario{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}

	log.Printf("Rotina de método ok.")
	return nil
}