# Desafio Prático - Tracing distribuído e span  

## Como rodar em ambiente DEV
* Alternativa 1: Na raiz do projeto usar o comando `docker-compose up`  
* Alternativa 2: Usar o comando `go run main.go servicoa.go servicob.go`    

## Realização de teste
Use o comando `docker-compose up` na raiz do projeto para subir os serviços e faça uma requisição HTTP POST para http://localhost:8080/cep com o conteúdo abaixo para o Body:  
`{ "cep": "29902555" }`  
Troque a numeração por qual desejar.   
O Postman pode ser uma alternativa de ferramenta para o teste.