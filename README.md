# Bitcoin Notification Server

A Go-based notification service

## Features

- **Create Notifications**: Generate notifications with fields
- **Send Notifications**: Email notifications to specified recipients
- **List Notifications**: Retrieve notifications with optional state filtering
- **Delete Notifications**: Soft delete functionality for notifications


## Installation & Setup

### Prerequisites
- Go 1.19 or later
- Git

### Clone the Repository
```bash
git clone https://github.com/Tazril/notification-server.git
cd notification-server
```

### Install Dependencies
```bash
go mod download
```

### Run the Application
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### 1. Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "notification-server"
}
```

### 2. Create Notification
```
POST /notifications
```

**Request Body:**
```json
{
  "current_btc_price": 45000.50,
  "market_trade_volume": 1500000000.75,
  "intra_day_high_price": 46200.00,
  "market_cap": 850000000000.25
}
```

**Response:**
```json
{
  "message": "Notification created successfully",
  "notification": {
    "ID": "uuid-string",
    "CurrentBTCPrice": 45000.50,
    "MarketTradeVolume": 1500000000.75,
    "IntraDayHighPrice": 46200.00,
    "MarketCap": 850000000000.25,
    "State": "CREATED",
    "Active": true,
    "CreatedAt": "2024-01-15T10:30:00Z",
    "UpdatedAt": "2024-01-15T10:30:00Z"
  }
}
```

### 3. Send Notification
```
POST /notifications/{id}/send
```

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "Notification sent successfully"
}
```

### 4. Get All Notifications
```
GET /notifications
```

**Optional Query Parameters:**
- `state`: Filter by notification state (`CREATED`, `SENT`, `FAILED`)

**Examples:**
- `GET /notifications` - Get all active notifications
- `GET /notifications?state=SENT` - Get only sent notifications

**Response:**
```json
{
  "notifications": [
    {
      "ID": "uuid-string",
      "CurrentBTCPrice": 45000.50,
      "MarketTradeVolume": 1500000000.75,
      "IntraDayHighPrice": 46200.00,
      "MarketCap": 850000000000.25,
      "State": "SENT",
      "Active": true,
      "CreatedAt": "2024-01-15T10:30:00Z",
      "UpdatedAt": "2024-01-15T10:35:00Z"
    }
  ],
  "count": 1
}
```

### 5. Delete Notification
```
DELETE /notifications/{id}
```

**Response:**
```json
{
  "message": "Notification deleted successfully"
}
```

## Example Usage

### Create and Send a Notification

1. **Create a notification:**
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "current_btc_price": 45000.50,
    "market_trade_volume": 1500000000.75,
    "intra_day_high_price": 46200.00,
    "market_cap": 850000000000.25
  }'
```

2. **Send the notification (use the ID from step 1):**
```bash
curl -X POST http://localhost:8080/notifications/{notification-id}/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

3. **Check notification status:**
```bash
curl http://localhost:8080/notifications
```

## Notification States

- **CREATED**: Notification has been created but not yet sent
- **SENT**: Notification has been successfully sent to the recipient
- **FAILED**: Notification sending failed

## Data Models

### Notification
```go
type Notification struct {
    ID                string
    CurrentBTCPrice   float64
    MarketTradeVolume float64
    IntraDayHighPrice float64
    MarketCap         float64
    State             NotificationState // CREATED, SENT, FAILED
    Active            bool
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

### NotificationState
```go
type NotificationState string

const (
    CREATED NotificationState = "CREATED"
    SENT    NotificationState = "SENT"
    FAILED  NotificationState = "FAILED"
)
```
