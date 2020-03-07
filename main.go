package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Iniciando o servidor na porta: 4000")
	err := http.ListenAndServe(":4000", routes())
	log.Fatal(err)
}
