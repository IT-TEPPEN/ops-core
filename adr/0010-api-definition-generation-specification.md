# ADR 0010: API Definition Generation and Specification

## Status

Accepted

## Context

To ensure consistency, maintainability, and ease of use for the backend API (initially defined in ADR 0003), we need to establish clear guidelines for how the API definition is created, stored, and specified. Key considerations include the location of the definition file, the method for generating it (manual vs. automated), and the specification format to use. We discussed storing the definition within the source code structure (e.g., `backend/api/`) versus a dedicated documentation folder or the root, using automated tools like `swaggo/swag` versus manual creation, and adopting standard formats like OpenAPI.

## Decision

We have decided on the following approach for managing the backend API definition:

1. **Storage Location:** The generated API definition file (e.g., `openapi.yaml` or `swagger.json`) will be stored in a dedicated directory: `backend/docs/`. This clearly separates generated documentation artifacts from source code.

2. **Generation Method:** We will use the `swaggo/swag` tool (`github.com/swaggo/swag`) to automatically generate the API definition from annotations written directly in the Go source code. This approach keeps the documentation closely tied to the implementation, reducing the risk of drift.

3. **Specification Format:** We will use the OpenAPI Specification (OAS) version 3.x. This is the industry standard, providing wide compatibility with various development and documentation tools.

## Consequences

* **Pros:**
  * Generated definition files are clearly separated in `backend/docs/`.
  * Automation via `swaggo/swag` ensures the API definition stays synchronized with the Go code implementation with minimal manual effort.
  * Adherence to OpenAPI 3.x standard facilitates the use of standard tooling for documentation UI, client generation, testing, etc.
  * `swaggo/swag` is a well-established and widely used tool in the Go ecosystem.

* **Cons:**
  * Developers need to learn the specific annotation syntax required by `swaggo/swag`.
  * The build or development process must incorporate the `swag init` command (or equivalent) to generate/update the definition file.
  * Initial setup and integration of `swaggo/swag` into the project require some effort.
