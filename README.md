# Desafio Prático - Tracing distribuído e span  

## Como rodar em ambiente DEV
* Na raiz do projeto usar o comando `docker-compose up --build`   

## Realização de teste
Use o comando `docker-compose up --build` na raiz do projeto para subir os serviços.  
Faça uma requisição HTTP POST para http://localhost:8080/cep com o Body abaixo:  
`{ "cep": "29902555" }`  
Troque a numeração do CEP por qual desejar.   
O Postman pode ser uma alternativa de ferramenta para o teste.

## Última solicitação de correção (Enviado 20/05/2024 • 13:32 - Gabriel Araujo Carneiro Junior)
- [x] Serviços separados.
- [x] docker-compose.yaml atualizado para atender a demanda de serviços.
- [x] main.go, servicoa.go e servicob.go também atualizado pelo motivo acima.

## Falta apenas este projeto para concluir a pós. Torcendo pela aprovação :)