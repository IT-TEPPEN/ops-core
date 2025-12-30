-- Drop tables in reverse order
-- Drop documents first because it has foreign key to document_versions
DROP TABLE IF EXISTS documents CASCADE;
DROP TABLE IF EXISTS document_versions CASCADE;
