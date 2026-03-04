# Go Boilerplate Project

A production-ready Go web API boilerplate using the **Factory Pattern** with clean layered architecture. It includes MongoDB, MySQL, Redis support, HTTP API helpers, transaction management, network service for external API calls, multilang (i18n) service for localized API responses, structured logging, and **global + application middlewares**.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Understanding the Architecture](#understanding-the-architecture)
- [Using Existing Features](#using-existing-features)
- [Middlewares](#middlewares)
- [Customization Guide](#customization-guide)
- [API Reference](#api-reference)

---

## Prerequisites

Before installing, ensure you have:

| Requirement | Version | Purpose |
|-------------|---------|---------|
| **Go** | 1.21+ | Runtime |
| **MongoDB** | 4.0+ | Document database |
| **MySQL** | 5.7+ or 8.0+ | Relational database |
| **Redis** | 6.0+ | Caching / session store |

---

## Installation

### Step 1: Clone or Download the Repository

```bash
# Using Git
git clone <repository-url> GoBoilerPlateFactoryPattern
cd GoBoilerPlateFactoryPattern

# Or download and extract the ZIP, then:
cd GoBoilerPlateFactoryPattern
```

### Step 2: Install Dependencies

```bash
go mod download
```

This fetches all Go modules listed in `go.mod` (Gin, GORM, MongoDB driver, Redis client, Zap logger, etc.).

### Step 3: Environment Setup

1. Copy the sample environment file:

   ```bash
   copy env_sample.text .env
   ```

   On Unix/macOS:

   ```bash
   cp env_sample.text .env
   ```

2. Edit `.env` and set your database credentials, ports, and other config (see [Configuration](#configuration)).

### Step 4: Ensure Databases Are Running

- **MongoDB**: Running on port `27017` (default)
- **MySQL**: Running on port `3306` with a database created
- **Redis**: Running on port `6379`

---

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ENV_MODE` | `development`, `staging`, or `production` | `development` |
| `API_BASE_URL` | Base URL of the API | `http://localhost:8080` |
| `APP_NAME` | Application name | `Go Boilerplate Project` |
| `APP_VERSION` | Version string | `1.0.0` |
| `APP_PORT` | Port the server listens on | `8080` |
| `APP_HOST` | Server host | `localhost` |
| `MONGO_DB_HOST` | MongoDB host | `localhost` |
| `MONGO_DB_PORT` | MongoDB port | `27017` |
| `MONGO_DB_DATABASE` | MongoDB database name | `go_boilerplate_project` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_DATABASE` | Redis DB index | `0` |
| `REDIS_PASSWORD` | Redis password (if any) | (leave empty if none) |
| `MYSQL_HOST` | MySQL host | `localhost` |
| `MYSQL_PORT` | MySQL port | `3306` |
| `MYSQL_DATABASE` | MySQL database name | `go_boilerplate_project` |
| `MYSQL_USERNAME` | MySQL username | `root` |
| `MYSQL_PASSWORD` | MySQL password | `your_password` |
| `LOGGER_LEVEL` | `debug`, `info`, `warn`, `error` | `info` |
| `LOGGER_FILE_PATH` | Log file path | `logs/app.log` |
| `LOGGER_MAX_SIZE` | Max log file size (MB) | `100` |
| `LOGGER_MAX_BACKUPS` | Number of rotated backups | `10` |
| `LOGGER_MAX_AGE` | Max age of backups (days) | `30` |
| `LOGGER_COMPRESS` | Compress rotated logs | `false` |

---

## Running the Application

```bash
go run apps/service_one/cmd/main.go
```

Or build and run:

```bash
go build -o bin/app apps/service_one/cmd/main.go
./bin/app
```

The API will start at `http://localhost:8080` (or your configured `APP_PORT`).

---

## Project Structure

```
GoBoilerPlateFactoryPattern/
├── apps/
│   └── service_one/                    # First microservice / app
│       ├── cmd/
│       │   ├── main.go                # Entry point
│       │   └── dependencies.go        # DI wiring (databases, helpers, services, layers)
│       ├── env/
│       │   └── env_service_one.go     # Loads .env for this service
│       ├── layers/
│       │   ├── data/                  # Data access (MySQL, MongoDB, Redis)
│       │   │   ├── brands/            # Brand CRUD operations
│       │   │   └── products/          # Product CRUD operations
│       │   ├── domain/                # Business logic
│       │   │   ├── brands/            # Brand domain rules
│       │   │   └── products/         # Product domain rules
│       │   └── http/                 # HTTP handlers (controllers)
│       │       ├── brands/           # Brand API handlers
│       │       └── products/         # Product API handlers
│       ├── middlewares/               # Application middlewares (service-specific)
│       │   ├── serviceone_middleware_repository.go
│       │   └── serviceone_middleware_service.go
│       └── router/                    # Route definitions (Gin routes, middleware)
│           ├── serviceone_router_repository.go
│           └── serviceone_router_service.go
├── middlewares/
│   └── global/                        # Applied to all service routes
│       ├── global_middleware_repository.go
│       └── global_middleware_service.go
├── constants/
│   ├── api/                           # API constants (pagination, base URLs)
│   ├── env/                           # Env mode constants
│   ├── lang/                          # Language header, default lang
│   ├── enums/                         # Enums (e.g. HTTP methods)
│   ├── errors/                        # Error definitions
│   ├── numeric/                       # Numeric constants
│   └── strings/                      # String constants
├── models/
│   ├── api/                           # API response, pagination models
│   ├── brands/                        # Brand entities & DTOs
│   ├── commons/                       # Shared models (context, etc.)
│   ├── dependencies/                  # DI model types
│   ├── env/                           # Environment config models
│   ├── network/                      # Network service input/output
│   ├── products/                      # Product entities & DTOs
│   └── databases/
│       ├── mysql/                     # MySQL helper input/output
│       ├── mongodb/                   # MongoDB helper input/output
│       └── redis/                     # Redis helper input/output
├── services/
│   ├── context/                      # Context factory (timeout, cancel)
│   ├── multilang/                    # Multi-language (go-i18n) for API errors
│   ├── logger/                       # Zap logger setup
│   ├── network/                      # HTTP fetch service (call external APIs)
│   ├── transactions/                 # MongoDB & MySQL transaction runner
│   └── helpers/
│       ├── api/                      # Parse body, send responses, validation
│       ├── custom/                    # Custom utilities
│       └── db/
│           ├── mysql/                # GORM CRUD helpers
│           ├── mongodb/               # MongoDB CRUD helpers
│           └── redis/                 # Redis helpers
├── storage/
│   ├── mongodb/                      # MongoDB connection & reconnect logic
│   ├── mysql/                        # MySQL connection & reconnect logic
│   └── redis/                        # Redis connection & reconnect logic
├── languages/                        # i18n JSON files (en.json, es.json, fr.json)
├── logs/                             # Log output directory
├── .env                              # Environment variables (create from env_sample.text)
├── env_sample.text                   # Sample .env template
├── go.mod
└── go.sum
```

---

## Understanding the Architecture

### Layer Flow

```
HTTP Request → Router → HTTP Layer → Domain Layer → Data Layer → DB/Storage
              ↓           ↓              ↓             ↓
           Gin       Validation      Business       MySQL/MongoDB
                     & Parsing       Logic         Redis helpers
```

| Layer | Purpose | Files |
|-------|---------|-------|
| **Router** | Defines routes, applies middleware | `router/serviceone_router_*` |
| **HTTP** | Parses requests, calls domain, sends responses | `http/brands`, `http/products` |
| **Domain** | Business rules, orchestrates data layer, transactions | `domain/brands`, `domain/products` |
| **Data** | Database operations via helpers | `data/brands`, `data/products` |
| **Storage** | DB connection management (connect, reconnect) | `storage/mysql`, `storage/mongodb`, `storage/redis` |
| **Services** | Cross-cutting: context, transactions, network, multilang, logger | `services/*` |
| **Helpers** | Reusable CRUD and API utilities | `services/helpers/*` |

### Factory Pattern

Components are initialized via `InitService` with an `Input` struct. Dependencies are injected and validated. Example:

```go
// dependencies.go - Services (e.g. transactions, multilang)
transactions := transactions_service.InitService(transactions_service.Input{
    Helpers:  &transactions_service.Helpers{ MongoDB: ..., MySQL: ... },
    Services: &transactions_service.Services{ Context: ... },
    Logger:   d.Logger,
})
d.Services.Transactions = transactions

multilang := multilang_service.InitService(multilang_service.Input{
    Logger:      d.Logger,
    DefaultLang: "en",
    Bundle:      bundle,
})
d.Services.Multilang = multilang
```

### Naming Conventions

| Pattern | Purpose |
|---------|---------|
| `*_repository.go` | Interface definitions, `InitService`, `Input` struct, validation |
| `*_service.go` | Implementation of repository methods |
| `*_models.go` | Data structures, request/response DTOs |
| `depedencies_models.go` | Dependency injection type definitions |

---

## Using Existing Features

### 1. Brands API

**Edit Brand** – Update a brand by ID.

```
POST /api/v1/brands
Content-Type: application/json

{
  "name": "New Brand Name",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

Validation: `name` (3–100 chars), `id` (UUID).  
Logic: Updates brand in MySQL and related product brand names in a transaction.

---

### 2. Products API

**Create Product** – Create a product linked to a brand.

```
POST /api/v1/products
Content-Type: application/json

{
  "name": "Product Name",
  "brandId": "550e8400-e29b-41d4-a716-446655440000"
}
```

Validation: `name` (3–100 chars), `brandId` required.  
Logic: Resolves brand name from MySQL, creates product in MySQL and MongoDB.

---

### 3. Network Service (External API Calls)

Use the network service to call external APIs from your code.

```go
var resp MyResponseStruct
output, err := d.Services.Network.Fetch(ctx, &network_models.FetchInput{
    Route:        "https://api.example.com/users",
    Method:       network_models.HTTPMethodGet,
    Headers:      map[string]string{"Authorization": "Bearer token"},
    QueryParams:  map[string]string{"page": "1"},
    Timeout:      10 * time.Second,
    ResponseModel: &resp,  // optional: unmarshal JSON into this struct
})
if err != nil {
    // handle error
}
// output.BodyBytes = raw response
// output.ParsedModel / resp = parsed struct
// output.StatusCode, output.Headers, output.Duration, etc.
```

---

### 4. Transactions

Run MongoDB or MySQL operations inside a transaction:

```go
// MySQL
err := d.Services.Transactions.RunMySQLTransaction(&mysql_models.TransactionInput{
    Callback: func(tx *gorm.DB) error {
        // use tx for all DB operations in this block
        return nil
    },
})

// MongoDB
err := d.Services.Transactions.RunMongoDBTransaction(&mongodb_models.TransactionInput{
    Callback: func(ctx context.Context) error {
        // use ctx in MongoDB helper calls
        return nil
    },
})
```

---

### 5. API Helpers

Use in HTTP handlers for parsing and responses. API helpers receive the multilang service as a dependency for translating validation errors and messages.

- `ParseJSONBody(c, dest)` – Bind JSON to struct, returns validation errors
- `SendSuccess(c, statusCode, message, data)` – Send success response
- `SendError(c, statusCode, message, errDetail)` – Send error response
- `SendApiResponse(c, response)` – Unified response handler
- `GetPaginationFromQuery(c)` – Parse `page_number` and `page_size`

---

### 6. Multi-Language Support (i18n)

API validation errors are translated using the **multilang service** ([go-i18n](https://github.com/nicksnyder/go-i18n)). The service is initialized in `dependencies.go` and injected into API helpers for localized responses. Language files live in `languages/` at the project root.

**Language header:** `X-Language` or `Accept-Language` (e.g. `en`, `es`, `fr`)

```bash
# English (default)
curl -X POST http://localhost:8080/api/v1/brands \
  -H "Content-Type: application/json" \
  -d '{"name":"x"}'   # name too short

# Spanish
curl -X POST http://localhost:8080/api/v1/brands \
  -H "Content-Type: application/json" \
  -H "X-Language: es" \
  -d '{"name":"x"}'
```

Response messages vary by language. Supported: **en** (default), **es**, **fr**. Add more by creating `languages/<code>.json` (e.g. `languages/de.json`) and registering it in `multilang_service.InitBundle` (called from `dependencies.go` before initializing the multilang service).

---

## Middlewares

The project uses two middleware layers, wired via dependencies and applied in the router.

### Global Middleware

Applied to **all routes** in `ConfigureRouter()`, before any handler runs.

| Middleware | Purpose |
|------------|---------|
| **RequestID** | Adds `X-Request-ID` header to each request (generates if not provided) for tracing and debugging. Stored in Gin context and echoed in response. |

Location: `middlewares/global/`

### Application Middleware

Applied only to **service-specific route groups** (e.g. `/api/v1/*`) in `SetupRoutes()`. Lives under each app so each service can have its own middlewares.

| Middleware | Purpose |
|------------|---------|
| **AppVersion** | Adds `X-App-Name` and `X-App-Version` response headers from environment config. |

Location: `apps/service_one/middlewares/`

### Execution Order

For a request to `/api/v1/brands`:

1. `gin.Recovery()` – panic recovery  
2. `gin.Logger()` – request logging  
3. **Global:** RequestID  
4. **Application:** AppVersion  
5. Handler (e.g. EditBrand)

### Response Headers

API responses include:

- `X-Request-ID` – unique request identifier  
- `X-App-Name` – application name (from `APP_NAME`)  
- `X-App-Version` – application version (from `APP_VERSION`)

### Adding More Middlewares

**Global** – In `middlewares/global/global_middleware_service.go`:

1. Implement a new method that returns `gin.HandlerFunc`.
2. Append it in `GetMiddlewares()`:
   ```go
   return []gin.HandlerFunc{
       s.RequestID(),
       s.YourNewMiddleware(),
   }
   ```

**Application** – In `apps/service_one/middlewares/serviceone_middleware_service.go`:

1. Implement a new method that returns `gin.HandlerFunc`.
2. Append it in `GetMiddlewares()`:
   ```go
   return []gin.HandlerFunc{
       s.AppVersion(),
       s.YourNewMiddleware(),
   }
   ```

The router applies these automatically via `GetMiddlewares()` in `ConfigureRouter` (global) and `SetupRoutes` (application).

---

## Customization Guide

### Add a New Entity (e.g., "Categories")

1. **Models** – `models/categories/category_models.go`
   - Define struct and request DTOs with validation tags.

2. **Data Layer**
   - `apps/service_one/layers/data/categories/categories_data_repository.go` – Input, validation, InitService.
   - `apps/service_one/layers/data/categories/categories_data_service.go` – CRUD via helpers.

3. **Domain Layer**
   - `apps/service_one/layers/domain/categories/categories_domain_repository.go`
   - `apps/service_one/layers/domain/categories/categories_domain_service.go` – Business rules.

4. **HTTP Layer**
   - `apps/service_one/layers/http/categories/categories_http_repository.go`
   - `apps/service_one/layers/http/categories/categories_http_service.go` – HTTP handlers.

5. **Dependencies**
   - Add `Categories` to `Data`, `Domain`, `Http` in `models/dependencies/depedencies_models.go`.
   - Wire in `dependencies.go`: `initDataLayers`, `initDomainLayers`, `initHttpLayers`.

6. **Router**
   - In `serviceone_router_service.go`, add routes:
     ```go
     v1.GET("/categories", s.Http.Categories.List)
     v1.POST("/categories", s.Http.Categories.Create)
     ```

7. **Auto-migrate**
   - Call `Helpers.MySQL.AutoMigrate(&categories_models.Category{})` in `categories_data_repository.InitService`.

---

### Add a New Microservice (e.g., `service_two`)

1. Copy `apps/service_one` → `apps/service_two`.
2. Add `env_service_two.go` if env differs.
3. Update `cmd/main.go` to start `service_two` instead of or alongside `service_one`.
4. Adjust `dependencies` for the new service’s needs.

---

### Add Environment Variables

1. In `models/env/env_models.go`, add fields to the appropriate struct.
2. In `apps/service_one/env/env_service_one.go`, read them with `getEnvAsString`, `getEnvAsInt`, etc.
3. Add new variables to `env_sample.text` and your `.env`.

---

### Add a New Database

1. Add storage package: `storage/postgres/postgres_service.go` (pattern similar to MySQL/MongoDB).
2. Add env config: `models/env/env_models.go`, `env_service_one.go`.
3. Add helper: `services/helpers/db/postgres/`.
4. Add models: `models/databases/postgres/`.
5. Wire in `dependencies.go` and pass to helpers/services as needed.

---

### Add a New Middleware

See [Middlewares – Adding More Middlewares](#adding-more-middlewares) above. Use `middlewares/global/` for routes-wide behavior, or `apps/service_one/middlewares/` for application-specific behavior.

---

### Add a New Route

In `serviceone_router_service.go` → `SetupRoutes()`:

```go
v1.GET("/my-endpoint", s.Http.MyHandler.MyMethod)
```

Ensure `MyHandler` is wired in `initHttpLayers` and added to `Http` in dependencies.

---

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/brands` | Edit brand by ID |
| `POST` | `/api/v1/products` | Create product |

Base URL: `http://localhost:8080` (configurable via `APP_PORT`).

Response headers for all `/api/v1` endpoints: `X-Request-ID`, `X-App-Name`, `X-App-Version`.

---

## Logging

Logs go to:

- **Console** – Colored output
- **File** – `logs/app.log` (rotation via Lumberjack)

Set `LOGGER_LEVEL` to `debug`, `info`, `warn`, or `error`.

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `Failed to connect to MongoDB/MySQL/Redis` | Ensure the service is running, check host/port in `.env`, verify firewall |
| `Failed to load environment` | Create `.env` from `env_sample.text`, ensure file is in project root |
| `validation failed` | Check request body matches DTO (e.g. `name` 3–100 chars, `id` as valid UUID for brands) |
| `brandId` not found on product create | Create the brand first or use an existing brand UUID |
| Logs not appearing | Check `LOGGER_LEVEL` (use `debug` for more output), ensure `logs/` directory exists |

---

## License

[Add your license here]
