package main

import (
	"encoding/json"
	"fmt"
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
	"net/url"
	"regexp"
	"strings"
)

type CEP struct {
	Cidade string `json:"localidade"`
}

type Clima struct {
	Main Temperatura `json:"main"`
}

type Temperatura struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
	Pressao float64 `json:"pressure"`
	Umidade float64 `json:"humidity"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Erro bool `json:"erro"`
}

const endpointURL_servicob = "http://zipkin:9411/api/v2/spans"

func main() {
	initTracer2()
	port := "8081"

	router := mux.NewRouter()
	router.HandleFunc("/weather/{cep}", GetWeatherByCep).Methods("GET")

	log.Printf("Servidor Serviço B rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initTracer2() {
	exporter, err := zipkin.New(endpointURL_servicob)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("servicob"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func GetWeatherByCep(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("servicob")
	ctx := r.Context()
	textMapPropagator := propagation.TraceContext{}
	ctx = textMapPropagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, span := tracer.Start(ctx, "GetWeatherByCep")
	defer span.End()
	log.Println("Span GetWeatherByCep iniciado")

	params := mux.Vars(r)
	cep := params["cep"]

	if !validaCEP(cep) {
		js, err := json.Marshal(map[string]string{"message": "invalid zipcode"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(js)
		return
	}

	dadosCidade, err := buscaCidade(cep, r)
	if err != nil {
		js, err := json.Marshal(map[string]string{"message": "can not find zipcode"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(js)
		return
	}
	log.Println("Cidade encontrada:", dadosCidade.Cidade)

	clima, err := buscaClima(strings.ToLower(dadosCidade.Cidade), r)
	if err != nil {
		fmt.Println("Erro na busca do clima:", err)
		return
	}
	log.Println("Clima encontrado para a cidade:", dadosCidade.Cidade)

	resp := WeatherResponse{
		City:  dadosCidade.Cidade,
		TempC: convKELtoC(clima.Main.Temp),
		TempF: convKELtoF(clima.Main.Temp),
		TempK: clima.Main.Temp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func validaCEP(cep string) bool {
	match, _ := regexp.MatchString("^[0-9]{8}$", cep)
	return match
}

func buscaCidade(cep string, r *http.Request) (*CEP, error) {

	tracer := otel.Tracer("servicob")
	ctx := r.Context()
	textMapPropagator := propagation.TraceContext{}
	ctx = textMapPropagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, span := tracer.Start(ctx, "buscaCidade")
	defer span.End()
	log.Println("Span buscaCidade iniciado")

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dados, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var erroResponse ErrorResponse
	err = json.Unmarshal(dados, &erroResponse)
	if err == nil && erroResponse.Erro {
		err = fmt.Errorf("A resposta apresentou erro: %v", erroResponse)
		return nil, err
	}

	var c CEP
	err = json.Unmarshal(dados, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func buscaClima(cidade string, r *http.Request) (*Clima, error) {

	tracer := otel.Tracer("servicob")
	ctx := r.Context()
	textMapPropagator := propagation.TraceContext{}
	ctx = textMapPropagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, span := tracer.Start(ctx, "buscaClima")
	defer span.End()
	log.Println("Span buscaClima iniciado")

	cidade = url.QueryEscape(cidade)
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + cidade + ",br&appid=904020cdcc44973b1dd0810487a25068")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dados, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var clima Clima
	err = json.Unmarshal(dados, &clima)
	if err != nil {
		return nil, err
	}
	return &clima, nil
}

func convKELtoC(tempK float64) float64 {
	return tempK - 273.15
}

func convKELtoF(tempK float64) float64 {
	return (tempK-273.15)*1.8 + 32
}
