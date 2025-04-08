-- Criação da tabela de pacientes
CREATE TABLE IF NOT EXISTS pacientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    telefone VARCHAR(20),
    email VARCHAR(100),
    data_nascimento DATE
);

-- Inserção de pacientes exemplo
INSERT INTO pacientes (nome, cpf, telefone, email, data_nascimento) VALUES
('João da Silva', '123.456.789-00', '(11) 99999-1111', 'joao.silva@example.com', '1985-03-10'),
('Maria Oliveira', '987.654.321-00', '(11) 98888-2222', 'maria.oliveira@example.com', '1990-07-22'),
('Carlos Santos', '321.654.987-00', '(11) 97777-3333', 'carlos.santos@example.com', '1978-11-05');

-- Criação da tabela de serviços
CREATE TABLE IF NOT EXISTS servicos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT NOT NULL,
    preco NUMERIC(10,2) NOT NULL
);

-- Inserção de serviços
INSERT INTO servicos (nome, descricao, preco) VALUES
('Limpeza', 'Limpeza e profilaxia dental', 150.00),
('Clareamento', 'Clareamento dental profissional', 600.00),
('Implantes', 'Implantes dentários por unidade', 2500.00),
('Ortodontia', 'Aparelho ortodôntico com manutenção mensal', 200.00),
('Odontopediatria', 'Consulta com dentista infantil', 180.00);