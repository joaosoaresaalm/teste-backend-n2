package controllers

import (
	"net/http"
	"github.com/joaosoaresa/fullstack/api/responses"
)

func (server *Servidor) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Bem-vindo à API de Clientes")

}