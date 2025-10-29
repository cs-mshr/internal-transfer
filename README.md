# Internal Transfers System

A high-performance internal transfers system built with Go, providing HTTP endpoints for account management and money transfers between accounts.

## Features

- **Account Management**: Create accounts with initial balance and query account balances
- **Money Transfers**: Process transfers between accounts with ACID guarantees
- **Decimal Precision**: Accurate financial calculations with decimal precision
- **Transaction Integrity**: Database transactions ensure consistency
- **Error Handling**: Comprehensive error handling and validation
- **API Documentation**: OpenAPI/Swagger specification included
- **Production Ready**: Clean architecture, structured logging, and best practices

## Requirements

- Go 1.24 or higher
- PostgreSQL 16+
- Make (optional, for using Taskfile)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/cs-mshr/internal-transfer.git
cd internal-transfers
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.sample .env
# Edit .env with your database configuration
```

4. Start PostgreSQL database (if not already running):
```bash
# Using Docker
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16

# Create database
docker exec -it postgres psql -U postgres -c "CREATE DATABASE internal_transfers;"
```

5. Run database migrations:
```bash
# Using task (if installed)
task migrations:up

# Or manually
go run cmd/internal-transfers/main.go
```

6. Start the server:
```bash
# Using task
task run

# Or directly
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

The application uses environment variables for configuration. Key variables:

- `INTERNAL_TRANSFERS_DATABASE_HOST` - PostgreSQL host (default: localhost)
- `INTERNAL_TRANSFERS_DATABASE_PORT` - PostgreSQL port (default: 5432)
- `INTERNAL_TRANSFERS_DATABASE_USER` - Database user (default: postgres)
- `INTERNAL_TRANSFERS_DATABASE_PASSWORD` - Database password
- `INTERNAL_TRANSFERS_DATABASE_NAME` - Database name (default: internal_transfers)
- `INTERNAL_TRANSFERS_SERVER_PORT` - Server port (default: 8080)

See `.env.sample` for all available configuration options.

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

To create a new migration:
```bash
task migrations:new name=your_migration_name
```

## Testing

Run tests with:
```bash
go test ./... -v
```

For integration tests (requires database):
```bash
go test -tags=integration ./... -v
```

## Assumptions

- All accounts use the same currency
- Account IDs are provided by the client
- Decimal precision is maintained up to 5 decimal places
- Transactions are processed synchronously
- No authentication/authorization is implemented (as per requirements)

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