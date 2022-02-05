CREATE TABLE contacts(
    contact_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    firstname TEXT,
    lastname TEXT,
    phone_no text not null,
    email text NOT NULL,
    created_by uuid REFERENCES users,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL

);

create UNIQUE index unique_contact on contacts (email) 
where deleted_at is null;