package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"chatbot-service/app/usecase/chat"
)

type ChatController struct {
	UseCase *chat.ChatUseCase
}

func NewChatController(useCase *chat.ChatUseCase) *ChatController {
	return &ChatController{UseCase: useCase}
}

type ChatRequest struct {
	ClienteID   string `json:"cliente_id"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Handle godoc
// @Summary Atendimento inteligente com IA
// @Description Recebe uma mensagem do usuário e responde com base na IA e no prompt dinâmico do cliente
// @Tags atendimento
// @Accept json
// @Produce json
// @Param message body ChatRequest true "Mensagem do usuário com cliente_id"
// @Success 200 {object} ChatResponse
// @Failure 400 {string} string "Requisição inválida"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /api/v1/atendimento [post]
func (c *ChatController) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Message == "" || req.ClienteID == "" || req.PhoneNumber == "" {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	resp, err := c.UseCase.ProcessarMensagem(ctx, req.ClienteID, req.PhoneNumber, req.Message)
	if err != nil {
		log.Printf("Erro no atendimento IA: %v", err)
		http.Error(w, "Erro ao processar atendimento", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ChatResponse{Response: resp})
}