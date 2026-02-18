# Payment Processing System

A multi-currency payment processing system built with Go and PostgreSQL.

## Features

- Internal payments between users
- External payments to external accounts
- Multi-currency support (USD, EUR, GBP)
- Double-entry bookkeeping
- Payment request staging and validation
- External account caching

## Project Structure

```
payment-system/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   ├── valueobjects/
│   │   └── repositories/
│   ├── application/
│   │   ├── services/
│   │   └── usecases/
│   ├── infrastructure/
│   │   ├── postgres/
│   │   ├── exchange/
│   │   └── external/
│   └── presentation/
│       ├── handlers/
│       └── dto/
├── migrations/
├── pkg/
│   └── errors/
└── go.mod
```

## Database

PostgreSQL.

## Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher

### Installation

1. Clone the repository and navigate to the project directory:
```bash
cd payment-system
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
createdb payment_system
```

4. Run migrations:
```bash
psql -d payment_system -f migrations/001_initial_schema.up.sql
psql -d payment_system -f migrations/002_test_data.up.sql
```

   The second migration (`002_test_data.up.sql`) creates test users and accounts needed for testing the API endpoints.

5. Set environment variables (optional):
```bash
export DATABASE_URL="postgres://postgres:yourpassword@localhost:5432/payment_system?sslmode=disable"
export PORT=4078
```

6. Run the server:
```bash
go run cmd/api/main.go
```

The server will start on port 4078 by default.

## Testing

### Test Data

The `002_test_data.up.sql` migration creates test users and accounts for testing:

**Test Users:**
- **Alice** (`11111111-1111-1111-1111-111111111111`)
  - USD: $10,000.00
  - EUR: €5,000.00
  - GBP: £3,000.00
- **Bob** (`22222222-2222-2222-2222-222222222222`)
  - USD: $5,000.00
  - EUR: €2,500.00
  - GBP: £1,500.00
- **Charlie** (`33333333-3333-3333-3333-333333333333`)
  - USD: $2,000.00
  - EUR: €1,000.00
  - GBP: £500.00

### Postman Collection

A Postman collection (`postman_collection.json`) is included with example requests. The collection is pre-configured with test user IDs from the test data migration. Make sure to run the `002_test_data.up.sql` migration before using the Postman collection.

## API Endpoints

### Internal Payment
- **POST** `/api/v1/payments/internal`
  - Process payment between internal users
  - Request body:
    ```json
    {
      "fromUserId": "uuid",
      "toUserId": "uuid",
      "amount": "100.00",
      "currency": "USD",
      "reference": "unique-reference"
    }
    ```

### External Payment Request
- **POST** `/api/v1/payments/external`
  - Create external payment request (requires validation)
  - Request body:
    ```json
    {
      "fromUserId": "uuid",
      "amount": "100.00",
      "currency": "USD",
      "targetCurrency": "EUR",
      "externalAccountNumber": "123456789",
      "externalIban": "GB82WEST12345698765432",
      "externalBankName": "Bank Name",
      "reference": "unique-reference",
      "idempotencyKey": "optional-key",
      "saveAccount": false
    }
    ```

### Process Payment Request
- **POST** `/api/v1/payments/requests/:id/process`
  - Process a validated external payment request

### Health Check
- **GET** `/health`
  - Check server health status

## Architecture

The system follows Domain-Driven Design (DDD) principles with:
- **Domain Layer**: Core business entities and value objects
- **Application Layer**: Use cases and services
- **Infrastructure Layer**: Database implementations and external integrations
- **Presentation Layer**: HTTP handlers and DTOs

## Features

- ✅ Multi-currency support (USD, EUR, GBP)
- ✅ Internal payments (instant processing)
- ✅ External payments (staged validation)
- ✅ External account caching
- ✅ Double-entry bookkeeping
- ✅ Payment holds for pending transactions
- ✅ Exchange rate handling
- ✅ Idempotency support
