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
    * **Folder:** `internal/domain/`
        * `model/`: Entities (with unexported fields), Value Objects (e.g., `user.go`, `user_id.go`)
        * `repository/`: Repository interface definitions
        * `service/`: Domain services (logic spanning multiple entities)

2. **Application Layer:**
    * **Responsibility:** Implements use cases. Orchestrates domain objects (retrieved via repositories and manipulated via their methods or domain services) and controls application-specific flows. Depends on infrastructure layer interfaces (e.g., repository interfaces).
    * **Dependencies:** Depends on the Domain Layer. Depends on Infrastructure Layer interfaces.
    * **Folder:** `internal/application/`
        * `usecase/`: Implementation of each use case
        * `dto/`: Data Transfer Objects (Request/Response)

3. **Infrastructure Layer:**
    * **Responsibility:** Implements technical details such as database access, external API integration, message queues, logging, etc. Implements interfaces defined in the Application or Domain layers (e.g., repository interfaces).
    * **Dependencies:** May depend on the Domain or Application layers (through the interfaces it implements), but not on specific business logic.
    * **Folder:** `internal/infrastructure/`
        * `persistence/`: Database-related implementations (repository implementation, ORM settings, etc.)
        * `external/`: External API clients
        * `messaging/`: Message queue related
        * `logging/`: Logging implementation

4. **Interfaces Layer (Presentation Layer):**
    * **Responsibility:** Provides interfaces to the outside world (HTTP clients, CLI, etc.). Receives requests, calls the appropriate Application Layer use case, and returns the result as a response. Also handles DTO conversion.
    * **Dependencies:** Depends on the Application Layer.
    * **Folder:** `internal/interfaces/` (or `internal/handler/`, `internal/delivery/`)
        * `http/`: HTTP handlers, routing configuration
        * `grpc/`: gRPC service definitions
        * `cli/`: Command-line interface

**Proposed Folder Structure:**

```text
backend/
├── cmd/
│   └── server/
│       └── main.go       # Entry point, DI container initialization, server startup
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   └── user.go   # Example: User entity
│   │   └── repository/
│   │       └── user_repository.go # Example: UserRepository interface
│   ├── application/
│   │   ├── usecase/
│   │   │   └── user_usecase.go # Example: Create user use case
│   │   └── dto/
│   │       └── user_dto.go     # Example: Create User request DTO
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── user_repository_impl.go # Example: GORM implementation of UserRepository
│   │   └── logging/
│   │       └── logger.go
│   └── interfaces/
│       └── http/
│           ├── handler/
│           │   └── user_handler.go # Example: User related HTTP handler
│           └── router.go       # Example: HTTP router configuration (e.g., Gin)
├── pkg/                  # Code potentially used outside the project (not expected to be used much this time)
├── go.mod
└── go.sum
```

**Dependency Rules:**

* Dependencies always point inwards, from outer layers to inner layers (Interfaces -> Application -> Domain).
* The Infrastructure layer implements interfaces defined in the Application or Domain layers, achieving Dependency Inversion.
* The Domain layer does not depend on any other layer.

**Sample Code (Conceptual):**

* **`internal/domain/model/user.go`**:

    ```go
    package model

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

* **`internal/domain/repository/user_repository.go`**:

    ```go
    package repository

    import "YOUR_PROJECT/internal/domain/model"

    type UserRepository interface {
        FindByID(id model.UserID) (*model.User, error) // Use Value Object
        Save(user *model.User) error
    }
    ```

* **`internal/application/usecase/user_usecase.go`**:

    ```go
    package usecase

    import (
        "YOUR_PROJECT/internal/domain/model"
        "YOUR_PROJECT/internal/domain/repository"
        "YOUR_PROJECT/internal/application/dto"
        "errors" // Example for error handling
    )

    type UserUsecase struct {
        userRepo repository.UserRepository
    }

    func NewUserUsecase(ur repository.UserRepository) *UserUsecase {
        return &UserUsecase{userRepo: ur}
    }

    func (uc *UserUsecase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
        // Generate Value Objects first if applicable
        // In a real scenario, ID might be generated by the DB or a UUID generator
        userID, err := model.NewUserID("some-generated-id") // Example ID generation
        if err != nil {
             return nil, errors.New("invalid user id format") // Or a more specific error type
        }

        // Use the factory function to create the entity
        user, err := model.NewUser(userID, req.Name /* ... other fields from req ... */)
        if err != nil {
            // Handle domain validation errors (e.g., invalid name)
            return nil, err // Propagate domain error
        }

        if err := uc.userRepo.Save(user); err != nil {
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

    func (uc *UserUsecase) GetUser(userIDStr string) (*dto.UserResponse, error) {
        userID, err := model.NewUserID(userIDStr)
        if err != nil {
             return nil, errors.New("invalid user id format")
        }
        user, err := uc.userRepo.FindByID(userID)
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

* **`internal/infrastructure/persistence/user_repository_impl.go`**:

    ```go
    package persistence

    import (
        "errors"
        "gorm.io/gorm"
        "YOUR_PROJECT/internal/domain/model"
        "YOUR_PROJECT/internal/domain/repository"
    )

    // GormUserModel represents the data structure in the database
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

    // toDomain converts Gorm model to domain model using the reconstructor
    func toDomain(gormUser *GormUserModel) (*model.User, error) {
        if gormUser == nil {
            return nil, errors.New("cannot convert nil gorm user model to domain model")
        }
        userID, err := model.NewUserID(gormUser.UserID)
        if err != nil {
            // This indicates data integrity issue in the DB or mapping
            return nil, errors.New("invalid user ID format in database: " + err.Error())
        }

        // Use the Reconstructor function
        user := model.ReconstructUser(userID, gormUser.Name /* ... other fields */)
        return user, nil
    }

    // fromDomain converts domain model to Gorm model for persistence
    func fromDomain(user *model.User) (*GormUserModel, error) {
         if user == nil {
             return nil, errors.New("cannot convert nil domain user to gorm model")
         }
         return &GormUserModel{
             UserID: user.ID().String(), // Use getter and convert UserID
             Name:   user.Name(),        // Use getter
             // ... map other fields using getters ...
         }, nil
    }

    func (r *UserRepositoryImpl) FindByID(id model.UserID) (*model.User, error) {
        var userModel GormUserModel
        // Search logic using GORM, converting UserID to string for query
        err := r.db.First(&userModel, "user_id = ?", id.String()).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil, nil // Return nil, nil for not found (idiomatic in Go)
            }
            return nil, err // Return other DB errors
        }
        return toDomain(&userModel)
    }

    func (r *UserRepositoryImpl) Save(user *model.User) error {
        // Convert domain model to persistence model
        userModel, err := fromDomain(user)
        if err != nil {
            return err // Handle conversion error
        }

        // Save logic using GORM (Create or Update)
        // GORM's Save handles upsert based on primary key presence
        // Ensure the GormUserModel has the correct primary key tag (`gorm:"primaryKey"`)
        return r.db.Save(userModel).Error
    }
    ```

* **`internal/interfaces/http/handler/user_handler.go`**:

    ```go
    package handler

    import (
        "github.com/gin-gonic/gin"
        "YOUR_PROJECT/internal/application/usecase"
        "YOUR_PROJECT/internal/application/dto"
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
        var req dto.CreateUserRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
            return
        }

        res, err := h.userUsecase.CreateUser(&req)
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
        c.JSON(http.StatusCreated, res)
    }

    func (h *UserHandler) GetUser(c *gin.Context) {
        userID := c.Param("id") // Assuming ID is a path parameter like /users/:id

        res, err := h.userUsecase.GetUser(userID)
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
        c.JSON(http.StatusOK, res)
    }
    ```

* **`cmd/server/main.go`**:

    ```go
    package main

    import (
        "github.com/gin-gonic/gin"
        "gorm.io/driver/postgres" // Example
        "gorm.io/gorm"
        "YOUR_PROJECT/internal/domain/repository"
        "YOUR_PROJECT/internal/application/usecase"
        infra "YOUR_PROJECT/internal/infrastructure/persistence"
        interfaces "YOUR_PROJECT/internal/interfaces/http/handler"
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
        userRepo := infra.NewUserRepositoryImpl(db)
        userUsecase := usecase.NewUserUsecase(userRepo)
        userHandler := interfaces.NewUserHandler(userUsecase)

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
