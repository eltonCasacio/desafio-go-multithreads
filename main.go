package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type CDN struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func main() {
	cdn := make(chan CDN)
	viacep := make(chan ViaCEP)

	go func() {
		c, err := http.Get("https://cdn.apicep.com/file/apicep/13277-705.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar CEP - https://cdn.apicep.com/file/apicep/13277-705.json: %v\n", err)
		}
		defer c.Body.Close()

		res, err := io.ReadAll(c.Body)
		if err != nil {
			panic(err)
		}

		var data CDN
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da responsa: %v\n", err)
		}
		cdn <- data
	}()

	go func() {
		c, err := http.Get("http://viacep.com.br/ws/13277705/json/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar CEP - http://viacep.com.br/ws/13277705/json/: %v\n", err)
		}
		defer c.Body.Close()

		res, err := io.ReadAll(c.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da responsa: %v\n", err)
		}

		var data ViaCEP
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da responsa: %v\n", err)
		}
		viacep <- data
	}()

	select {
	case cdnMessage := <-cdn:
		fmt.Fprintf(os.Stdout, "CDN:: %v\n", cdnMessage)
	case viacepMessage := <-viacep:
		fmt.Fprintf(os.Stdout, "VIACEP:: %v\n", viacepMessage)
	case <-time.After(time.Second):
		fmt.Println("Excedeu tempo limite de resposta")
	}
}
