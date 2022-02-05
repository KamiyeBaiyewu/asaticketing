

DROP TYPE IF EXISTS object_action;
CREATE TYPE object_action AS ENUM (
'create',
'view',
'list',
'update',
'delete',
'import',
'export'
);

CREATE TABLE IF NOT EXISTS object_policies(
    policy_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID REFERENCES roles,
    object_id UUID REFERENCES objects,
    action object_action NOT NULL,
    created_by UUID REFERENCES users,
    is_standard bool NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX object_policies_unique ON object_policies USING btree (role_id,object_id,action)
WHERE (deleted_at IS NULL);