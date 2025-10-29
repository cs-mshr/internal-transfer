# Internal Transfers System

HTTP service for managing accounts and processing money transfers.

## Requirements

- Go 1.21+
- PostgreSQL 14+
- Task (optional, [installation](https://taskfile.dev/installation/))

## Quick Start

1. **Clone and install**
```bash
git clone <repository-url>
cd internal-transfers
go mod download
```

2. **Configure**
Create `.env` file:
```bash
INTERNAL_TRANSFERS_PRIMARY_ENV=production
INTERNAL_TRANSFERS_DATABASE_HOST=your-db-host
INTERNAL_TRANSFERS_DATABASE_PORT=5432
INTERNAL_TRANSFERS_DATABASE_USER=your-db-user
INTERNAL_TRANSFERS_DATABASE_PASSWORD=your-db-password
INTERNAL_TRANSFERS_DATABASE_NAME=your-db-name
INTERNAL_TRANSFERS_DATABASE_SSL_MODE=require
```

Additional settings (optional):
```bash
INTERNAL_TRANSFERS_SERVER_PORT=8080
INTERNAL_TRANSFERS_SERVER_READ_TIMEOUT=30
INTERNAL_TRANSFERS_SERVER_WRITE_TIMEOUT=30
INTERNAL_TRANSFERS_SERVER_IDLE_TIMEOUT=60
INTERNAL_TRANSFERS_DATABASE_MAX_OPEN_CONNS=25
INTERNAL_TRANSFERS_DATABASE_MAX_IDLE_CONNS=5
```

3. **Run**
```bash
task run
# or without task:
go run ./cmd/internal-transfers
```

API runs at `http://localhost:8080`

## API Endpoints

### Create Account
```
POST /api/v1/accounts
{
  "account_id": 123,
  "initial_balance": "100.23344"
}
```

### Get Account
```
GET /api/v1/accounts/{account_id}
```

### Create Transaction
```
POST /api/v1/transactions
{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "50.12345"
}
```

## Development

**With Task:**
```bash
task help           # Show available tasks
task run            # Run application
task tidy           # Format and tidy code
task migrations:new name=<name>  # Create migration
```

**Without Task:**
```bash
go run ./cmd/internal-transfers       # Run application
go fmt ./... && go mod tidy          # Format and tidy code
```

API docs: `http://localhost:8080/docs`

## Key Features

- Single currency, decimal precision (5 places)
- ACID compliant transactions
- Clean architecture design
- No authentication (internal service)
- Automatic migrations on startup

## Testing

```bash
go test ./... -v
```