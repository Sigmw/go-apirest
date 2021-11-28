package modelos

import (
	"api/src/banco"
	"api/src/seguranca"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário utilizando a rede social
type Usuario struct {
	ID         uint64    `json:"id,omitempty"`
	Nome       string    `json:"nome,omitempty"`
	Nick       string    `json:"nick,omitempty"`
	Email      string    `json:"email,omitempty"`
	Senha      string    `json:"senha,omitempty"`
	CriadoEm   time.Time `json:"CriadoEm,omitempty"`
	Desativado bool      //é meio obvio
	Codigo     string    //codigo para captcha
	Cargo      int       //0 = usuario 1 = admin 2 = admin master
}

// Preparar vai chamar os métodos para validar e formatar o usuário recebido
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}

	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}

	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if usuario.Nick == "" {
		return errors.New("O nick é obrigatório e não pode estar em branco")
	}

	if usuario.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("O e-mail inserido é inválido")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("A senha é obrigatória e não pode estar em branco")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash)
	}

	return nil
}

//se vc nao entendeu sugiro consultar um medico
func PegarUsuarioPeloID(id uint64) (Usuario, error) {
	db, err := banco.Conectar()
	if err != nil {
		return Usuario{}, err
	}

	query := `select nome, nick, email, cargo, desativado from usuarios where id = ?`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return Usuario{}, err
	}
	defer stmt.Close()

	user := Usuario{}

	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&user.Nome, &user.Nick, &user.Email, &user.Cargo, &user.Desativado); err != nil {
		return Usuario{}, err
	}
	return user, nil
}
