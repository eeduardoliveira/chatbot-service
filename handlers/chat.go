package handlers

import (
	"encoding/json"
	"net/http"

	"chatbot-service/openai"
	"chatbot-service/types"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req types.ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Message == "" {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	reply, err := openai.GetChatResponse(req.Message)
	if err != nil {
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(types.ChatResponse{Response: reply})
}