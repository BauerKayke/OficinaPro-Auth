-- Script de inicialização do banco de dados para desenvolvimento local

-- Criar tabela pessoa
CREATE TABLE IF NOT EXISTS pessoa (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    documento VARCHAR(14) NOT NULL UNIQUE,
    tipo_pessoa VARCHAR(20) NOT NULL,
    data_nascimento DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar tabela cliente
CREATE TABLE IF NOT EXISTS cliente (
    id SERIAL PRIMARY KEY,
    pessoa_id INTEGER NOT NULL REFERENCES pessoa(id) ON DELETE CASCADE,
    ativo BOOLEAN DEFAULT TRUE,
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pessoa_id)
);

-- Criar índices
CREATE INDEX IF NOT EXISTS idx_pessoa_documento ON pessoa(documento);
CREATE INDEX IF NOT EXISTS idx_cliente_ativo ON cliente(ativo);
CREATE INDEX IF NOT EXISTS idx_cliente_pessoa_id ON cliente(pessoa_id);

-- Inserir dados de teste
INSERT INTO pessoa (nome, email, documento, tipo_pessoa, data_nascimento)
VALUES
    ('João Silva', 'joao.silva@example.com', '12345678909', 'FISICA', '1990-01-15'),
    ('Maria Santos', 'maria.santos@example.com', '98765432100', 'FISICA', '1985-05-20'),
    ('Cliente Teste', 'teste@example.com', '11122233344', 'FISICA', '1995-03-10')
ON CONFLICT (documento) DO NOTHING;

INSERT INTO cliente (pessoa_id, ativo)
SELECT id, true FROM pessoa WHERE documento IN ('12345678909', '98765432100', '11122233344')
ON CONFLICT (pessoa_id) DO NOTHING;

-- Verificar dados inseridos
SELECT
    c.id as cliente_id,
    p.nome,
    p.documento,
    c.ativo
FROM cliente c
JOIN pessoa p ON c.pessoa_id = p.id;

