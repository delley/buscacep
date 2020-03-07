package main

type cep struct {
	Cep        string `json:"cep"`
	Cidade     string `json:"cidade,omitempty"`
	Bairro     string `json:"bairro,omitempty"`
	Logradouro string `json:"logradouro,omitempty"`
	UF         string `json:"uf,omitempty"`
}
