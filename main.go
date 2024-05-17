package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"sync"
)

func main() {

	initTracer()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		servicob()
	}()

	go func() {
		defer wg.Done()
		servicoa()
	}()

	wg.Wait()
}

func initTracer() {
	// Inicializa o exportador Zipkin com a URL do servidor Zipkin
	log.Printf("Servidor Zipkin rodando na porta 9411")
	exp, err := zipkin.NewRawExporter("http://zipkin:9411/api/v2/spans")

	if err != nil {
		log.Fatalf("erro ao criar o exportador Zipkin: %v", err)
	}

	// Cria um novo provedor de trace utilizando o exportador Zipkin.
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exp), sdktrace.WithSampler(sdktrace.AlwaysSample()))

	// Defina o provedor de trace global como o provider que acabamos de criar.
	otel.SetTracerProvider(tp)
}
