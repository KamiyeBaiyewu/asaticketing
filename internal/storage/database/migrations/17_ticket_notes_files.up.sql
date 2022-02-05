CREATE TABLE IF NOT EXISTS ticket_notes_files(
    file_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename text NOT NULL,
    mime_type VARCHAR(255) NOT NULL DEFAULT '',
    size INT NOT NULL DEFAULT 0,
    note_id UUID REFERENCES ticket_notes,
    created_by UUID REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

