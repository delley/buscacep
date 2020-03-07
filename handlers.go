package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var endpoints = map[string]string{
	"viacep":           "https://viacep.com.br/ws/%s/json/",
	"postmon":          "https://api.postmon.com.br/v1/cep/%s",
	"republicavirtual": "https://republicavirtual.com.br/web_cep.php?cep=%s&formato=json",
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem vindo a API de CEPs"))
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	rCep := r.URL.Path[len("/cep/"):]
	rCep, err := sanitizeCEP(rCep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch := make(chan []byte, 1)
	for _, url := range endpoints {
		endpoint := fmt.Sprintf(url, rCep)
		go request(endpoint, ch)
	}

	w.Header().Set("Content-Type", "application/json")
	for index := 0; index < 3; index++ { // poderia ser de 0 até len(endpoints)
		cepInfo, err := parseResponse(<-ch)
		log.Println(index)
		if err != nil {
			continue
		}

		if cepInfo.exist() {
			cepInfo.Cep = rCep
			json.NewEncoder(w).Encode(cepInfo)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}
