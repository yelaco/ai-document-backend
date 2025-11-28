CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_name TEXT NOT NULL,
	save_path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_documents_user
    FOREIGN KEY (user_id) REFERENCES users (id)
    ON DELETE CASCADE
);

-- Index to speed up queries by user_id
CREATE INDEX idx_documents_user_id ON documents (user_id);
