package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestLogarDB(t *testing.T) {

	err := atualizarTabelaUsuario()
	if err != nil {
		log.Fatal(err)
	}
	user, err := popularUmUsuario()
	if err != nil {
		fmt.Printf("Este é o erro %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", //Observe que a senha deve ser esta, não aquela com hash do banco de dados
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Email errado",
			errorMessage: "crypto/bcrypt: hashedPassword não é o hash da senha fornecida",
		},
		{
			email:        "Email errado",
			password:     "password",
			errorMessage: "Recurso não encontrado",
		},
	}

	for _, v := range samples {

		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {

	atualizarTabelaUsuario()

	_, err := popularUmUsuario()
	if err != nil {
		fmt.Printf("Este é o erro %v\n", err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "vitao@gmail.com", "senha": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "vitao@gmail.com", "senha": "atualizarTabelaUsuario"}`,
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			inputJSON:    `{"email": "jo@gmail.com", "senha": "password"}`,
			statusCode:   422,
			errorMessage: "Infos incorretas",
		},
		{
			inputJSON:    `{"email": "kangmail.com", "senha": "Akfpt6sg"}`,
			statusCode:   422,
			errorMessage: "Email inválido",
		},
		{
			inputJSON:    `{"email": "", "senha": "Akfpt6sg"}`,
			statusCode:   422,
			errorMessage: "Email obrigatório",
		},
		{
			inputJSON:    `{"email": "joaosoaresa.alm@gmail.com", "senha": ""}`,
			statusCode:   422,
			errorMessage: "Senha obrigatória",
		},
		{
			inputJSON:    `{"email": "", "senha": "password"}`,
			statusCode:   422,
			errorMessage: "Email obrigatório",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Este é o erro: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Não é possível converter em JSON: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
