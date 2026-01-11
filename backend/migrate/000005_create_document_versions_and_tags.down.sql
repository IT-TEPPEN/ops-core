-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS document_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS document_current_version CASCADE;
DROP TABLE IF EXISTS document_versions CASCADE;
