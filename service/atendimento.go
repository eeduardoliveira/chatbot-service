package service

import (
	"chatbot-service/db"
	"chatbot-service/openai"
	"fmt"
)

var historico []openai.Message

func AtendimentoIA(prompt string) (string, error) {

	servicos, err := db.BuscarServicos()
	if err != nil {
		return  "", fmt.Errorf("Erro ao buscar serviços: %w", err)
	}
	
	descricaoServicos := "A clínica oferece os seguintes serviços:\n"
	for _, s := range servicos {
		descricaoServicos += fmt.Sprintf("- %s: %s — R$ %.2f\n", s.Nome, s.Descricao, s.Preco)
	}
	
	systemPrompt := fmt.Sprintf(`Você é uma atendente virtual chamada Ana, da Clínica Sorriso Ideal.
	Fale de forma simpática e acolhedora.
	Evite ser prolixa, seja simpatica mas direta.
	Converse de forma dinamica
	Falas inglês e português
	%s
	
	A clínica atende das 08h às 18h, de segunda a sábado.
	
	Você deve responder com base nessas informações. Se o paciente perguntar algo fora do escopo ou solicitar valores exatos e personalizados, oriente que será necessário entrar em contato com a recepção.`, descricaoServicos)
	
	
	if len(historico) == 0 {
		historico = append(historico, openai.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	var resposta string
	historico, resposta, err = openai.PerguntarAoChatbot(prompt, historico)
	return resposta, err
}