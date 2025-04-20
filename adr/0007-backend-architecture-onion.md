# ADR 0007: Backend Architecture - Onion Architecture

## Status

Proposed

## Context

In OpsCore backend development, it is necessary to introduce a consistent architectural pattern to enhance maintainability, testability, and scalability. Currently, there are no clear architectural guidelines, which could make future feature additions and changes difficult. By thoroughly separating concerns based on the principles of Clean Architecture, we aim to address these challenges.

## Decision

We will adopt the **Onion Architecture** for the backend architecture. This is based on the principle of directing dependencies unidirectionally from the inside (Domain) to the outside (Infrastructure).

**Layer Structure and Responsibilities:**

1. **Domain Layer:**
    * **Responsibility:** Represents the core business logic and rules of the application. Includes entities, value objects, domain events, and repository interfaces.
    * **Dependencies:** Does not depend on external layers.
    * **Folder:** `internal/domain/`
        * `model/`: Entities, Value Objects
        * `repository/`: Repository interface definitions
        * `service/`: Domain services (logic spanning multiple entities)

2. **Application Layer:**
    * **Responsibility:** Implements use cases. Manipulates domain objects and controls application-specific flows. Depends on infrastructure layer interfaces (e.g., repository interfaces).
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

    type User struct {
        ID   string
        Name string
        // ... other fields
    }

    // Factory functions like NewUser
    ```

* **`internal/domain/repository/user_repository.go`**:

    ```go
    package repository

    import "YOUR_PROJECT/internal/domain/model"

    type UserRepository interface {
        FindByID(id string) (*model.User, error)
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
        "net/http"
    )

    type UserUsecase struct {
        userRepo repository.UserRepository
    }

    func NewUserUsecase(ur repository.UserRepository) *UserUsecase {
        return &UserUsecase{userRepo: ur}
    }

    func (uc *UserUsecase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
        user := model.User{ /* ... generate User from req ... */ }
        if err := uc.userRepo.Save(&user); err != nil {
            return nil, err // Error handling
        }
        // ... generate UserResponse DTO from User ...
        return &dto.UserResponse{ /* ... */ }, nil
    }
    ```

* **`internal/infrastructure/persistence/user_repository_impl.go`**:

    ```go
    package persistence

    import (
        "gorm.io/gorm"
        "YOUR_PROJECT/internal/domain/model"
        "YOUR_PROJECT/internal/domain/repository"
    )

    type UserRepositoryImpl struct {
        db *gorm.DB
    }

    func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
        return &UserRepositoryImpl{db: db}
    }

    func (r *UserRepositoryImpl) FindByID(id string) (*model.User, error) {
        // Search logic using GORM
    }

    func (r *UserRepositoryImpl) Save(user *model.User) error {
        // Save logic using GORM
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
    )

    type UserHandler struct {
        userUsecase *usecase.UserUsecase
    }

    func NewUserHandler(uu *usecase.UserUsecase) *UserHandler {
        return &UserHandler{userUsecase: uu}
    }

    func (h *UserHandler) CreateUser(c *gin.Context) {
        var req dto.CreateUserRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        res, err := h.userUsecase.CreateUser(&req)
        if err != nil {
            // Error handling (e.g., domain error vs system error)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }
        c.JSON(http.StatusCreated, res)
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
