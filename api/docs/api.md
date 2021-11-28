# Documentação da API

## Usuarios

### Criar usuario

<table>
  <tr>
    <td>URI</td>
    <td>Autenticação</td>
    <td>Metodo</td>
  </tr>
  <tr>
    <td>/usuarios</td>
    <td>false</td>
    <td>POST</td>
  </tr>
</table>

body
```json
{
    "nome": "string",
    "nick": "string",
    "senha": "string",
    "email": "string",
    "codigo": "string" //codigo de verificação
}
```

### Alterar Cargo

<table>
  <tr>
    <td>URI</td>
    <td>Autenticação</td>
    <td>Metodo</td>
  </tr>
  <tr>
    <td>/usuarios/cargo</td>
    <td>true</td>
    <td>POST</td>
  </tr>
</table>

body
```json
{
    "nick": "string",
    "new_cargo": 0 - 2 //0 ate o 2
}
```

### Alistar Usuario

<table>
  <tr>
    <td>URI</td>
    <td>Autenticação</td>
    <td>Metodo</td>
  </tr>
  <tr>
    <td>/usuarios/alistar</td>
    <td>true</td>
    <td>PUT</td>
  </tr>
</table>

body
```json
{
    "nick": "string",
}
```
