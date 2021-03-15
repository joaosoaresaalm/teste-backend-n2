package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "nome") {
		return errors.New("Nome já existe")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email já existe")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Senha incorreta")
	}
	return errors.New("Dados incorretos")
}