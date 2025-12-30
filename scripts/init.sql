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

-- Criar tabela usuario (nova - sistema de autenticação email+senha)
CREATE TABLE IF NOT EXISTS usuario (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL, -- Hash bcrypt
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    pessoa_id INTEGER REFERENCES pessoa(id) ON DELETE SET NULL,
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_desativacao TIMESTAMP,
    is_ativo BOOLEAN DEFAULT TRUE
);

-- Criar índices
CREATE INDEX IF NOT EXISTS idx_pessoa_documento ON pessoa(documento);
CREATE INDEX IF NOT EXISTS idx_cliente_ativo ON cliente(ativo);
CREATE INDEX IF NOT EXISTS idx_cliente_pessoa_id ON cliente(pessoa_id);
CREATE INDEX IF NOT EXISTS idx_usuario_email ON usuario(email);
CREATE INDEX IF NOT EXISTS idx_usuario_ativo ON usuario(is_ativo);

-- Inserir dados de teste (pessoa e cliente - sistema legado CPF)
INSERT INTO pessoa (nome, email, documento, tipo_pessoa, data_nascimento)
VALUES
    ('João Silva', 'joao.silva@example.com', '12345678909', 'FISICA', '1990-01-15'),
    ('Maria Santos', 'maria.santos@example.com', '98765432100', 'FISICA', '1985-05-20'),
    ('Cliente Teste', 'teste@example.com', '11122233344', 'FISICA', '1995-03-10')
ON CONFLICT (documento) DO NOTHING;

INSERT INTO cliente (pessoa_id, ativo)
SELECT id, true FROM pessoa WHERE documento IN ('12345678909', '98765432100', '11122233344')
ON CONFLICT (pessoa_id) DO NOTHING;

-- Inserir usuários de teste (novo sistema email+senha)
-- Senhas:
--   admin@oficinapro.com: Admin@123
--   client@oficinapro.com: Client@123
--   user@oficinapro.com: senha123
--   joao.silva@example.com: senha123
INSERT INTO usuario (email, senha, role, pessoa_id, is_ativo)
VALUES
    ('admin@oficinapro.com', '$2a$10$1tR72UJaZdEqCKWpCxbqY.g727uQ.yL/JLyyxM2zweH6Kd7wy4yHW', 'ADMIN', NULL, true),
    ('client@oficinapro.com', '$2a$10$KZmncud3SXrMU4rQ8osCNe1BpptXP1RRl9AtctWCqssV7.ONEC7NG', 'USER', NULL, true),
    ('user@oficinapro.com', '$2a$10$rO7Y6bC9vZF8kz5TqzN5zeN.VQJZKxRGQxC/rLvP2x4YC5xZ8yVWS', 'USER', NULL, true),
    ('joao.silva@example.com', '$2a$10$rO7Y6bC9vZF8kz5TqzN5zeN.VQJZKxRGQxC/rLvP2x4YC5xZ8yVWS', 'USER', (SELECT id FROM pessoa WHERE email = 'joao.silva@example.com'), true)
ON CONFLICT (email) DO NOTHING;

-- Verificar dados inseridos (cliente - sistema legado)
SELECT
    c.id as cliente_id,
    p.nome,
    p.documento,
    c.ativo
FROM cliente c
JOIN pessoa p ON c.pessoa_id = p.id;

-- Verificar usuários inseridos (novo sistema)
SELECT
    id,
    email,
    role,
    pessoa_id,
    is_ativo,
    data_criacao
FROM usuario
ORDER BY id;

