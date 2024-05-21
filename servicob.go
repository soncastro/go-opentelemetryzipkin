package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/propagation"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// CEP representa os dados retornados pela API ViaCep
type CEP struct {
	Cidade string `json:"localidade"`
}

// Clima representa os dados retornados pela API OpenWeatherMap
type Clima struct {
	Main Temperatura `json:"main"`
}

// Temperatura representa a temperatura e outras informações meteorológicas
type Temperatura struct {
	Temp    float64 `json:"temp"`
	TempMin float64 `json:"temp_min"`
	TempMax float64 `json:"temp_max"`
	Pressao float64 `json:"pressure"`
	Umidade float64 `json:"humidity"`
	TempKf  float64 `json:"temp_kf"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Erro bool `json:"erro"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/weather/{cep}", GetWeatherByCep).Methods("GET")
	port := "8081"
	log.Printf("Servidor Serviço B rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func GetWeatherByCep(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	textoMapPropagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	ctx = textoMapPropagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
	cep := params["cep"]

	// Valida se o CEP possui formato correto
	if !validaCEP(cep) {
		js, err := json.Marshal(map[string]string{"message": "invalid zipcode"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Write(js)
		return
	}

	// Busca cidade do CEP
	dadosCidade, err := buscaCidade(cep)
	if err != nil {
		js, err := json.Marshal(map[string]string{"message": "can not find zipcode"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound) // 404
		w.Write(js)
		return
	}

	// Busca o clima da cidade
	clima, err := buscaClima(strings.ToLower(dadosCidade.Cidade))
	if err != nil {
		fmt.Println("Erro na busca do clima:", err)
		return
	}

	resp := WeatherResponse{
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

func buscaCidade(cep string) (*CEP, error) {
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

func buscaClima(cidade string) (*Clima, error) {
	cidade = url.QueryEscape(cidade)
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + cidade + ",br&appid=904020cdcc44973b1dd0810487a25068")
	fmt.Println(cidade)
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

func convCtoK(tempC float64) float64 {
	return tempC + 273.15
}
