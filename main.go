package main

import (
	"log"
	"net/http"
	"os"

	"chatbot-service/app/usecase/chat"
	"chatbot-service/dependencies/bucket"
	"chatbot-service/dependencies/openai"
	"chatbot-service/presentation/controller"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	_ "chatbot-service/docs"
)

// @title Chatbot 
// @version 2.0
// @description API para atendimento inteligente com IA e configuração via bucket
// @host localhost:8080
// @BasePath /
// @contact.name Equipe Técnica
// @contact.email suporte@chatbotwhitelabel.com
func main() {
	// Carrega variáveis do .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	// Injeção de dependências
	bucketService := bucket.NewHTTPBucket()
	promptRepo := bucket.NewBucketPromptRepository(bucketService)
	openaiClient := openai.NewClient()
	chatUseCase := chat.NewChatUseCase(promptRepo, openaiClient)
	chatController := controller.NewChatController(chatUseCase)

	// Configura endpoints
	http.HandleFunc("/atendimento", chatController.Handle)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando em http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}