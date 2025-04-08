
# 🤖 Chatbot Whitelabel com OpenAI + Bucket

Este projeto é um serviço de chatbot inteligente, multilíngue e whitelabel, que utiliza a API da OpenAI e busca configurações específicas de cada cliente a partir de arquivos JSON armazenados em um bucket (como S3, GCS, etc).

---

## 🚀 Tecnologias Utilizadas

- **Go (Golang)** — backend rápido e eficiente
- **OpenAI API** — geração de linguagem natural
- **Docker / Docker Compose** — ambiente padronizado
- **Swagger (Swaggo)** — documentação da API
- **Arquitetura limpa (DDD)** — domínio separado da infraestrutura
- **Buckets configuráveis** — fonte externa de prompts e serviços

---

## 🐳 Como Rodar com Docker

1. Crie um arquivo `.env` (veja exemplo abaixo)

2. Execute:

```bash
docker-compose up --build

	3.	Acesse a API via:

http://localhost:8080/atendimento

	4.	Acesse a documentação Swagger:

http://localhost:8080/swagger/index.html



⸻

📄 Exemplo de .env

PORT=8080

# OpenAI
OPENAI_API_KEY=chave_com_fornecimento_interno
OPENAI_MODEL=gpt-4o
OPENAI_API_URL=https://api.openai.com/v1/chat/completions

# Bucket
BUCKET_BASE_URL=https://seu-bucket.s3.amazonaws.com/prompts
PROMPT_FILE_PATTERN=%s-prompt.json



⸻

🛡️ Segurança
	•	As chaves da OpenAI NÃO estão inclusas neste repositório.
	•	Solicite a OPENAI_API_KEY diretamente ao departamento de segurança da informação.
	•	Nunca versionar .env em ambientes sensíveis.

⸻

🧪 Testes

go test ./... -v

Incluímos mocks e dados em app/test/testdata/ para testes automatizados com precisão.

⸻

📂 Estrutura de Pastas

├── app
│   ├── domain         # Contratos de domínio
│   ├── usecases       # Regras de negócio
│   └── test           # Testes e mocks
├── dependencies       # Integrações externas (OpenAI, Bucket)
├── presentation       # Controllers e rotas HTTP
├── main.go            # Entry point
├── docker-compose.yml
└── .env               # Variáveis de ambiente (local)



⸻

👨‍💻 Contato

Para dúvidas técnicas ou integração com n8n / WhatsApp, procure o time de engenharia ou envie um e-mail para:

suporte@chatbotwhitelabel.com



⸻



"Simples, simpático e configurável. Um chatbot que entende seu cliente." 💬

---

Env Example

# Porta do servidor HTTP
PORT=8080

# -------------------------------
# 🔐 Configuração da OpenAI
# -------------------------------
# Solicite a chave ao departamento de segurança
OPENAI_API_KEY=sua_openai_api_key_aqui
OPENAI_MODEL=gpt-4o
OPENAI_API_URL=https://api.openai.com/v1/chat/completions

# -------------------------------
# ☁️ Configuração do Bucket
# -------------------------------
BUCKET_BASE_URL=https://seu-bucket.s3.amazonaws.com/prompts
PROMPT_FILE_PATTERN=%s-prompt.json