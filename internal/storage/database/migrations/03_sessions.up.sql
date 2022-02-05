 CREATE TABLE IF NOT EXISTS sessions(
	user_id UUID REFERENCES users,
	device_id text, 
	refresh_token text, 
	expires_at integer,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
	PRIMARY KEY (user_id, device_id)
);