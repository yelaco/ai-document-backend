DROP INDEX IF EXISTS idx_documents_user_id;

ALTER TABLE documents
DROP CONSTRAINT IF EXISTS fk_documents_user;

DROP TABLE IF EXISTS documents;
