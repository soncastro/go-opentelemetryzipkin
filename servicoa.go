package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const endpointURL = "http://zipkin:9411/api/v2/spans"

type InputCEP struct {
	CEP string `json:"cep"`
}

func main() {
	initTracer()
	port := "8080"

	router := mux.NewRouter()
	router.HandleFunc("/cep", processaCEP).Methods("POST")

	log.Printf("Servidor Serviço A rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initTracer() {
	exporter, err := zipkin.New(endpointURL)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("servicoa"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func processaCEP(w http.ResponseWriter, req *http.Request) {
	tracer := otel.Tracer("servicoa")
	ctx := req.Context()
	ctx, processSpan := tracer.Start(ctx, "processaCEP")
	defer processSpan.End()
	log.Println("Span processaCEP iniciado")

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var cep InputCEP
	err = json.Unmarshal(bodyBytes, &cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regexp, err := regexp.Compile("^[0-9]{8}$")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if regexp.MatchString(cep.CEP) {
		cepCtx, cepSpan := tracer.Start(ctx, "validaCEP")
		log.Println("Span validaCEP iniciado")
		defer cepSpan.End()

		client := &http.Client{}
		request, err := http.NewRequestWithContext(cepCtx, "GET", "http://servicob:8081/weather/"+cep.CEP, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		propagation.TraceContext{}.Inject(cepCtx, propagation.HeaderCarrier(request.Header))
		log.Println("Contexto de tracing injetado na requisição para servicob")

		resp, err := client.Do(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else {
		js, err := json.Marshal(map[string]string{"message": "invalid zipcode"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(js)
	}
}
