package services

import (
	"chatbot-service/db"
	"chatbot-service/models"
	"database/sql"
	"fmt"
)

func AtenderCliente() {
	var entrada string

	fmt.Println("Olá, seja bem-vindo à clínica!")
	fmt.Print("Por favor, informe seu CPF ou código de paciente: ")
	fmt.Scanln(&entrada)

	var paciente models.Paciente
	err := db.DB.QueryRow(`
		SELECT id, nome, cpf, telefone, email, to_char(data_nascimento, 'DD/MM/YYYY')
		FROM pacientes
		WHERE cpf = $1 OR CAST(id AS TEXT) = $1`, entrada).
		Scan(&paciente.ID, &paciente.Nome, &paciente.CPF, &paciente.Telefone, &paciente.Email, &paciente.Nascimento)

	if err == sql.ErrNoRows {
		fmt.Println("Paciente não encontrado.")
		return
	} else if err != nil {
		fmt.Println("Erro ao buscar paciente:", err)
		return
	}

	fmt.Println("\n--- Dados do Paciente ---")
	fmt.Printf("ID: %d\n", paciente.ID)
	fmt.Printf("Nome: %s\n", paciente.Nome)
	fmt.Printf("CPF: %s\n", paciente.CPF)
	fmt.Printf("Telefone: %s\n", paciente.Telefone)
	fmt.Printf("Email: %s\n", paciente.Email)
	fmt.Printf("Nascimento: %s\n", paciente.Nascimento)
	fmt.Println("--------------------------")
}