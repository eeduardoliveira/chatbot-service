package handlers

import (
	"encoding/json"
	"net/http"

	"chatbot-service/openai"
)

type AtendimentoRequest struct {
	Message string `json:"message"`
}

type AtendimentoResponse struct {
	Response string `json:"response"`
}

// AtendimentoHandler godoc
// @Summary Atendimento inteligente com IA
// @Description Recebe uma mensagem do usuário e responde com base na IA e nos dados da clínica
// @Tags atendimento
// @Accept json
// @Produce json
// @Param message body AtendimentoRequest true "Mensagem do usuário"
// @Success 200 {object} AtendimentoResponse
// @Failure 400 {string} string "Requisição inválida"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /atendimento [post]
func AtendimentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // útil se rodar em outro domínio

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req AtendimentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Message == "" {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	resp, err := openai.PerguntarAoChatbot(req.Message)
	if err != nil {
		http.Error(w, "Erro ao processar atendimento", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AtendimentoResponse{Response: resp})
}
