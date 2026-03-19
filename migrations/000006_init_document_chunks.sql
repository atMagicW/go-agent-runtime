CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS kb_document_chunks (
    chunk_id TEXT PRIMARY KEY,
    doc_id TEXT NOT NULL,
    kb_id TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    embedding vector(1536),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_kb_document_chunks_doc
        FOREIGN KEY (doc_id)
        REFERENCES kb_documents(doc_id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_kb_document_chunks_kb
        FOREIGN KEY (kb_id)
        REFERENCES knowledge_bases(kb_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_kb_document_chunks_kb_id ON kb_document_chunks(kb_id);