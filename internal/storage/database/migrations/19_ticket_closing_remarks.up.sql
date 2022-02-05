CREATE TABLE IF NOT EXISTS ticket_closing_remarks(
    remark_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID REFERENCES tickets,
    cause_id UUID REFERENCES ticket_causes,
    closed_by UUID REFERENCES users,
    remark text ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ticket_closing_remarks_unique ON ticket_closing_remarks USING btree (ticket_id)
WHERE
    (deleted_at IS NULL);