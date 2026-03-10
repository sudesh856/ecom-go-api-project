# E-Commerce REST API

A production-ready e-commerce REST API built with **Go**, demonstrating advanced software architecture patterns, database design, and modern API development practices.

## 📋 Overview

This project implements a scalable e-commerce backend service with comprehensive product management and transactional order processing. The API showcases professional software engineering practices including clean architecture, type-safe database interactions, and robust error handling.

## 🏗️ Architecture

The application follows a **hexagonal architecture** (ports and adapters) pattern with clear separation of concerns:

```
cmd/                          # Entry point
├── main.go                    # Application bootstrap & configuration
└── api.go                     # Server setup & route mounting

internal/                      # Business logic (non-exported)
├── adaptors/                  # External service integrations
│   └── postgresql/            # Database layer
│       ├── migrations/        # Schema version control
│       └── sqlc/              # Type-safe generated queries
├── env/                       # Environment configuration
├── json/                      # Serialization utilities
├── products/                  # Product domain
│   ├── handlers.go            # HTTP handlers
│   └── service.go             # Business logic
└── orders/                    # Orders domain
    ├── handlers.go            # HTTP handlers
    ├── service.go             # Business logic with transaction management
    └── types.go               # Domain models
```

**Key Architectural Decisions:**
- **Dependency Injection**: Services receive repositories through constructor injection for testability
- **Interface-Based Design**: Service interfaces define contracts, enabling mock implementations
- **Domain Segregation**: Separate packages for products and orders enforce bounded contexts
- **Database Abstraction**: SQLC-generated repositories provide type-safe queries without ORM overhead

## ✨ Features

### Product Management
- **List all products** with pagination support
- **Fetch individual product** details including pricing and inventory
- **Create new products** with validation (price constraints, default quantities)
- Real-time inventory tracking

### Order Processing
- **Transactional order placement** ensuring data consistency
- **Stock validation** preventing overselling
- **Atomic order item creation** with price snapshots
- **Comprehensive error handling** with custom error types
- **Customer order retrieval** with full order history

### Technical Features
- 🔒 **ACID Compliance**: Database transactions for order atomicity
- 📊 **Type Safety**: SQLC-generated query methods eliminate SQL injection risks
- 🎯 **Middleware Stack**: Request ID tracking, Real IP detection, structured logging, panic recovery
- ⏱️ **Request Timeouts**: 60-second global timeout with per-handler timeouts
- 📝 **Structured Logging**: slog-based logging for production observability
- 🔄 **Connection Pooling**: PostgreSQL connection management through pgx

## 🛠️ Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Runtime** | Go 1.25 | Type-safe, concurrent execution |
| **Web Framework** | Chi mux (v5.2.5) | Fast, lightweight HTTP router |
| **Database** | PostgreSQL 16 | Relational data with ACID guarantees |
| **Database Driver** | pgx/v5 (v5.8.0) | High-performance PostgreSQL driver with transaction support |
| **Query Generation** | SQLC | Type-safe SQL queries with compile-time validation |
| **Migration Tool** | Goose | Database schema versioning |
| **Containerization** | Docker Compose | Development environment consistency |

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sudesh856/ecom-go-api-project.git
   cd ecom-go-api-project
   ```

2. **Configure environment variables:**
   ```bash
   # .env file
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=ecom
   GOOSE_DBSTRING=host=localhost user=postgres password=your_secure_password dbname=ecom sslmode=disable
   ```

3. **Start PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations:**
   ```bash
   goose -dir = internal/adaptors/postgresql/migrations postgres "$GOOSE_DBSTRING" up
   ```

5. **Start the API server:**
   ```bash
   go run cmd/main.go
   ```
   Server runs on `http://localhost:8080`

## 📡 API Endpoints

### Health Check
```http
GET /health
```
**Response:** `200 OK` - "Everything normal."

### Products

#### List All Products
```http
GET /products
```
**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Laptop",
    "price_in_rupees": 50000,
    "quantity": 10,
    "created_at": "2025-03-10T10:30:00Z"
  }
]
```

#### Get Product by ID
```http
GET /products/{id}
```
**Response:** `200 OK` - Single product object  
**Errors:** `400 Bad Request` (invalid ID), `500 Internal Server Error`

#### Create Product
```http
POST /products
Content-Type: application/json

{
  "name": "Wireless Mouse",
  "price_in_rupees": 1500,
  "quantity": 50
}
```
**Response:** `201 Created` - Created product object  
**Validation:** Price ≥ 0

### Orders

#### Place Order
```http
POST /orders
Content-Type: application/json

{
  "customer_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```
**Response:** `201 Created`
```json
{
  "id": 1,
  "customer_id": 1,
  "created_at": "2025-03-10T10:35:00Z"
}
```
**Error Handling:**
- `400 Bad Request` - Missing customer ID or items
- `400 Bad Request` - Product not found
- `400 Bad Request` - Insufficient stock
- `500 Internal Server Error` - Transaction failure

#### Get Order
```http
GET /orders/{id}
```
**Response:** `200 OK` - Order details with items  
**Errors:** `400 Bad Request` (invalid ID), `500 Internal Server Error`

## 🗄️ Database Schema

### Products Table
```sql
CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price_in_rupees INTEGER NOT NULL CHECK (price_in_rupees >= 0),
  quantity INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Orders & Order Items Tables
```sql
CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE order_items (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  quantity INTEGER NOT NULL,
  price_in_rupees INTEGER NOT NULL,
  CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id)
);
```

## 💡 Key Implementation Details

### Transaction Management
The order service implements database transactions to ensure order atomicity:
- Creates order record
- Validates product availability
- Creates order items with price snapshots
- Automatically rolls back on any failure

```go
tx, err := s.db.Begin(ctx)
defer tx.Rollback(ctx) // Auto-rollback on function exit
// ... operations ...
tx.Commit(ctx)
```

### Error Handling
Custom error types for domain-specific error conditions:
```go
ErrProductNotFound = errors.New("Product not found.")
ErrProductNoStock = errors.New("No stock.")
```

### Type Safety
SQLC generates strongly-typed query methods:
- `ListProducts(ctx context.Context) ([]Product, error)`
- `CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)`
- Compile-time SQL validation

### Request Processing Pipeline
1. **Middleware Stack**: Request ID → Real IP → Logging → Panic Recovery → Timeout
2. **Handler Layer**: HTTP parsing and validation
3. **Service Layer**: Business logic and data operations
4. **Repository Layer**: Type-safe database queries
5. **Data Layer**: PostgreSQL transaction handling

## 🧪 Error Handling & Validation

- **Input Validation**: Customer ID and order items validation
- **Business Logic Validation**: Stock availability checks
- **Database Constraints**: Price non-negativity, referential integrity
- **Structured Error Responses**: Consistent error messages and status codes
- **Panic Recovery**: Middleware-level recovery prevents server crashes

## 📦 Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| chi/v5 | v5.2.5 | HTTP routing |
| pgx/v5 | v5.8.0 | PostgreSQL driver |
| pgpassfile | v1.0.0 | PostgreSQL auth |
| golang.org/x/text | v0.29.0 | Unicode handling |

## 🔒 Security Considerations

- **SQL Injection Prevention**: SQLC type-safe queries eliminate manual SQL
- **Connection String Management**: Environment variable configuration
- **No Hard-coded Credentials**: Configuration via environment variables
- **Request Timeout Protection**: 60-second global timeout prevents resource exhaustion
- **Panic Recovery**: Prevents information leakage through stack traces

## 🎯 Learning Outcomes

This project demonstrates proficiency in:

✅ **Go Web Development**: Chi router, HTTP handlers, middleware patterns  
✅ **Relational Databases**: Schema design, migrations, transactions, constraints  
✅ **Clean Architecture**: Separation of concerns, dependency injection, interfaces  
✅ **Type Safety**: SQLC generation, compile-time query validation  
✅ **Production Practices**: Error handling, logging, transaction management  
✅ **Container Technology**: Docker Compose, development environment setup  
✅ **API Design**: RESTful principles, HTTP status codes, error handling  
✅ **Concurrent Safe Operations**: Context propagation, transaction handling  

## 📈 Scalability & Future Enhancements

- **Caching Layer**: Redis integration for product catalog caching
- **Connection Pooling**: pgxpool for high-concurrency scenarios
- **Metrics Collection**: Prometheus integration for monitoring
- **Authentication**: JWT-based customer authentication
- **Authorization**: Role-based access control (RBAC)
- **API Documentation**: OpenAPI/Swagger specification
- **Pagination**: Limit/offset parameters for list endpoints
- **Soft Deletes**: Archive orders and products instead of deletion

## 📄 License

This project is open source and available under the MIT License.

## 👤 Author

**Sudesh** - [GitHub](https://github.com/sudesh856)

---

**Last Updated:** March 10, 2026  
**Go Version:** 1.25  
**Status:** ✅ Production Ready