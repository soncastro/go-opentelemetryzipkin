# Dockerfile
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080
EXPOSE 8081

CMD ["./main"]