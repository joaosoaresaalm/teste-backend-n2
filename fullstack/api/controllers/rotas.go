package controllers

import "github.com/joaosoaresa/fullstack/api/middlewares"

func (s *Servidor) inicializarRotas() {

	// Rota Inicial
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")


	//Rotas de Usuários
	s.Router.HandleFunc("/clientes", middlewares.SetMiddlewareJSON(s.Criar)).Methods("POST")
	s.Router.HandleFunc("/clientes", middlewares.SetMiddlewareJSON(s.ObterTodos)).Methods("GET")
	s.Router.HandleFunc("/clientes/{id}", middlewares.SetMiddlewareJSON(s.Obter)).Methods("GET")
	s.Router.HandleFunc("/clientes/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.Atualizar))).Methods("PUT")
	s.Router.HandleFunc("/clientes/{id}", middlewares.SetMiddlewareAuthentication(s.Deletar)).Methods("DELETE")

	// Rota de Login
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	
}