CREATE TABLE IF NOT EXISTS ticket_causes(
    cause_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(30) NOT NULL DEFAULT '',
    description varchar(300) NOT NULL DEFAULT '',
    weight int4 NOT NULL DEFAULT 10,
    created_by UUID REFERENCES users,
    is_standard bool NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ticket_causes_name ON ticket_causes USING btree (name)
WHERE
    (deleted_at IS NULL);