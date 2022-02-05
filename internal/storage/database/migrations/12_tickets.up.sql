CREATE SEQUENCE IF NOT EXISTS ticket_number_seq;

SELECT
    setval('ticket_number_seq', 10000);

CREATE TABLE IF NOT EXISTS tickets(
    ticket_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subject VARCHAR(150) NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    number int  DEFAULT nextval('ticket_number_seq'),
    created_by UUID REFERENCES users,
    category_id UUID REFERENCES ticket_categories,
    status_id UUID REFERENCES ticket_statuses,
    priority_id UUID REFERENCES ticket_priorities,
    source_id UUID REFERENCES ticket_sources,
    sla_id UUID REFERENCES ticket_slas,
    assigned_to UUID REFERENCES users,
    deadline TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);