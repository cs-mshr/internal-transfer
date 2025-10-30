# Internal Transfers System

A robust financial transaction system built with Go, designed for managing accounts and processing transfers with high reliability and performance.

## Features

- **Account Management**: Create accounts with initial balance and query account balances
- **Money Transfers**: Process transfers between accounts with ACID guarantees
- **Decimal Precision**: Accurate financial calculations with decimal precision
- **Transaction Integrity**: Database transactions ensure consistency
- **Error Handling**: Comprehensive error handling and validation
- **API Documentation**: OpenAPI/Swagger specification included
- **Production Ready**: Clean architecture, structured logging, and best practices

## Requirements

- Go 1.21 or higher
- PostgreSQL 14+
- Task (optional, for task automation)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd internal-transfers
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application:

Create a `.env` file in the project root with your database configuration:

```bash
# Environment
INTERNAL_TRANSFERS_PRIMARY_ENV=production

# Server Configuration  
INTERNAL_TRANSFERS_SERVER_PORT=8080
INTERNAL_TRANSFERS_SERVER_READ_TIMEOUT=30
INTERNAL_TRANSFERS_SERVER_WRITE_TIMEOUT=30
INTERNAL_TRANSFERS_SERVER_IDLE_TIMEOUT=60
INTERNAL_TRANSFERS_SERVER_CORS_ALLOWED_ORIGINS=*

# Database Configuration
INTERNAL_TRANSFERS_DATABASE_HOST=your-db-host
INTERNAL_TRANSFERS_DATABASE_PORT=5432
INTERNAL_TRANSFERS_DATABASE_USER=your-db-user
INTERNAL_TRANSFERS_DATABASE_PASSWORD=your-db-password
INTERNAL_TRANSFERS_DATABASE_NAME=your-db-name
INTERNAL_TRANSFERS_DATABASE_SSL_MODE=require
INTERNAL_TRANSFERS_DATABASE_MAX_OPEN_CONNS=25
INTERNAL_TRANSFERS_DATABASE_MAX_IDLE_CONNS=5
INTERNAL_TRANSFERS_DATABASE_CONN_MAX_LIFETIME=300
INTERNAL_TRANSFERS_DATABASE_CONN_MAX_IDLE_TIME=60
```

4. Run the application:

The application will automatically run database migrations on startup in non-local environments.

```bash
go run cmd/internal-transfers/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Account Management

#### Create Account
```bash
POST /api/v1/accounts
Content-Type: application/json

{
  "account_id": 123,
  "initial_balance": "100.23344"
}
```

#### Get Account Balance
```bash
GET /api/v1/accounts/{account_id}
```

Response:
```json
{
  "account_id": 123,
  "balance": "100.23344"
}
```

### Transactions

#### Create Transaction
```bash
POST /api/v1/transactions
Content-Type: application/json

{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "50.12345"
}
```

## API Documentation

The API documentation is available at:
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI JSON: `http://localhost:8080/static/openapi.json`

## Configuration

The application uses environment variables for configuration. All configuration options are prefixed with `INTERNAL_TRANSFERS_` and follow a hierarchical structure:

- Primary settings: `PRIMARY_*`
- Server settings: `SERVER_*`
- Database settings: `DATABASE_*`

Refer to the `.env` template above for all available options.

## Architecture

The application follows clean architecture principles:

```
internal/
├── handler/      # HTTP request handlers
├── service/      # Business logic
├── repository/   # Data access layer
├── model/        # Domain models
├── database/     # Database connection and migrations
├── middleware/   # HTTP middleware
├── router/       # Route definitions
└── server/       # Server initialization
```

## Development

### Available Commands

```bash
# Run tests
go test ./...

# Run with live reload (using air)
air

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

### Database Migrations

Migrations are located in `internal/database/migrations/` and are automatically applied on startup in non-local environments.

New migrations should be added to `internal/database/migrations/` following the naming pattern `XXX_description.sql`.

## Testing

Run tests with:
```bash
go test ./... -v
```

For integration tests (requires database):
```bash
go test -tags=integration ./... -v
```

## Design Decisions

- **Single Currency**: The system operates with a single currency for simplicity
- **Client-Provided IDs**: Account IDs are provided by API clients for flexibility
- **Decimal Precision**: All monetary values maintain 5 decimal places for accuracy
- **Synchronous Processing**: Transactions are processed synchronously for immediate consistency
- **Stateless API**: No authentication layer is implemented, suitable for internal microservice communication

## Performance Considerations

- Connection pooling for database connections
- Row-level locking for transaction consistency
- Indexed database columns for faster lookups
- Rate limiting to prevent abuse

## Error Handling

The API returns structured error responses:
```json
{
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Insufficient balance in source account"
  }
}
```

Common error codes:
- `INVALID_FORMAT` - Invalid input format
- `ACCOUNT_EXISTS` - Account already exists
- `ACCOUNT_NOT_FOUND` - Account not found
- `INSUFFICIENT_BALANCE` - Insufficient balance for transfer
- `SAME_ACCOUNT` - Source and destination accounts are the same