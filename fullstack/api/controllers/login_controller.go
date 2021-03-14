package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/joaosoaresa/fullstack/api/auth"
	"github.com/joaosoaresa/fullstack/api/models"
	"github.com/joaosoaresa/fullstack/api/responses"
	"github.com/joaosoaresa/fullstack/api/utils/format-error"
	"golang.org/x/crypto/bcrypt"
)

func (server *Servidor) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.Usuario{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Preparar()
	err = user.Validar("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Senha)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Servidor) SignIn(email, password string) (string, error) {

	var err error

	user := models.Usuario{}

	err = server.DB.Debug().Model(models.Usuario{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerificarSenha(user.Senha, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CriarToken(user.ID)
}
