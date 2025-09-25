package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type BrasilCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func callViaCep(cep string, client *http.Client, c chan<- ViaCep) (*ViaCep, error) {
	viaCepUrl := "http://viacep.com.br/ws/" + cep + "/json/"
	requestViaCep, _ := http.NewRequest("GET", viaCepUrl, nil)
	result, err := client.Do(requestViaCep)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var viaCep ViaCep
	err = json.Unmarshal(body, &viaCep)
	if err != nil {
		return nil, err
	}

	c <- viaCep

	return &viaCep, nil
}

func callBrasilCep(cep string, client *http.Client, c chan BrasilCep) (*BrasilCep, error) {
	brasilCepUrl := "https://brasilapi.com.br/api/cep/v1/" + cep
	requestBrasilCep, _ := http.NewRequest("GET", brasilCepUrl, nil)

	result, err := client.Do(requestBrasilCep)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	var brasilCep BrasilCep
	err = json.Unmarshal(body, &brasilCep)
	if err != nil {
		return nil, err
	}

	c <- brasilCep

	return &brasilCep, nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	viaCep := make(chan ViaCep)
	brasilCep := make(chan BrasilCep)

	r.Get("/search/{cep}", func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")
		if cep == "" {
			w.Write([]byte("missing param: cep"))
			w.WriteHeader(http.StatusBadRequest)
		}

		client := http.Client{}
		go callViaCep(cep, &client, viaCep)
		go callBrasilCep(cep, &client, brasilCep)

		select {
		case vc := <-viaCep:
			json.NewEncoder(w).Encode(vc)
			log.Println("ViaCep wins!")
			w.WriteHeader(http.StatusOK)
		case bc := <-brasilCep:
			json.NewEncoder(w).Encode(bc)
			log.Println("BrasilCep wins!")
			w.WriteHeader(http.StatusOK)
		case <-time.After(time.Second):
			w.Write([]byte("Timeout Error"))
			log.Println("Both loses!")
			w.WriteHeader(http.StatusGatewayTimeout)
		}
	})

	http.ListenAndServe(":3000", r)
}
