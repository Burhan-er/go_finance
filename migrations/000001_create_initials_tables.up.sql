-- 000001_create_initial_tables.up.sql

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE balances (
    user_id UUID PRIMARY KEY,
    amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID, 
    action VARCHAR(50) NOT NULL,
    details VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS balance_history (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION fn_balance_history_trigger()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO balance_history(user_id, amount, recorded_at)
    VALUES (NEW.user_id, NEW.amount, NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_balance_update
AFTER UPDATE ON balances
FOR EACH ROW
EXECUTE FUNCTION fn_balance_history_trigger();

CREATE TRIGGER trg_balance_insert
AFTER INSERT ON balances
FOR EACH ROW
EXECUTE FUNCTION fn_balance_history_trigger();

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language plpgsql;

CREATE OR REPLACE FUNCTION update_last_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated_at = NOW();
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_balances_last_updated_at
BEFORE UPDATE ON balances
FOR EACH ROW
EXECUTE FUNCTION update_last_updated_at_column();

INSERT INTO users (username, email, password_hash, role)
VALUES 
    ('user1', 'ahmet@example.com', '$2a$14$2E2JqjBF3MG5omtKKH1C5OP5crlkjfu2DkI9SHAOuozaT/AGZ6RWC', 'user'),
    ('user2', 'tuncay@example.com', '$2a$14$lSgFh/GeTSOrAinC5y6E4elS9N98zKLRcXfigbo4MH1I622ts/iy2', 'user'),
    ('admin', 'admin@example.com', '$2a$14$2i/8ilIoeLLE/ron26DQZOQm5CU9YHMAD8pMcDbSsc7bdWhJuwpAq', 'admin');

INSERT INTO balances (user_id, amount)
SELECT id, 1000.00 FROM users WHERE username = 'user1';
INSERT INTO balances (user_id, amount)
SELECT id, 500.00 FROM users WHERE username = 'user2';
INSERT INTO balances (user_id, amount)
SELECT id, 999999.99 FROM users WHERE username = 'admin';
