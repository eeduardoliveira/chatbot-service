package main

import (
	"log"
	"net/http"
	"os"

	"chatbot-service/app/usecase/chat"
	"chatbot-service/dependencies/bucket"
	"chatbot-service/dependencies/openai"
	redisrepo "chatbot-service/dependencies/redis"
	"chatbot-service/presentation/controller"

	_ "chatbot-service/docs"

	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
	"github.com/redis/go-redis/v9"
)

// @title Chatbot
// @version 2.0.5
// @description API para atendimento inteligente com IA e configuração via bucket
// @host chatbot.syphertech.com.br
// @BasePath /
// @contact.name Equipe Técnica
// @contact.email sypher.infraestrutura@gmail.com
func main() {

	// Injeção de dependências
	bucketService := bucket.NewHTTPBucket()
	promptRepo := bucket.NewBucketPromptRepository(bucketService)
	openaiClient := openai.NewClient()

	// Configuração do Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		panic("REDIS_ADDR não configurado")

	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	sessionRepo := redisrepo.NewRedisSessionRepository(redisClient)
	contextRepo := redisrepo.NewRedisContextRepository(redisClient)
	chatUseCase := chat.NewChatUseCase(promptRepo, openaiClient, sessionRepo, contextRepo)
	chatController := controller.NewChatController(chatUseCase)

	// Configura endpoints
	http.HandleFunc("/api/v1/atendimento", chatController.Handle)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}