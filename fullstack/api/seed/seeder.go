package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/joaosoaresa/fullstack/api/models"
)

var usuarios = []models.Usuario{
	models.Usuario{
		Nome: "Joao Soares",
		Email:    "joaosoaresa.alm@gmail.com",
		Senha: "senha",
	},
	models.Usuario{
		Nome: "Tarcisio do Acordeon",
		Email:    "tarcisio.acordeon@gmail.com",
		Senha: "senha",
	},
}

func Load(db *gorm.DB) {

	erro := db.Debug().DropTableIfExists(&models.Usuario{}).Error
	if erro != nil {
		log.Fatalf("não é possível fazer drop: %v", erro)
	}
	erro = db.Debug().AutoMigrate(&models.Usuario{}).Error
	if erro != nil {
		log.Fatalf("não é possível realizar migração: %v", erro)
	}

	for i, _ := range usuarios {
		erro = db.Debug().Model(&models.Usuario{}).Create(&usuarios[i]).Error
		if erro != nil {
			log.Fatalf("não é possível popular a tabela de usuários: %v", erro)
		}
	}
}