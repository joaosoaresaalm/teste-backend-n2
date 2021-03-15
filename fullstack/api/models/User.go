package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"
	"os"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nome      string    `gorm:"size:255;not null;unique" json:"nome"`
	Endereco  string    `gorm:"size:255;not null;" json:endereco`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Senha     string    `gorm:"size:100;not null;" json:"senha"`
	CriadoEm  string     `gorm:json:"criado_em"`
	AtualizadoEm string  `gorm:json:"atualizado_em"`
}

func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

func VerificarSenha(hashedSenha, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedSenha), []byte(senha))
}

func (u *Usuario) salvarAntes() error {
	hashedSenha, err := Hash(u.Senha)
	if err != nil {
		return err
	}
	u.Senha = string(hashedSenha)
	return nil
}

func (u *Usuario) Preparar() {
	u.ID = 0;
	u.Nome = html.EscapeString(strings.TrimSpace(u.Nome))
	u.Endereco = html.EscapeString(strings.TrimSpace(u.Endereco))
	u.CriadoEm = manipularData()
	u.AtualizadoEm = manipularData()

}

func manipularData() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}


func (u *Usuario) Validar(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nome == "" {
			return errors.New("Nome Obrigatório")
		}
		if u.Endereco == "" {
			return errors.New("Endereço Obrigatório")
		}
		if u.Senha == "" {
			return errors.New("Senha Obrigatória")
		}
		return nil
	case "login":
		if u.Senha == "" {
			return errors.New("Nome Obrigatório")
		}
		return nil

	default:
		if u.Nome == "" {
			return errors.New("Nome Obrigatório")
		}
		if u.Endereco == "" {
			return errors.New("Endereço obrigatório")
		}
		if u.Senha == "" {
			return errors.New("Senha Obrigatória")
		}
		return nil
	}
}

func (u *Usuario) Salvar(db *gorm.DB) (*Usuario, error) {
	var err error

	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	
	criarDiretorio()
	escreverJSON(u)
	
	return u, nil
}

// Função com estratégia de montar o arquivo .json
// Por exemplo, joaosoaresa.json
func escreverJSON(u *Usuario) {
	const PATH = "NOVOS_CLIENTES/"
	var NOME_ARQUIVO = u.Nome + ".json"
	var arquivo string = indentarJSON(u)

	nomeArquivo, erro := os.Create(PATH + NOME_ARQUIVO)
	escreverArquivo, err := nomeArquivo.WriteString(string(arquivo))
	fmt.Println(escreverArquivo)
	
	checarErro(erro)

	defer nomeArquivo.Close()
	if err != nil {
		log.Fatal(err)
	}

}

// Função com estratégia de obter o objeto da requisição e formatá-lo
func indentarJSON(u *Usuario) string {
	var usuarioFormatado Usuario

	usuarioFormatado = Usuario{
		ID: u.ID,
		Nome: u.Nome,
		Endereco: u.Endereco,
		Email: u.Email,
		Senha: u.Senha,
		CriadoEm: u.CriadoEm,
		AtualizadoEm: u.AtualizadoEm,
	}
	arquivo, _ := json.MarshalIndent(usuarioFormatado, "", "")
	return string(arquivo)
}


func checarErro(erro error) {
	if erro != nil {
		log.Fatal(erro)
	}
}

// Funcão para criar DIR raiz, caso não exista
func criarDiretorio() error {
	const PATH = "NOVOS_CLIENTES"
	_, err := os.Stat(PATH)
	if os.IsNotExist(err){
		err := os.Mkdir(PATH, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil 
}

func (u *Usuario) ObterTodos(db *gorm.DB) (*[]Usuario, error) {
	var err error
	users := []Usuario{}
	err = db.Debug().Model(&Usuario{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]Usuario{}, err
	}
	return &users, err
}

func (u *Usuario) ObterPorId(db *gorm.DB, uid uint32) (*Usuario, error) {
	var err error
	err = db.Debug().Model(Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Usuario{}, errors.New("Usuário não encontrado")
	}
	return u, err
}

func (u *Usuario) Atualizar(db *gorm.DB, uid uint32) (*Usuario, error) {

	// Prepara um hash para a senha
	err := u.salvarAntes()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&Usuario{}).UpdateColumns(
		map[string]interface{}{
			"senha":  u.Senha,
			"nome":  u.Nome,
			"endereco": u.Endereco,
			"criadoEm": time.Now(),
		},
	)
	if db.Error != nil {
		return &Usuario{}, db.Error
	}
	// Esta é a tela do usuário atualizada
	err = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	return u, nil
}

func (u *Usuario) Deletar(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&Usuario{}).Delete(&Usuario{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}