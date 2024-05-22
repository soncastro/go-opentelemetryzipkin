# Desafio Prático - Tracing distribuído e span

## Como rodar em ambiente DEV
* Na raiz do projeto usar o comando `docker-compose up --build`

## Realização de teste
- Use o comando `docker-compose up --build` na raiz do projeto para subir os serviços.
- Faça uma requisição HTTP POST para http://localhost:8080/cep com o Body abaixo:  
  `{ "cep": "29902555" }`
- Troque a numeração do CEP por qual desejar.
- O Postman pode ser uma alternativa de ferramenta para o teste.
- Acesse pelo browser o endereço http://localhost:9411/zipkin para visualizar os traces da aplicação após os requests.

## Última solicitação de correção (Enviado 21/05/2024 • 10:13 - Gabriel Araujo Carneiro Junior)
- [x] Colocar no padrão: `{ "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`
- [x] Nesta versão foi conferido que está sendo possível visualizar os traces ao entrar no zipkin na porta 9411.

## ~~Solicitação de correção (Enviado 20/05/2024 • 13:32 - Gabriel Araujo Carneiro Junior)~~
- [x] Serviços separados.
- [x] docker-compose.yaml atualizado para atender a demanda de serviços.
- [x] main.go, servicoa.go e servicob.go também atualizado pelo motivo acima.
