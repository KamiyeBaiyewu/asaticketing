CREATE TABLE IF NOT EXISTS ticket_categories (
    category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(60) NOT NULL DEFAULT '' :: character varying,
    description varchar(300) NOT NULL DEFAULT '' :: character varying,
    weight int4 NOT NULL DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ticket_categories_name ON ticket_categories USING btree (name)
WHERE
    (deleted_at IS NULL);