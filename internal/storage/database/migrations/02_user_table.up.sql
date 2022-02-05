CREATE TABLE IF NOT EXISTS users(
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    firstname text NOT NULL,
    lastname text NOT NULL,
    email text NOT NULL,
    password_hash bytea,
    user_type text NOT NULL,
    role_id uuid REFERENCES roles,
    is_active bool NOT NULL DEFAULT TRUE,
    is_system bool NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX user_email ON users (email) WHERE deleted_at IS NULL;