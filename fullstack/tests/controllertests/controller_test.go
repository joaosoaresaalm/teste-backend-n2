package controllertests

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
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")


	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Não é possível conectar à %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {

	err := server.DB.DropTableIfExists(&models.Usuario{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Usuario{}).Error
	if err != nil {
		return err
	}

	log.Printf("Tabelas utilizadas com sucesso")
	return nil
}

func seedOneUser() (models.Usuario, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.Usuario{
		Nome: "Vitao",
		Email:    "vitao@gmail.com",
		Senha: "password",
	}

	err = server.DB.Model(&models.Usuario{}).Create(&user).Error
	if err != nil {
		return models.Usuario{}, err
	}
	return user, nil
}

func seedUsers() ([]models.Usuario, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.Usuario{
		models.Usuario{
			Nome: "Joao Soares",
			Email:    "joaosoaresa.alm@gmail.com",
			Senha: "password",
		},
		models.Usuario{
			Nome: "Tarcisio",
			Email:    "tarcisio@gmail.com",
			Senha: "password",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.Usuario{}).Create(&users[i]).Error
		if err != nil {
			return []models.Usuario{}, err
		}
	}
	return users, nil
}

