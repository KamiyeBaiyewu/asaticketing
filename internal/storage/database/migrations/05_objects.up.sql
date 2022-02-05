CREATE TABLE IF NOT EXISTS objects(
    object_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    is_standard bool NOT NULL DEFAULT FALSE,
    created_by UUID REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX objects_name ON objects USING btree (name)
WHERE (deleted_at IS NULL);
   