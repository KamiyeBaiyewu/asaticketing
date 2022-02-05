CREATE TABLE  IF NOT EXISTS ticket_slas (
    agreement_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(30) NOT NULL DEFAULT '' :: character varying,
    grace_period int4 NOT NULL DEFAULT 24,
    weight int4 NOT NULL DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


CREATE UNIQUE INDEX ticket_slas_name ON ticket_slas USING btree (name)
WHERE
    (deleted_at IS NULL);