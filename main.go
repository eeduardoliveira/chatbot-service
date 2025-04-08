package main

import (
	"log"
	"net/http"
	"os"
	"chatbot-service/handlers"
	"github.com/joho/godotenv"
	"chatbot-service/db"
	"github.com/swaggo/http-swagger"
	_ "chatbot-service/docs"
)

// @title Chatbot da Clínica Odontológica
// @version 1.0
// @description API para integração com OpenAI e base de pacientes/serviços
// @host localhost:8080
// @BasePath /
// @contact.name Suporte Técnico
// @contact.email suporte@clinicasorriso.com
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	err = db.Conectar()
	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	http.HandleFunc("/atendimento", handlers.AtendimentoHandler)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando em http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}