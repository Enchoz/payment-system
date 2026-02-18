# Postman Collection Setup Guide

This guide explains how to import and use the Postman collection for testing the Payment Processing System API.

## Importing the Collection

1. Open Postman
2. Click **Import** button (top left)
3. Select **File** tab
4. Choose `postman_collection.json` from this directory
5. Click **Import**

## Environment Variables

The collection uses the following variables that you can customize:

- `base_url`: Base URL for the API (default: `http://localhost:4078`)
- `from_user_id`: UUID of the sender user
- `to_user_id`: UUID of the recipient user (for internal payments)
- `external_account_id`: UUID of a saved external account (optional)
- `payment_request_id`: UUID of a payment request to process (set after creating external payment request)

### Setting Variables

1. Click on the collection name
2. Go to **Variables** tab
3. Update the values as needed
4. Variables can also be set at the environment level or globally

## Testing Workflow

### 1. Health Check
First, verify the server is running:
- Run **Health Check** request
- Should return `{"status": "ok"}`

### 2. Internal Payment
Test immediate payment between users:
- Update `from_user_id` and `to_user_id` variables with actual user UUIDs
- Run **Process Internal Payment** request
- Payment is processed immediately

### 3. External Payment Request
Test external payment with validation:
- Update `from_user_id` with actual user UUID
- Run **Create External Payment Request** request
- Copy the `id` from the response
- Set `payment_request_id` variable to this ID

### 4. Process Payment Request
After validation is complete:
- Run **Process Payment Request** request
- This will execute the payment

## Example Scenarios

### Scenario 1: Internal Payment (Same Currency)
```json
{
  "fromUserId": "11111111-1111-1111-1111-111111111111",
  "toUserId": "22222222-2222-2222-2222-222222222222",
  "amount": "100.00",
  "currency": "USD",
  "reference": "payment-001"
}
```

### Scenario 2: External Payment with Currency Conversion
```json
{
  "fromUserId": "11111111-1111-1111-1111-111111111111",
  "amount": "500.00",
  "currency": "USD",
  "targetCurrency": "EUR",
  "externalIban": "FR1420041010050500013M02606",
  "externalBankName": "French Bank",
  "externalAccountHolderName": "Jean Dupont",
  "externalCountryCode": "FR",
  "reference": "payment-002",
  "saveAccount": true
}
```

### Scenario 3: External Payment Using Saved Account
```json
{
  "fromUserId": "11111111-1111-1111-1111-111111111111",
  "externalAccountId": "550e8400-e29b-41d4-a716-446655440002",
  "amount": "200.00",
  "currency": "USD",
  "targetCurrency": "GBP",
  "reference": "payment-003"
}
```

## Expected Responses

### Internal Payment Response (200 OK)
```json
{
  "id": "uuid",
  "paymentReference": "internal-payment-1234567890",
  "fromUserId": "uuid",
  "toUserId": "uuid",
  "amount": "100.00",
  "currency": "USD",
  "status": "COMPLETED",
  "paymentType": "INTERNAL",
  "createdAt": "2024-01-15T10:30:00Z",
  "completedAt": "2024-01-15T10:30:00Z"
}
```

### External Payment Request Response (201 Created)
```json
{
  "id": "uuid",
  "requestReference": "external-payment-1234567890",
  "fromUserId": "uuid",
  "amount": "500.00",
  "currency": "USD",
  "targetCurrency": "EUR",
  "status": "VALIDATED",
  "paymentType": "EXTERNAL",
  "validationResults": {
    "ibanValidation": {
      "checkName": "IBAN Validation",
      "passed": true,
      "message": "Valid IBAN format"
    }
  },
  "exchangeRate": "0.92",
  "convertedAmount": "460.00",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "amount must be greater than zero"
}
```

### 500 Internal Server Error
```json
{
  "error": "insufficient balance"
}
```

## Tips

1. **Idempotency**: Use unique `idempotencyKey` values to prevent duplicate payments
2. **References**: Use unique `reference` values for tracking payments
3. **Save Account**: Set `saveAccount: true` to cache validated external accounts
4. **Currency Codes**: Use ISO 4217 currency codes (USD, EUR, GBP)
5. **UUIDs**: Ensure user IDs exist in the database before testing

## Troubleshooting

- **Connection Refused**: Make sure the server is running on the configured port
- **User Not Found**: Create users in the database first
- **Insufficient Balance**: Ensure the sender account has sufficient balance
- **Validation Failed**: Check that external account details are valid (IBAN format, etc.)
