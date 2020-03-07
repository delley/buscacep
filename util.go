package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func (c cep) exist() bool {
	return len(c.UF) != 0
}

func parseResponse(content []byte) (payload cep, err error) {
	response := make(map[string]interface{})
	_ = json.Unmarshal(content, &response)

	if err := isValidResponse(response); !err {
		return payload, errors.New("invalid response")
	}

	if _, ok := response["localidade"]; ok {
		payload.Cidade = response["localidade"].(string)
	} else {
		payload.Cidade = response["cidade"].(string)
	}

	if _, ok := response["estado"]; ok {
		payload.UF = response["estado"].(string)
	} else {
		payload.UF = response["uf"].(string)
	}

	if _, ok := response["logradouro"]; ok {
		payload.Logradouro = response["logradouro"].(string)
	}

	if _, ok := response["tipo_logradouro"]; ok {
		payload.Logradouro = response["tipo_logradouro"].(string) + " " + payload.Logradouro
	}

	payload.Bairro = response["bairro"].(string)

	return
}

func isValidResponse(requestContent map[string]interface{}) bool {
	if len(requestContent) <= 0 {
		return false
	}

	if _, ok := requestContent["erro"]; ok {
		return false
	}

	if _, ok := requestContent["fail"]; ok {
		return false
	}

	return true
}

func sanitizeCEP(cep string) (string, error) {
	re := regexp.MustCompile(`[^0-9]`)
	sanitizedCEP := re.ReplaceAllString(cep, `$1`)

	if len(sanitizedCEP) < 8 {
		return "", errors.New("O CEP deve conter apenas números e no minimo 8 digitos")
	}

	return sanitizedCEP[:8], nil
}

func request(endpoint string, ch chan []byte) {
	start := time.Now()

	c := http.Client{Timeout: time.Duration(time.Millisecond * 300)}
	resp, err := c.Get(endpoint)
	if err != nil {
		log.Printf("Ops! ocorreu um erro: %s", err.Error())
		ch <- nil
		return
	}
	defer resp.Body.Close()

	requestContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ops! ocorreu um erro: %s", err.Error())
		ch <- nil
		return
	}

	if len(requestContent) != 0 && resp.StatusCode == http.StatusOK {
		log.Printf("O endpoint respondeu com sucesso - source: %s, Duração: %s", endpoint, time.Since(start).String())
		ch <- requestContent
	}
}
