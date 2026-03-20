#!/usr/bin/env bash
set -e

psql "${POSTGRES_DSN}" -f migrations/000001_init_sessions.sql
psql "${POSTGRES_DSN}" -f migrations/000002_init_session_messages.sql
psql "${POSTGRES_DSN}" -f migrations/000003_init_conversation_state.sql
psql "${POSTGRES_DSN}" -f migrations/000004_init_knowledge_bases.sql
psql "${POSTGRES_DSN}" -f migrations/000005_init_documents.sql
psql "${POSTGRES_DSN}" -f migrations/000006_init_document_chunks.sql