package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/joaosoaresa/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestObterTodos(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatalf("Error seeding user table %v\n", err)
	}

	users, err := userInstance.ObterTodos(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestSalvar(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error user refreshing table %v\n", err)
	}
	newUser := models.Usuario{
		ID:       1,
		Email:    "test@gmail.com",
		Nome: "test",
		Senha: "password",
	}
	savedUser, err := newUser.Salvar(server.DB)
	if err != nil {
		t.Errorf("Error while saving a user: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nome, savedUser.Nome)
}

func TestObterPorId(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error user refreshing table %v\n", err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := userInstance.ObterPorId(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nome, user.Nome)
}

func TestAtualizar(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Não é possivel popular um usuario: %v\n", err)
	}

	userUpdate := models.Usuario{
		ID:       1,
		Nome: "modiUpdate",
		Email:    "modiupdate@gmail.com",
		Senha: "password",
	}
	updatedUser, err := userUpdate.Atualizar(server.DB, user.ID)
	if err != nil {
		t.Errorf("Este é o erro após tentar atualizar um usuário: %v\n", err)
		return
	}
	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nome, userUpdate.Nome)
}

func TestDeletar(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()

	if err != nil {
		log.Fatalf("Não é possivel popular usuário: %v\n", err)
	}

	isDeleted, err := userInstance.Deletar(server.DB, user.ID)
	if err != nil {
		t.Errorf("Este é o erro ao excluir o usuário: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
