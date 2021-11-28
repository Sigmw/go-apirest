package modelos

import (
	"api/src/config"
	"fmt"

	"github.com/gorilla/securecookie"
)

// DadosAutenticacao contém o token e o id do usuário autenticado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

var s *securecookie.SecureCookie

//configura o secureCookie
func Configurar() {
	s = securecookie.New(config.HashKey, config.BlockKey)
}

//pega id e token pelo token
func UserDataPeloToken(tk string) DadosAutenticacao {
	userdata := DadosAutenticacao{}
	fmt.Println(tk[1:])
	err := s.Decode("Authorization", tk[1:], &userdata)
	if err != nil {
		fmt.Println(err)
	}
	return userdata
}
