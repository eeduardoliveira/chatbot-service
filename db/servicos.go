package db

import (
	"fmt"
)

type Servico struct {
	Nome  string
	Descricao string
	Preco float64
}

func BuscarServicos() ([]Servico, error) {
	rows, err := DB.Query("SELECT nome, descricao, preco FROM servicos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servicos []Servico
	for rows.Next() {
		var s Servico
		err := rows.Scan(&s.Nome, &s.Descricao, &s.Preco)
		if err != nil {
			return nil, err
		}
		servicos = append(servicos, s)
	}
	return servicos, nil
}