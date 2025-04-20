# ADR 0005: Database Schema for Repository Configurations

## Status

Accepted

## Context

ADR 0004 established PostgreSQL as the chosen database for OpsCore. This ADR defines the initial database schema required to store the external repository configurations, as outlined in ADR 0002 and needed for the API specified in ADR 0003.

## Decision

We define the following table schema for storing repository configurations in the PostgreSQL database:

*   **`repository_configurations` Table:** Stores details about configured external repositories.
    *   `id`: UUID (Primary Key) - Unique identifier for the configuration. Generated automatically (e.g., using `gen_random_uuid()`).
    *   `name`: VARCHAR(255) (Unique, Not Null) - User-defined name for easy identification.
    *   `provider_type`: VARCHAR(50) (Not Null) - Type of the Git provider (e.g., 'github', 'gitlab'). Should likely be constrained to allowed values.
    *   `repository_url`: VARCHAR(2048) (Not Null) - Base URL of the repository (e.g., `https://github.com/owner/repo`, `https://gitlab.com/group/project`). Validation should ensure it's a valid URL format.
    *   `credentials`: TEXT (Not Null) - **Encrypted** access credentials or a reference identifier to retrieve credentials from a secure vault/secrets manager. **Security Critical:** Raw credentials must *never* be stored directly in this field. The encryption method or vault reference strategy needs separate consideration.
    *   `created_at`: TIMESTAMP WITH TIME ZONE (Not Null, Default: `CURRENT_TIMESTAMP`) - Timestamp of when the record was created.
    *   `updated_at`: TIMESTAMP WITH TIME ZONE (Not Null, Default: `CURRENT_TIMESTAMP`) - Timestamp of the last update. A trigger should update this automatically on row modification.

**Indexing:**
*   A unique index should be created on the `name` column.
*   An index might be beneficial on `provider_type` if filtering by provider is common.

**Future Considerations:**
*   Additional tables for users, audit logs, etc., will be defined in separate ADRs as needed.
*   Schema migration tooling (e.g., Flyway, sql-migrate) should be adopted to manage schema changes over time.

## Consequences

*   **Pros:**
    *   Provides a clear, structured definition for storing essential repository configuration data.
    *   Defines data types and constraints, promoting data integrity.
    *   Highlights the critical security requirement for handling credentials.
*   **Cons:**
    *   This schema is specific to repository configurations; other application data will require additional schema definitions.
    *   The chosen encryption strategy for the `credentials` field needs careful implementation and management.
    *   Requires implementation of schema migration management.
