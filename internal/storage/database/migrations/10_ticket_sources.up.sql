CREATE TABLE  IF NOT EXISTS ticket_sources (
    source_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(30) NOT NULL DEFAULT '' :: character varying,
    description varchar(300) NOT NULL DEFAULT '' :: character varying,
    weight int4 NOT NULL DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


CREATE UNIQUE INDEX ticket_sources_name ON ticket_sources USING btree (name)
WHERE
    (deleted_at IS NULL);