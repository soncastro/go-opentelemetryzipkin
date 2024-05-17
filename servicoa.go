package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type InputCEP struct {
	CEP string `json:"cep"`
}

func servicoa() {
	router := mux.NewRouter()
	router.HandleFunc("/cep", processaCEP).Methods("POST")
	port := "8080"
	log.Printf("Servidor Serviço A rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func processaCEP(w http.ResponseWriter, req *http.Request) {
	tracer := otel.GetTracerProvider().Tracer("servicoa")
	ctx := context.Background()
	ctx, processSpan := tracer.Start(ctx, "processaCEP")
	defer processSpan.End()

	fmt.Println(req.Body)

	var cep InputCEP
	err := json.NewDecoder(req.Body).Decode(&cep)
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
		cepSpan.End()

		var client http.Client
		request, err := http.NewRequestWithContext(cepCtx, "GET", "http://localhost:8081/weather/"+cep.CEP, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		textMapPropagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		textMapPropagator.Inject(cepCtx, propagation.HeaderCarrier(request.Header))

		_, weatherSpan := tracer.Start(cepCtx, "buscaTemperatura")
		resp, err := client.Do(request)
		weatherSpan.End()

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
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Write(js)
	}
}
