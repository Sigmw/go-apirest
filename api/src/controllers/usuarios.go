package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func AlistarUsuario(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	usuario, err := modelos.PegarUsuarioPeloID(usuarioID)

	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	type Alistamento struct {
		Nick string
	}

	if usuario.Cargo < 1 {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Permissão negada."))
		return
	}

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	usuario_a_ser_alistado := Alistamento{}
	if erro = json.Unmarshal(corpoRequest, &usuario_a_ser_alistado); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		fmt.Println("fudeu bixo banco caiu")
		return
	}

	if usuario_a_ser_alistado.Nick == usuario.Nick {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Voce nao pode alistar a vc mesmo"))
		return
	}

	query := `SELECT COUNT(nick) AS users FROM alistados WHERE nick = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, err)
			fmt.Println("fudeu bixo deu merda no banco")
			return
		}
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, usuario_a_ser_alistado.Nick)
	number_of_users := 0
	if err := row.Scan(&number_of_users); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		fmt.Println("fudeu bixo deu merda no banco")
		return
	}

	if number_of_users > 0 {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Usuario ja alistado"))
		return
	}

	query = fmt.Sprintf("insert into alistados(nick) values (\"%v\")", usuario_a_ser_alistado.Nick)
	fmt.Println(query)
	_, err = db.Exec(query)

	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		fmt.Println("fudeu bixo deu merda no banco")
		return
	}

	respostas.JSON(w, http.StatusOK, "alistado com sucesso")
}

// porra eu preciso dizer oq isso faz?
func MudarCargo(w http.ResponseWriter, r *http.Request) {

	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
	}

	usuario, err := modelos.PegarUsuarioPeloID(usuarioID)

	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	type new_cargo struct {
		New_cargo                   int
		Usuario_para_mudar_cargo_id int
	}

	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	cargo_data := new_cargo{}
	if erro = json.Unmarshal(corpoRequest, &cargo_data); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuario_para_mudar_cargo, err := modelos.PegarUsuarioPeloID(uint64(cargo_data.Usuario_para_mudar_cargo_id))

	if erro = json.Unmarshal(corpoRequest, &cargo_data); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if cargo_data.New_cargo < 0 || cargo_data.New_cargo > 3 {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Cargo invalido"))
		return
	}

	//dahora que isso handla se ele tentar mudar o proprio cargo lol (racismo sexual com arvores)
	if usuario_para_mudar_cargo.Cargo >= usuario.Cargo {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Voce nao tem permissão para fazer isso!"))
		return
	}

	db, err := banco.Conectar() //banco inter?
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		fmt.Println("fudeu bixo banco caiu")
		return
	}
	query := fmt.Sprintf("update usuarios set cargo = %v where id = %v", cargo_data.New_cargo, cargo_data.Usuario_para_mudar_cargo_id)
	db.Exec(query)
	respostas.JSON(w, http.StatusOK, "mudado com sucesso")
}

// CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario

	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		fmt.Println("sexo2")
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		fmt.Println("sexo3")
		return
	}

	query := `SELECT COUNT(nick) AS users FROM alistados WHERE nick = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		if err != nil {
			respostas.Erro(w, http.StatusInternalServerError, err)
			fmt.Println("fudeu bixo deu merda no banco")
			return
		}
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, usuario.Nick)
	number_of_users := 0
	if err := row.Scan(&number_of_users); err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		fmt.Println("fudeu bixo deu merda no banco")
		return
	}
	fmt.Println(number_of_users)
	if number_of_users == 0 {
		respostas.Erro(w, http.StatusBadRequest, errors.New("Voce nao esta alistado"))
		fmt.Println("sexo2")
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	fmt.Println("b: ", usuario.Senha)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		fmt.Printf("sexo4 %v\n", erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios busca todos os usuários salvos no banco
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, erro := banco.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuario busca um usuário salvo no banco
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario altera as informações de um usuário no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario exclui as informações de um usuário no banco
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja o seu"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuario permite que um usuário siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// PararDeSeguirUsuario permite que um usuário pare de seguir outro
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível parar de seguir você mesmo"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.PararDeSeguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores traz todos os seguidores de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo traz todos os usuários que um determinado usuário está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha permite alterar a senha de um usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não condiz com a que está salva no banco"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
