CREATE TABLE IF NOT EXISTS kb_documents (
    doc_id TEXT PRIMARY KEY,
    kb_id TEXT NOT NULL,
    title TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT '',
    metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

        CONSTRAINT fk_kb_documents_kb
        FOREIGN KEY (kb_id)
        REFERENCES knowledge_bases(kb_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_kb_documents_kb_id ON kb_documents(kb_id);