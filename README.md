# Internal Transfers System

HTTP service for managing accounts and processing money transfers.

## Requirements

- Go 1.21+
- PostgreSQL 14+
- Task 3.x ([installation](https://taskfile.dev/installation/))
- tern v2.x (database migration tool)

## Quick Start

1. **Clone and install dependencies**
```bash
git clone <repository-url>
cd internal-transfers
go mod download
```

2. **Install required tools**
```bash
# Install Task (macOS with Homebrew)
brew install go-task/tap/go-task

# Install tern migration tool
go install github.com/jackc/tern/v2@latest

# Ensure Go binaries are in PATH
export PATH=$PATH:~/go/bin
```

3. **Configure database**
```bash
cp env.sample .env
# Edit .env with your database credentials
```

See `env.sample` for all configuration options. At minimum, update the database connection settings.

4. **Run database migrations**
```bash
# Set database connection string
export INTERNAL_TRANSFERS_DB_DSN="postgresql://user:pass@host:port/dbname?sslmode=require"

# Run migrations
task migrations:up
```

5. **Run the application**
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
task migrations:up  # Apply database migrations
```

**Without Task:**
```bash
go run ./cmd/internal-transfers       # Run application
go fmt ./... && go mod tidy          # Format and tidy code
```

API docs: `http://localhost:8080/docs`

## Assumptions

- All accounts operate in a single currency
- Account IDs are unique and provided by the client
- Balance precision is maintained at 5 decimal places
- Negative balances are not allowed
- All transactions are processed synchronously
- No authentication/authorization is implemented (internal service)
- Database migrations must be run manually before starting the application

## Key Features

- Single currency, decimal precision (5 places)
- ACID compliant transactions
- Clean architecture design
- No authentication (internal service)
- Manual migration management for better control

## Testing

```bash
go test ./... -v
```

## Project Structure

```
.
├── cmd/internal-transfers/    # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connection and migrations
│   ├── handler/              # HTTP request handlers
│   ├── middleware/           # HTTP middleware (logging, CORS, etc.)
│   ├── model/                # Domain models
│   ├── repository/           # Data access layer
│   ├── router/               # Route definitions
│   ├── server/               # Server setup
│   └── service/              # Business logic
├── static/                   # OpenAPI documentation
├── env.sample               # Environment configuration template
└── Taskfile.yml             # Task automation
```