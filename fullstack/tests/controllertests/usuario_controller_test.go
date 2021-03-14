package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joaosoaresa/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"nickname":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			nickname:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"nickname":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Nickname Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"nickname": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Nickname",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Criar)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestObter(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/usuarios", nil)
	if err != nil {
		t.Errorf("Este é o erro: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ObterTodos)
	handler.ServeHTTP(rr, req)

	var users []models.Usuario
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Não é possível converter para JSON: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			nickname:   user.Nome,
			email:      user.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/usuarios", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Obter)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Nome, responseMap["nome"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestAtualizar(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() // Precisamos de pelo menos dois usuários para verificar adequadamente a atualização
	if err != nil {
		log.Fatalf("Erro ao popular usuário: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" // Observe que a senha no banco de dados já está com hash, queremos sem hash
	}
	// Faça o login do usuário e obtenha o token de autenticação
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Não é possível realizar login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		atualizarJSON     string
		statusCodigo   int
		atualizarNome string
		atualizarEmail    string
		tokenGiven     string
		erroMessagem   string
	}{
		{
			// Converta int32 em int primeiro antes de converter em string
			id:             strconv.Itoa(int(AuthID)),
			atualizarJSON:     `{"nickname":"Grande", "email": "grande@gmail.com", "password": "password"}`,
			statusCodigo:     200,
			atualizarNome: "Grande",
			atualizarEmail:    "grande/@gmail.com",
			tokenGiven:     tokenString,
			erroMessagem:   "",
		},
		{
			// Quando o campo de senha está vazio
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"johnnyboy", "email": "johnnyboy@gmail.com", "password": ""}`,
			statusCodigo:   422,
			tokenGiven:   tokenString,
			erroMessagem: "Senha obrigatória",
		},
		{
			// Quando nenhum token foi passado
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"joe", "email": "joe@gmail.com", "password": "password"}`,
			statusCodigo:   401,
			tokenGiven:   "",
			erroMessagem: "Não autorizado",
		},
		{
			// Quando o token incorreto foi passado
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"johnnyboy", "email": "johnnyboy@gmail.com", "password": "password"}`,
			statusCodigo:   401,
			tokenGiven:   "Token informado incorreto",
			erroMessagem: "Não autorizado",
		},
		{
			// Lembre-se de que joao@gmail.com" pertence ao usuário 2
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"joao", "email": "joao@gmail.com", "password": "password"}`,
			statusCodigo:   500,
			tokenGiven:   tokenString,
			erroMessagem: "Email com este token já existe",
		},
		{
			// Lembre-se de que "joao" pertence ao usuário 2
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"João Soares", "email": "joaosoaresa.alm@gmail.com", "password": "password"}`,
			statusCodigo:   500,
			tokenGiven:   tokenString,
			erroMessagem: "Nome já escolhido ",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname":"Victor", "email": "victorgmail.com", "password": "password"}`,
			statusCodigo:   422,
			tokenGiven:   tokenString,
			erroMessagem: "Email inválido",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname": "", "email": "vitin@gmail.com", "password": "password"}`,
			statusCodigo:   422,
			tokenGiven:   tokenString,
			erroMessagem: "Nome obrigatório",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			atualizarJSON:   `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCodigo:   422,
			tokenGiven:   tokenString,
			erroMessagem: "Email obrigatório",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCodigo: 400,
		},
		{
			//Quando o usuário 2 está usando o token do usuário 1
			id:           strconv.Itoa(int(2)),
			atualizarJSON:   `{"nickname": "Victor", "email": "victor@gmail.com", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCodigo:   401,
			erroMessagem: "Não autorizado",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(v.atualizarJSON))
		if err != nil {
			t.Errorf("Este é o erro: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Atualizar)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Não é possível converter em JSON: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCodigo)
		if v.statusCodigo == 200 {
			assert.Equal(t, responseMap["nome"], v.atualizarNome)
			assert.Equal(t, responseMap["email"], v.atualizarEmail)
		}
		if v.statusCodigo == 401 || v.statusCodigo == 422 || v.statusCodigo == 500 && v.erroMessagem != "" {
			assert.Equal(t, responseMap["error"], v.erroMessagem)
		}
	}
}

func TestDeletar(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	usuarios, err := seedUsers() // Precisamos de pelo menos dois usuários para verificar adequadamente a atualização
	if err != nil {
		log.Fatalf("Erro ao propagar usuário: %v\n", err)
	}
	// Pegue apenas o primeiro e faça o login
	for _, usuario := range usuarios {
		if usuario.ID == 2 {
			continue
		}
		AuthID = usuario.ID
		AuthEmail = usuario.Email
		AuthPassword = "password" // Observe que a senha no banco de dados já está com hash, queremos sem hash
	}
	//	Faça o login do usuário e obtenha o token de autenticação
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("Não é possível fazer login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	usuarioModel := []struct {
		id           string
		tokenGiven   string
		statusCodigo   int
		erroMessagem string
	}{
		{
			// Converta int32 em int primeiro antes de converter em string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCodigo:   204,
			erroMessagem: "",
		},
		{
			// Quando nenhum token é fornecido
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCodigo:   401,
			erroMessagem: "Não autorizado",
		},
		{
			// Quando um token incorreto é fornecido
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "Este é um token inválido",
			statusCodigo:   401,
			erroMessagem: "Não autorizado",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCodigo: 400,
		},
		{
			// O usuário 2 está tentando usar o token do usuário 1
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCodigo:   401,
			erroMessagem: "Não autorizado",
		},
	}
	for _, v := range usuarioModel {

		req, err := http.NewRequest("GET", "/usuarios", nil)
		if err != nil {
			t.Errorf("Este é o erro: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Deletar)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCodigo)

		if v.statusCodigo == 401 && v.erroMessagem != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Não é possível converter para json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.erroMessagem)
		}
	}
}
