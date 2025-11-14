# ADR 0007: Backend Architecture - Onion Architecture

## Status

Accepted

## Context

In OpsCore backend development, it is necessary to introduce a consistent architectural pattern to enhance maintainability, testability, and scalability. Currently, there are no clear architectural guidelines, which could make future feature additions and changes difficult. By thoroughly separating concerns based on the principles of Clean Architecture, we aim to address these challenges.

## Decision

We will adopt the **Onion Architecture** for the backend architecture. This is based on the principle of directing dependencies unidirectionally from the inside (Domain) to the outside (Infrastructure).

**Layer Structure and Responsibilities:**

1. **Domain Layer:**
    * **Responsibility:** Represents the core business logic and rules of the application. Includes entities (created *only* via factory or reconstructor functions to ensure invariants), value objects, domain events, and repository interfaces. Ensures data integrity and enforces business rules.
    * **Dependencies:** Does not depend on external layers.
    * **Folder:** `internal/<context>/domain/`
        * `entity/`: Entities (with unexported fields), e.g., `repository.go`, `file_node.go`
        * `value_object/`: Value Objects (e.g., `repository_id.go`, `file_path.go`)
        * `domain_service/`: Domain services (logic spanning multiple entities)
        * `repository/`: Repository interface definitions
        * `error/`: Domain-specific custom errors

2. **Application Layer:**
    * **Responsibility:** Implements use cases. Orchestrates domain objects (retrieved via repositories and manipulated via their methods or domain services) and controls application-specific flows. Defines DTOs that represent the data needed for use cases, independent of API formats. Depends on infrastructure layer interfaces (e.g., repository interfaces).
    * **Dependencies:** Depends on the Domain Layer. Depends on Infrastructure Layer interfaces.
    * **Folder:** `internal/<context>/application/`
        * `usecase/`: Implementation of each use case
        * `dto/`: Data Transfer Objects for use case requests/responses (without API-specific details)
        * `error/`: Application-specific custom errors

3. **Infrastructure Layer:**
    * **Responsibility:** Implements technical details such as database access, external API integration, message queues, logging, etc. Implements interfaces defined in the Application or Domain layers (e.g., repository interfaces).
    * **Dependencies:** May depend on the Domain or Application layers (through the interfaces it implements), but not on specific business logic.
    * **Folder:** `internal/<context>/infrastructure/`
        * `persistence/`: Database-related implementations (repository implementation, ORM models, etc.)
        * `external/`: External API clients (e.g., `git/` for Git operations)
        * `messaging/`: Message queue related
        * `error/`: Infrastructure-specific custom errors

4. **Interfaces Layer (Presentation Layer):**
    * **Responsibility:** Provides interfaces to the outside world (HTTP clients, CLI, etc.). Receives requests, calls the appropriate Application Layer use case, and returns the result as a response. Handles conversion between API schemas and application DTOs.
    * **Dependencies:** Depends on the Application Layer.
    * **Folder:** `internal/<context>/interfaces/`
        * `api/handlers/`: HTTP handlers and routing configuration
        * `api/schema/`: API request/response schemas (with API-specific details like json/binding tags)
        * `grpc/`: gRPC service definitions (if applicable)
        * `cli/`: Command-line interface (if applicable)
        * `error/`: Interface-specific custom errors (e.g., HTTP error responses)

**Proposed Folder Structure:**

Given Go's package management characteristics (all files in the same directory must belong to the same package), we introduce **Bounded Context** layering between `internal/` and architectural layers (domain, application, etc.). This prevents the proliferation of files directly under layer directories and provides clear boundaries aligned with DDD principles.

```text
backend/
├── cmd/
│   └── server/
│       └── main.go       # Entry point, DI container initialization, server startup
├── internal/
│   ├── <context-name>/   # Bounded Context (e.g., git_repository, blog, user)
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   │   └── <entity>.go   # Example: repository.go, file_node.go
│   │   │   ├── value_object/     # Value Objects
│   │   │   ├── domain_service/   # Domain services
│   │   │   ├── repository/
│   │   │   │   └── <repository_interface>.go # Example: repository.go
│   │   │   └── error/            # Domain-specific custom errors
│   │   ├── application/
│   │   │   ├── usecase/
│   │   │   │   └── <usecase>.go # Example: repository_usecase.go
│   │   │   ├── dto/              # Data Transfer Objects (use case data, no API details)
│   │   │   └── error/            # Application-specific custom errors
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   └── <repository_impl>.go # Example: postgres_repository.go
│   │   │   ├── <external_system>/ # Example: git/
│   │   │   └── error/            # Infrastructure-specific custom errors
│   │   └── interfaces/
│   │       ├── api/
│   │       │   ├── handlers/
│   │       │   │   └── <handler>.go # Example: repository_handler.go
│   │       │   └── schema/       # API request/response schemas
│   │       │       └── <schema>.go # Example: repository_schema.go
│   │       └── error/            # Interface-specific custom errors
│   │
│   └── shared/           # Shared context for common utilities
│       ├── domain/
│       │   └── value_object/     # Common value objects
│       └── infrastructure/
│           └── middleware/       # Shared middleware
│
├── pkg/                  # Code potentially used outside the project (not expected to be used much this time)
├── go.mod
└── go.sum
```

**Example with concrete context:**

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── git_repository/                # Git repository management context
│   │   ├── domain/
│   │   │   ├── entity/                # Entities (formerly model/)
│   │   │   │   ├── repository.go      # package entity
│   │   │   │   ├── repository_test.go
│   │   │   │   ├── file_node.go
│   │   │   │   └── file_node_test.go
│   │   │   ├── value_object/          # Value Objects
│   │   │   ├── domain_service/        # Domain Services
│   │   │   ├── repository/
│   │   │   │   ├── repository.go      # package repository (interface)
│   │   │   │   └── mock_repository.go
│   │   │   └── error/                 # Domain-specific errors
│   │   ├── application/
│   │   │   ├── usecase/
│   │   │   │   ├── repository_usecase.go
│   │   │   │   └── repository_usecase_test.go
│   │   │   ├── dto/                # Use case DTOs (no API-specific details)
│   │   │   └── error/                 # Application-specific errors
│   │   ├── infrastructure/
│   │   │   ├── persistence/
│   │   │   │   ├── postgres_repository.go
│   │   │   │   └── postgres_repository_test.go
│   │   │   ├── git/
│   │   │   │   ├── git_manager.go
│   │   │   │   ├── cli_git_manager.go
│   │   │   │   └── github_api_manager.go
│   │   │   └── error/                 # Infrastructure-specific errors
│   │   └── interfaces/
│   │       ├── api/
│   │       │   ├── handlers/
│   │       │   │   ├── repository_handler.go
│   │       │   │   └── repository_handler_test.go
│   │       │   └── schema/           # API schemas with json/binding tags
│   │       │       ├── repository_schema.go
│   │       │       ├── file_schema.go
│   │       │       └── converter.go  # Schema ↔ DTO conversion
│   │       └── error/                 # Interface-specific errors
│   │
│   └── shared/
│       └── infrastructure/
│           └── middleware/
│               └── logger_middleware.go
│
├── pkg/
├── go.mod
└── go.sum
```

**Bounded Context Guidelines:**

* **Context Identification:** Each context represents a cohesive business capability or subdomain (e.g., repository management, blog, user management).
* **Context Independence:** Each context should be as independent as possible, with its own domain models, use cases, and infrastructure implementations.
* **Inter-Context Communication:**
  * Avoid direct circular dependencies between contexts.
  * Shared concepts should be placed in the `shared/` context.
  * Consider using domain events or anti-corruption layers for complex inter-context communication.
* **Package Naming:** Import paths will include the context name, e.g., `github.com/IT-TEPPEN/ops-core/backend/internal/git_repository/domain/entity`.
* **Naming Conventions:**
  * Use `entity/` instead of `model/` to clearly distinguish from ORM models in infrastructure layer.
  * Use descriptive context names that avoid conflicts with layer names (e.g., `git_repository` instead of `repository`).
  * Each layer should have its own `error/` package for custom error types specific to that layer.

**Dependency Rules:**

* Dependencies always point inwards, from outer layers to inner layers (Interfaces -> Application -> Domain).
* The Infrastructure layer implements interfaces defined in the Application or Domain layers, achieving Dependency Inversion.
* The Domain layer does not depend on any other layer.
* **Context-level dependencies:** Contexts should minimize dependencies on other contexts. When needed, depend on `shared/` context or use well-defined interfaces.

**Schema and DTO Separation:**

To maintain proper separation of concerns and unidirectional dependency flow, the architecture distinguishes between API schemas and application DTOs:

* **Application DTOs** (`application/dto/`):
  * Represent data structures needed for use case logic
  * Do not contain API-specific details (no json tags, binding validations, or swagger annotations)
  * Focus solely on what the application layer needs to function
  * Independent of API format changes
  
* **API Schemas** (`interfaces/api/schema/`):
  * Define the exact structure of API requests and responses
  * Include API-specific details (json tags, binding validations, swagger examples)
  * Handle API versioning and format requirements
  * Convert to/from application DTOs within the interfaces layer

* **Conversion Flow:**
  ```
  API Request → Schema (interfaces) → DTO (application) → Use Case
  Use Case → DTO (application) → Schema (interfaces) → API Response
  ```

This separation ensures:
* Application layer remains independent of API specifications
* API format changes don't require modifying use case logic
* DTOs can evolve based on business needs without breaking API contracts
* Clear unidirectional dependency: interfaces depends on application, not vice versa

**Sample Code (Conceptual):**

* **`internal/user/domain/entity/user.go`**:

    ```go
    package entity

    import "errors"

    // UserID represents the value object for User ID
    type UserID string

    // NewUserID creates a new UserID (example factory)
    func NewUserID(id string) (UserID, error) {
        if id == "" {
            return "", errors.New("user ID cannot be empty")
        }
        // Add more validation logic if needed
        return UserID(id), nil
    }

    // String returns the string representation of UserID
    func (uid UserID) String() string {
        return string(uid)
    }

    // User entity with unexported fields
    type User struct {
        id   UserID // Unexported
        name string // Unexported
        // ... other unexported fields
    }

    // NewUser is a factory function for creating *new* User entities.
    // It ensures all invariants for a new user are met.
    func NewUser(id UserID, name string /* ... other params */) (*User, error) {
        if name == "" {
            return nil, errors.New("user name cannot be empty")
        }
        // Add other validation logic for creating a *new* user
        return &User{
            id:   id,
            name: name,
            // ... initialize other fields
        }, nil
    }

    // ReconstructUser is used to reconstruct a User entity from persistence.
    // It bypasses some validation that only applies to *new* entities,
    // but still ensures the data represents a valid state.
    func ReconstructUser(id UserID, name string /* ... other fields */) *User {
        // Minimal validation might still be needed, or trust the data source
        return &User{
            id:   id,
            name: name,
            // ... initialize other fields
        }
    }

    // Getter methods to access unexported fields
    func (u *User) ID() UserID {
        return u.id
    }

    func (u *User) Name() string {
        return u.name
    }

    // Methods on User to encapsulate behavior...
    // e.g., func (u *User) ChangeName(newName string) error { ... }
    ```

* **`internal/user/domain/repository/user_repository.go`**:

    ```go
    package repository

    import (
        "context"
        "YOUR_PROJECT/internal/user/domain/entity"
    )

    type UserRepository interface {
        FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) // Use Value Object
        Save(ctx context.Context, user *entity.User) error
    }
    ```

* **`internal/user/application/usecase/user_usecase.go`**:

    ```go
    package usecase

    import (
        "context"
        "YOUR_PROJECT/internal/user/domain/entity"
        "YOUR_PROJECT/internal/user/domain/repository"
        "YOUR_PROJECT/internal/user/application/dto"
        "errors" // Example for error handling
    )

    type UserUsecase struct {
        userRepo repository.UserRepository
    }

    func NewUserUsecase(ur repository.UserRepository) *UserUsecase {
        return &UserUsecase{userRepo: ur}
    }

    func (uc *UserUsecase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
        // Generate Value Objects first if applicable
        // In a real scenario, ID might be generated by the DB or a UUID generator
        userID, err := entity.NewUserID("some-generated-id") // Example ID generation
        if err != nil {
             return nil, errors.New("invalid user id format") // Or a more specific error type
        }

        // Use the factory function to create the entity
        user, err := entity.NewUser(userID, req.Name /* ... other fields from req ... */)
        if err != nil {
            // Handle domain validation errors (e.g., invalid name)
            return nil, err // Propagate domain error
        }

        if err := uc.userRepo.Save(ctx, user); err != nil {
            // Handle infrastructure errors (e.g., database connection issue)
            return nil, errors.New("failed to save user") // Or a more specific error type
        }

        // Map entity back to DTO using getter methods
        res := &dto.UserResponse{
            ID:   user.ID().String(), // Use getter and convert UserID
            Name: user.Name(),        // Use getter
            // ... map other fields using getters ...
        }
        return res, nil
    }

    func (uc *UserUsecase) GetUser(ctx context.Context, userIDStr string) (*dto.UserResponse, error) {
        userID, err := entity.NewUserID(userIDStr)
        if err != nil {
             return nil, errors.New("invalid user id format")
        }
        user, err := uc.userRepo.FindByID(ctx, userID)
        if err != nil {
            // Handle not found or other repo errors
            return nil, err
        }
        if user == nil {
             return nil, errors.New("user not found") // Or specific error type
        }

        // Map entity back to DTO using getter methods
        res := &dto.UserResponse{
            ID:   user.ID().String(), // Use getter and convert UserID
            Name: user.Name(),        // Use getter
             // ... map other fields using getters ...
        }
        return res, nil
    }
    ```

* **`internal/user/infrastructure/persistence/user_repository_impl.go`**:

    ```go
    package persistence

    import (
        "context"
        "errors"
        "gorm.io/gorm"
        "YOUR_PROJECT/internal/user/domain/entity"
        "YOUR_PROJECT/internal/user/domain/repository"
    )

    // GormUserModel represents the data structure in the database (ORM model)
    type GormUserModel struct {
     // gorm.Model // Uncomment if using standard GORM fields like ID, CreatedAt etc.
     UserID string `gorm:"primaryKey"` // Store UserID as string
     Name   string
     // ... other fields matching DB columns
    }

    func (GormUserModel) TableName() string {
        return "users" // Explicitly set table name
    }

    type UserRepositoryImpl struct {
        db *gorm.DB
    }

    func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
        // AutoMigrate should ideally be handled by a separate migration tool/process
        // db.AutoMigrate(&GormUserModel{}) // Example migration
        return &UserRepositoryImpl{db: db}
    }

    // toDomain converts Gorm ORM model to domain entity using the reconstructor
    func toDomain(gormUser *GormUserModel) (*entity.User, error) {
        if gormUser == nil {
            return nil, errors.New("cannot convert nil gorm user model to domain entity")
        }
        userID, err := entity.NewUserID(gormUser.UserID)
        if err != nil {
            // This indicates data integrity issue in the DB or mapping
            return nil, errors.New("invalid user ID format in database: " + err.Error())
        }

        // Use the Reconstructor function
        user := entity.ReconstructUser(userID, gormUser.Name /* ... other fields */)
        return user, nil
    }

    // fromDomain converts domain entity to Gorm ORM model for persistence
    func fromDomain(user *entity.User) (*GormUserModel, error) {
         if user == nil {
             return nil, errors.New("cannot convert nil domain user to gorm model")
         }
         return &GormUserModel{
             UserID: user.ID().String(), // Use getter and convert UserID
             Name:   user.Name(),        // Use getter
             // ... map other fields using getters ...
         }, nil
    }

    func (r *UserRepositoryImpl) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
        var userModel GormUserModel
        // Search logic using GORM with context, converting UserID to string for query
        err := r.db.WithContext(ctx).First(&userModel, "user_id = ?", id.String()).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil, nil // Return nil, nil for not found (idiomatic in Go)
            }
            return nil, err // Return other DB errors
        }
        return toDomain(&userModel)
    }

    func (r *UserRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
        // Convert domain entity to persistence ORM model
        userModel, err := fromDomain(user)
        if err != nil {
            return err // Handle conversion error
        }

        // Save logic using GORM with context (Create or Update)
        // GORM's Save handles upsert based on primary key presence
        // Ensure the GormUserModel has the correct primary key tag (`gorm:"primaryKey"`)
        return r.db.WithContext(ctx).Save(userModel).Error
    }
    ```

* **`internal/user/interfaces/api/schema/user_schema.go`**:

    ```go
    package schema

    import "YOUR_PROJECT/internal/user/application/dto"

    // CreateUserRequest represents the API request for creating a user
    type CreateUserRequest struct {
        Name string `json:"name" binding:"required" example:"John Doe"`
        // ... other fields with API-specific tags
    }

    // UserResponse represents the API response for a user
    type UserResponse struct {
        ID   string `json:"id" example:"user-123"`
        Name string `json:"name" example:"John Doe"`
        // ... other fields with API-specific tags
    }

    // ToCreateUserDTO converts API schema to application DTO
    func ToCreateUserDTO(req CreateUserRequest) dto.CreateUserRequest {
        return dto.CreateUserRequest{
            Name: req.Name,
            // ... map other fields
        }
    }

    // FromUserDTO converts application DTO to API schema
    func FromUserDTO(dtoResp dto.UserResponse) UserResponse {
        return UserResponse{
            ID:   dtoResp.ID,
            Name: dtoResp.Name,
            // ... map other fields
        }
    }
    ```

* **`internal/user/interfaces/api/handlers/user_handler.go`**:

    ```go
    package handlers

    import (
        "github.com/gin-gonic/gin"
        "YOUR_PROJECT/internal/user/application/usecase"
        "YOUR_PROJECT/internal/user/interfaces/api/schema"
        "net/http"
        "errors" // For example error checking
    )

    type UserHandler struct {
        userUsecase *usecase.UserUsecase
    }

    func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
        return &UserHandler{userUsecase: uc}
    }

    func (h *UserHandler) CreateUser(c *gin.Context) {
        var req schema.CreateUserRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
            return
        }

        // Convert schema to DTO
        dtoReq := schema.ToCreateUserDTO(req)

        // Pass request context to use case
        dtoResp, err := h.userUsecase.CreateUser(c.Request.Context(), &dtoReq)
        if err != nil {
            // Basic error handling - map domain/app errors to HTTP status codes
            // This could be more sophisticated (e.g., checking error types)
            if errors.Is(err, /* some specific domain validation error type */ nil) {
                 c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            } else {
                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            }
            return
        }

        // Convert DTO to schema
        response := schema.FromUserDTO(*dtoResp)
        c.JSON(http.StatusCreated, response)
    }

    func (h *UserHandler) GetUser(c *gin.Context) {
        userID := c.Param("id") // Assuming ID is a path parameter like /users/:id

        // Pass request context to use case
        dtoResp, err := h.userUsecase.GetUser(c.Request.Context(), userID)
         if err != nil {
            // Handle errors like invalid ID format or user not found
             if errors.Is(err, /* specific not found error type */ nil) {
                 c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
             } else if errors.Is(err, /* specific invalid format error type */ nil) {
                 c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
             } else {
                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
             }
            return
        }

        // Convert DTO to schema
        response := schema.FromUserDTO(*dtoResp)
        c.JSON(http.StatusOK, response)
    }
    ```

* **`cmd/server/main.go`**:

    ```go
    package main

    import (
        "github.com/gin-gonic/gin"
        "gorm.io/driver/postgres" // Example
        "gorm.io/gorm"
        "YOUR_PROJECT/internal/user/domain/repository"
        "YOUR_PROJECT/internal/user/application/usecase"
        "YOUR_PROJECT/internal/user/infrastructure/persistence"
        "YOUR_PROJECT/internal/user/interfaces/api/handlers"
        // ... router configuration, etc.
    )

    func main() {
        // Initialize database connection (Example)
        dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Tokyo"
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            panic("failed to connect database")
        }

        // Dependency Injection (DI)
        userRepo := persistence.NewUserRepositoryImpl(db)
        userUsecase := usecase.NewUserUsecase(userRepo)
        userHandler := handlers.NewUserHandler(userUsecase)

        // Router configuration (Gin example)
        r := gin.Default()
        // ... Routing configuration (e.g., r.POST("/users", userHandler.CreateUser))

        r.Run() // Start server
    }
    ```

## Consequences

* **Pros:**
  * **High Testability:** Each layer, especially the Domain and Application layers, is independent of infrastructure details, making unit testing easier.
  * **Separation of Concerns:** Business logic, application logic, and infrastructure are clearly separated, making code easier to understand and maintain.
  * **Interchangeability:** Infrastructure implementations (DB, external API clients, etc.) can be changed with minimal impact on other layers.
  * **Scalability:** Adding new use cases or features is easier with reduced impact on existing code.
* **Cons:**
  * **Learning Curve:** Understanding the concepts of Onion Architecture and its dependency rules is required.
  * **Initial Development Complexity:** The initial amount of code and configuration might increase due to interface definitions between layers and DTO creation.
  * **Boilerplate Code:** Repetitive code might increase due to data mapping between layers (can be mitigated with code generation tools).
