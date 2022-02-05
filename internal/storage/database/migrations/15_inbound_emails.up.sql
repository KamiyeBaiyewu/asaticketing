CREATE TABLE IF NOT EXISTS inbound_emails(
    email_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL DEFAULT '',
    status VARCHAR(500) NOT NULL DEFAULT '',
    address VARCHAR(100) NOT NULL DEFAULT '',
    email_user VARCHAR(100) NOT NULL DEFAULT '',
    email_secret VARCHAR(100) NOT NULL DEFAULT '',
    port INT NOT NULL DEFAULT 993,
    secured bool NOT NULL DEFAULT TRUE,
    mailbox VARCHAR(30) NOT NULL DEFAULT '',
    is_primary bool NOT NULL DEFAULT FALSE,
    last_seq INT NOT NULL DEFAULT 0,
	poll_period int4 NOT NULL DEFAULT 5,
    last_synced TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users,
    delete_seen bool NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX inbound_emails_unique ON inbound_emails USING btree (address, email_user, email_secret, port, mailbox)
WHERE
    (deleted_at IS NULL);