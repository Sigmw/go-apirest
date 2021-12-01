var captcha_message = document.getElementById("captcha_message")
var captcha_code = document.getElementById("captcha_code")

//codigo roubado fodase, vai fzr oq?
const random = (length = 8) => {
    // Declare all characters
    let chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

    // Pick characers randomly
    let str = '';
    for (let i = 0; i < length; i++) {
        str += chars.charAt(Math.floor(Math.random() * chars.length));
    }

    return str;

};
$('#formulario-cadastro').on('submit', criarUsuario);
//cria um codigo
let code = `SIGMA-${random(5)}`;
captcha_message.innerText = `CAPTCHA: coloque isso como missão no habbo`;
captcha_code.innerText = code;

function criarUsuario(evento) {
    evento.preventDefault();


    //sexo
    if ($('#senha').val() != $('#confirmar-senha').val()) {
        Swal.fire("Ops...", "As senhas não coincidem!", "error");
        return;
    }
    //envia o request pra cadastrar o
    $.ajax({
        url: "/usuarios",
        method: "POST",
        data: {
           nome: $('#nome').val(), 
           email: $('#email').val(),
           nick: $('#nick').val(),
           senha: $('#senha').val(),
           codigo: code
        }
    
    }) /*se der certo ele chama essa funcao*/.done(function() {
        Swal.fire("Sucesso!", "Usuário cadastrado com sucesso!", "success")
            .then(function() {
                $.ajax({
                    url: "/login",
                    method: "POST",
                    data: {
                        email: $('#email').val(),
                        senha: $('#senha').val()
                    }
                }).done(function() {
                    window.location = "/home";
                }).fail(function() {
                    Swal.fire("Ops...", "Erro ao autenticar o usuário!", "error");
                })
            })
    })/*senao ele chama essa*/.fail(function(err) {
        Swal.fire /*sexo*/("Ops...", `Erro ao cadastrar o usuário! ${err.responseJSON.erro}`, "error" /*sexo final*/);
    });
}