 
CREATE TABLE IF NOT EXISTS users_roles(
	user_id UUID REFERENCES users,
	role_id UUID REFERENCES roles,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- craete INDEX FOR roles
CREATE UNIQUE INDEX IF NOT EXISTS users_roles_user ON users_roles (user_id,role_id) WHERE deleted_at IS NULL;