# Gofermmart Loyalty System API

Simple HTTP API for loyalty system with following endpoints:

## Auth

### POST /api/user/register
Request body:
{
    "login": "string",     // Unique user identifier
    "password": "string"   // User password
}
Arguments:
- login: required field, cannot be empty
- password: required field, cannot be empty

Response: JWT token for authentication

### POST /api/user/login
Request body:
{
    "login": "string",     // Existing user login
    "password": "string"   // User password
}
Arguments:
- login: must match existing user
- password: must match registration password

Response: JWT token for authentication

## Orders 

### POST /api/user/orders
Content-Type: text/plain
Body: order number (string)

Arguments:
- order number: must pass Luhn algorithm validation

### GET /api/user/orders
Response:
[
    {
        "number": "string",      // Order number (Luhn algorithm)
        "status": "string",      // Order processing status
        "accrual": float,        // Awarded points (can be null)
        "uploaded_at": "string"  // Upload time in RFC3339 format
    }
]

Possible statuses:
- NEW: order uploaded to system
- PROCESSING: calculating points
- INVALID: order rejected
- PROCESSED: points calculation completed

## Balance

### GET /api/user/balance
Response:
{
    "current": float,    // Current points balance
    "withdrawn": float   // Total points used
}

### POST /api/user/balance/withdraw
Request:
{
    "order": "string",   // Order number for withdrawal
    "sum": float         // Amount to withdraw
}

Arguments:
- order: must pass Luhn algorithm validation
- sum: must be greater than 0 and not exceed current balance

### GET /api/user/withdrawals
Response:
[
    {
        "order": "string",      // Withdrawal order number
        "sum": float,           // Withdrawn amount
        "processed_at": string  // Withdrawal time in RFC3339 format
    }
]

## Response Codes
- 200: Successful operation
- 202: Accepted for processing
- 400: Invalid request format
- 401: Authentication required
- 409: Conflict (e.g. duplicate order)
- 422: Invalid order number format
- 500: Internal server error

## Authentication
All protected endpoints require header:
Authorization: Bearer <token>