
CREATE TABLE IF NOT EXISTS system_object_policies(
    policy_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID REFERENCES roles,
    object_id UUID REFERENCES system_objects,
    created_by UUID REFERENCES users,
    is_standard bool NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX system_object_policies_unique ON system_object_policies USING btree (role_id,object_id)
WHERE (deleted_at IS NULL);