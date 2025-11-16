# Error Code Reference

This document provides a comprehensive reference for all error codes used in the OpsCore backend application.

## Error Code Format

Error codes follow the format: `<LAYER>_<CATEGORY>_<NUMBER>`

**Layer Prefixes:**
- `DOM`: Domain Layer
- `APP`: Application Layer
- `INF`: Infrastructure Layer

**Severity Levels:**
- `LOW`: Minor issues that don't significantly impact functionality
- `MEDIUM`: Issues that may impact some functionality
- `HIGH`: Serious issues that impact core functionality
- `CRITICAL`: Critical failures that prevent system operation

## Error Codes

| Code | Layer | Category | Description | Severity | Retryable |
|------|-------|----------|-------------|----------|-----------|
| GRA0001 | Application | Resource | Requested resource not found | LOW | No |
| GRA0002 | Application | Resource | Resource conflict (duplicate) | MEDIUM | No |
| GRA0003 | Application | Authentication | Unauthorized access | HIGH | No |
| GRA0004 | Application | Authorization | Forbidden access | HIGH | No |
| GRA0005 | Application | Authentication | Invalid credentials | MEDIUM | No |
| GRA0006 | Application | Validation | Application-level validation failed | MEDIUM | No |
| GRA0007 | Application | Validation | Invalid request format | MEDIUM | No |
| GRD0001 | Domain | Validation | Invalid entity field value | MEDIUM | No |
| GRD0002 | Domain | Validation | Required field is missing | MEDIUM | No |
| GRD0003 | Domain | Validation | Invalid field format | MEDIUM | No |
| GRD0004 | Domain | Validation | Field value out of range | MEDIUM | No |
| GRD0005 | Domain | Validation | Invalid URL format | MEDIUM | No |
| GRD0006 | Domain | Validation | Unsupported URL scheme (only HTTPS is supported) | MEDIUM | No |
| GRD0007 | Domain | Business | Business rule violation | HIGH | No |
| GRD0008 | Domain | Business | Invalid state transition | HIGH | No |
| GRD0009 | Domain | Business | Invariant violation | HIGH | No |
| GRI0001 | Infrastructure | Database | Database connection error | CRITICAL | Yes |
| GRI0002 | Infrastructure | Database | Database query error | HIGH | No |
| GRI0003 | Infrastructure | Database | Database constraint violation | MEDIUM | No |
| GRI0004 | Infrastructure | Database | Database query timeout | HIGH | Yes |
| GRI0005 | Infrastructure | ExternalAPI | External API error | HIGH | No |
| GRI0006 | Infrastructure | ExternalAPI | External API timeout | HIGH | Yes |
| GRI0007 | Infrastructure | ExternalAPI | External API resource not found | MEDIUM | No |
| GRI0008 | Infrastructure | Connection | Connection failed | CRITICAL | Yes |
| GRI0009 | Infrastructure | Connection | Connection timeout | HIGH | Yes |
| GRI0010 | Infrastructure | Storage | Storage operation failed | HIGH | No |
| GRI0011 | Infrastructure | Storage | Storage resource not found | MEDIUM | No |

## Usage

### In Code

```go
// Domain Layer
return &error.ValidationError{
    Code:    error.CodeInvalidURL,
    Field:   "url",
    Message: "must be a valid HTTPS URL",
}

// Application Layer
return &error.NotFoundError{
    Code:         error.CodeResourceNotFound,
    ResourceType: "Repository",
    ResourceID:   id,
}
```

### In Logs

Error codes are automatically included in log messages:

```
[DOM_VAL_005] validation failed for field 'url': must be a valid HTTPS URL
```

### In HTTP Responses

Error codes are mapped to appropriate HTTP status codes and included in responses:

```json
{
  "code": "APP_RES_001",
  "message": "Resource not found",
  "details": {
    "resource_type": "Repository",
    "resource_id": "repo-123"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2025-11-15T00:10:18Z"
}
```

## Related ADRs

- [ADR 0007: Backend Architecture - Onion Architecture](../adr/0007-backend-architecture-onion.md)
- [ADR 0008: Backend Logging Strategy](../adr/0008-backend-logging-strategy.md)
- [ADR 0015: Backend Custom Error Design](../adr/0015-backend-custom-error-design.md)
