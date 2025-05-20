# Gofermmart Loyalty System API

Simple HTTP API for loyalty system with following endpoints:

## Auth
```
POST /api/user/register - Register user
POST /api/user/login - Login user
```

## Orders 
```
POST /api/user/orders - Upload order number
GET /api/user/orders - Get user orders
```

## Balance
```
GET /api/user/balance - Get balance
POST /api/user/balance/withdraw - Withdraw points
GET /api/user/withdrawals - Get withdrawals history
```

## Features
- User registration and auth with JWT
- Order processing and tracking
- Loyalty points calculation
- Balance management
- Withdrawal system

All authenticated endpoints require: `Bearer <token>`

Status codes: 200 (Success), 202 (Accepted), 400 (Bad Request), 401 (Unauthorized), 409 (Conflict), 422 (Invalid Order), 500 (Error)