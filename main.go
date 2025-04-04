package main

import (
	"log"
	"github.com/joho/godotenv"
	"chatbot-service/db"
	"chatbot-service/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	err = db.Conectar()
	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	services.AtendimentoIA()
}