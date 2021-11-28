package seguranca

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

// Hash recebe uma string e coloca um hash nela
func Hash(senha string) ([]byte, error) {
	h := sha512.New()

	h.Write([]byte(senha))
	password_hashed := hex.EncodeToString(h.Sum(nil))
	return []byte(password_hashed), nil
}

// VerificarSenha compara uma senha e um hash e retorna se elas são iguais
func VerificarSenha(senhaComHash, senhaString string) error {
	a, _ := Hash(senhaString)

	if string(a) != senhaComHash {

		return errors.New("invalid password")
	}
	return nil
}
