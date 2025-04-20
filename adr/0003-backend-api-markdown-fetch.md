# ADR 0003: Backend API for Fetching External Markdown

## Status

Accepted

## Context

Based on ADR 0001 and ADR 0002, OpsCore needs to fetch operational procedure Markdown documents from external private GitLab and GitHub repositories. To achieve this functionality, a clear API specification is required for the frontend to communicate with the backend. This API will be responsible for retrieving the content of specified Markdown files from configured external repositories.

## Decision

We define the backend API endpoint for fetching Markdown files from external repositories as follows:

**Endpoint:**

```
GET /api/v1/repositories/{config_id}/file
```

**Description:**

Retrieves the content of a Markdown file at the specified path using an external repository configuration (`identified by config_id`) previously set up and stored in the OpsCore backend.

**Path Parameters:**

*   `config_id` (string, required): The unique identifier for the target repository configuration. This configuration must securely store the repository URL, provider type (GitHub/GitLab), and authentication credentials (such as access tokens or application authentication details) as defined in ADR 0002.

**Query Parameters:**

*   `path` (string, required): The path to the Markdown file from the repository root (e.g., `docs/procedures/backup.md`).
*   `ref` (string, optional): The Git reference (branch name, tag name, commit hash) from which to fetch the file. If not specified, the repository's default branch will be used.

**Authentication:**

*   The API endpoint itself must be protected by OpsCore's standard authentication mechanism (e.g., JWT token, session).
*   The backend service will use the credentials associated with the `config_id` to authenticate against the external Git repository.

**Responses:**

*   **Success (200 OK):**
    *   `Content-Type: text/markdown; charset=utf-8`
    *   Body: The raw content of the requested Markdown file.

*   **Error:**
    *   `400 Bad Request`: If required parameters (`path`) are missing.
    *   `401 Unauthorized / 403 Forbidden`: If the OpsCore backend does not have access rights to the external repository (e.g., invalid credentials, insufficient permissions). This also includes authentication failure for the API endpoint itself.
    *   `404 Not Found`: If the repository configuration for the specified `config_id` does not exist, or if the file at the specified `path` or `ref` is not found within the repository.
    *   `500 Internal Server Error`: If an unexpected error occurs within the backend (e.g., failure executing Git commands, communication errors with external APIs).

**Related Features (Future Considerations):**

*   An endpoint to list files/directories within a directory (`GET /api/v1/repositories/{config_id}/tree?path={dir_path}&ref={ref}`) might also be useful in the future.

**Prerequisites:**

*   Backend functionality and UI for securely managing (creating, updating, deleting) external repository connection information (URL, provider type, credentials) are required separately. This API assumes that such information is accessible via the `config_id`.

## Consequences

*   **Pros:**
    *   Clearly separates concerns between the frontend and backend. The frontend doesn't need to be aware of the complexities of file fetching.
    *   Follows RESTful design principles, providing an understandable and easy-to-use interface.
    *   Enhances security by centralizing credential management in the backend.
*   **Cons:**
    *   Requires implementing Git operation logic and authentication handling with external repositories in the backend.
    *   Secure storage and management of credentials associated with `config_id` are crucial.
    *   Proper error handling (e.g., connection failures to external repositories, permission errors) needs to be implemented and communicated clearly to the frontend.
