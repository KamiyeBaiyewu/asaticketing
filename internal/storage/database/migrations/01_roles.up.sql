CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles(
	role_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


CREATE UNIQUE INDEX unique_role ON roles (name) WHERE deleted_at IS NULL;
