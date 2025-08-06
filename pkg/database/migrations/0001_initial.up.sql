-- Создание таблицы пользователей
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Создание таблицы сессий
CREATE TABLE sessions (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    ip INET,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_session_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Индексы
CREATE INDEX idx_user_email ON users (email);

CREATE INDEX idx_session_expires_at ON sessions (expires_at);
CREATE INDEX idx_session_user_id ON sessions (user_id);

                                                       
