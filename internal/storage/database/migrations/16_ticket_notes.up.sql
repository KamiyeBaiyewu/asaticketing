CREATE TABLE IF NOT EXISTS ticket_notes(
    note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    note text NOT NULL,
    created_by UUID REFERENCES users,
    ticket_id UUID REFERENCES tickets,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

