# ExpenseTracker Backend API Documentation

Welcome to the production-grade ExpenseTracker API documentation. This API is built with a mobile-first, clean architecture, featuring a secure **AWS S3 Client-Direct Upload flow via Presigned URLs**, robust ownership validation, and index-optimized cursor pagination.

---

## 1. Authentication Requirements

All protected endpoints require a JWT Access Token.
- **Header format**: `Authorization: Bearer <access_token>`
- **Token Injection**: The token is validated server-side by the `AuthMiddleware`. It automatically extracts the `user_id` and injects it into the request context. 
- **Security Guard**: Handlers never trust a `user_id` provided in the request body; it is strictly parsed from the secure request context.

---

## 2. Route List

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| **POST** | `/api/v1/auth/register` | Register a new user | No |
| **POST** | `/api/v1/auth/login` | Login and get access/refresh tokens | No |
| **POST** | `/api/v1/auth/refresh` | Refresh access token | No |
| **POST** | `/api/v1/auth/verify-email` | Verify user email | No |
| **POST** | `/api/v1/auth/forgot-password` | Request password reset email | No |
| **POST** | `/api/v1/auth/reset-password` | Reset password using token | No |
| **POST** | `/api/v1/auth/logout` | Logout user | Yes |
| **GET** | `/api/v1/users/me` | Fetch active user profile | Yes |
| **PATCH** | `/api/v1/users/me` | Update name/bio/photo_url (No file uploads) | Yes |
| **POST** | `/api/v1/users/me/profile-picture/presign` | Generate AWS S3 upload pre-signed PUT URL | Yes |
| **POST** | `/api/v1/users/me/profile-picture/complete` | Finalize profile picture DB metadata & delete old photo | Yes |
| **POST** | `/api/v1/transactions` | Create a new transaction | Yes |
| **GET** | `/api/v1/transactions` | List user transactions with cursor pagination & filtering | Yes |
| **GET** | `/api/v1/transactions/info` | Fetch transaction statistics & category breakdown | Yes |
| **GET** | `/api/v1/transactions/{id}` | Fetch a specific transaction by ID | Yes |
| **PATCH** | `/api/v1/transactions/{id}` | Update an existing transaction | Yes |
| **PATCH** | `/api/v1/transactions/{id}/status` | Soft delete a transaction | Yes |

---

## 3. Standard Headers

All JSON requests must include the following header:
```http
Content-Type: application/json
```
All protected endpoints must include:
```http
Authorization: Bearer <jwt_access_token>
```

---

## 4. Query Parameters (GET /api/v1/transactions & GET /api/v1/transactions/info)

| Parameter | Type | Required | Description | Default |
| :--- | :--- | :--- | :--- | :--- |
| `category` | `string` / `int` | No | Filter by transaction category (1 - 21). Supports single ID or comma-separated list (e.g., `1,2,3`) | - |
| `category_ids` | `string` | No | Comma-separated list of category IDs (e.g., `1,2,3`) for multi-category filter | - |
| `type` | `int` | No | Filter by transaction type (1: Expense, 2: Income) | - |
| `status` | `int` | No | Filter by status (1: Active, 2: Deleted) | 1 (Active) |
| `from_date` | `string` | No | Start date filter (RFC3339 format) | - |
| `to_date` | `string` | No | End date filter (RFC3339 format) | - |
| `after_id` | `int` | No | Cursor ID for pagination (**list only**) | - |
| `limit` | `int` | No | Max items returned (Min 1, Max 100) (**list only**) | 20 |
| `sort` | `string` | No | Sort type (`newest`, `amount_asc`, `amount_desc`) (**list only**) | `newest` |

---

## 5. Request Body Examples

### PATCH /api/v1/users/me
```json
{
  "name": "Jane Doe",
  "bio": "Personal finance enthusiast."
}
```

### POST /api/v1/users/me/profile-picture/presign
```json
{
  "content_type": "image/jpeg"
}
```

### POST /api/v1/users/me/profile-picture/complete
```json
{
  "object_key": "ExpenseTracker/Users/12/Profile/profile_1748039200.jpg"
}
```

### POST /api/v1/transactions
```json
{
  "amount": 25.50,
  "category": 5,
  "type": 1,
  "description": "Weekly grocery run at Market",
  "date": "2026-05-21T08:00:00Z"
}
```

### PATCH /api/v1/transactions/{id}
```json
{
  "amount": 28.00,
  "description": "Adjusted grocery run cost"
}
```

### PATCH /api/v1/transactions/{id}/status (Soft Delete)
```json
{
  "status": 2
}
```

---

## 6. Response Examples

### GET /api/v1/users/me
```json
{
  "id": 12,
  "email": "user@example.com",
  "name": "Jane Doe",
  "photo_url": "https://expense-tracker-bucket.s3.amazonaws.com/ExpenseTracker/Users/12/Profile/profile_1748039200.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...&X-Amz-Date=...&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=...",
  "photo_key": "ExpenseTracker/Users/12/Profile/profile_1748039200.jpg",
  "photo_size": 204850,
  "bio": "Personal finance enthusiast.",
  "is_email_verified": true,
  "is_transactions_enabled": true,
  "is_reporting_enabled": false,
  "is_notification_enabled": false,
  "created_at": "2026-05-19T04:16:41Z",
  "updated_at": "2026-05-21T08:30:00Z"
}
```

### POST /api/v1/users/me/profile-picture/presign
```json
{
  "upload_url": "https://expense-tracker-bucket.s3.amazonaws.com/ExpenseTracker/Users/12/Profile/profile_1748039200.jpg?AWSAccessKeyId=...",
  "object_key": "ExpenseTracker/Users/12/Profile/profile_1748039200.jpg"
}
```

### POST /api/v1/transactions
```json
{
  "id": 105,
  "amount": 25.5,
  "category": 5,
  "type": 1,
  "description": "Weekly grocery run at Market",
  "date": "2026-05-21T08:00:00Z",
  "status": 1
}
```

---

## 7. Error Response Examples

### 401 Unauthorized
```json
{
  "error": "unauthorized"
}
```

### 400 Bad Request (Validation Error)
```json
{
  "error": "invalid file type, allowed: image/jpeg, image/png, image/webp, image/gif"
}
```

### 404 Not Found (Access Denied / Not Found)
```json
{
  "error": "transaction not found or access denied"
}
```

---

## 8. Cursor Pagination Examples

To optimize cellular connections and support smooth FlatList rendering, pagination uses `after_id` instead of offsets.

### Initial Request:
`GET /api/v1/transactions?limit=2`

Response:
```json
{
  "items": [
    {
      "id": 102,
      "amount": 1500.00,
      "category": 15,
      "type": 2,
      "description": "Salary paycheck",
      "date": "2026-05-21T02:00:00Z",
      "status": 1
    },
    {
      "id": 99,
      "amount": 45.00,
      "category": 8,
      "type": 1,
      "description": "Uber ride",
      "date": "2026-05-20T18:30:00Z",
      "status": 1
    }
  ],
  "next_cursor": 99,
  "has_more": true,
  "limit": 2
}
```

### Next Request (Passing `next_cursor` as `after_id`):
`GET /api/v1/transactions?limit=2&after_id=99`

---

## 9. Status Code Table

| Status Code | Description | Occurrence |
| :--- | :--- | :--- |
| `200 OK` | Request completed successfully | Fetching / updating profiles, list queries |
| `201 Created` | Resource created successfully | Creating a transaction |
| `400 Bad Request` | Invalid inputs or failed validation | Exceeding S3 limits, invalid content-type |
| `401 Unauthorized` | Missing, invalid, or expired JWT | Request without valid Authorization header |
| `404 Not Found` | Resource not found or no ownership | Attempting to fetch/update another user's transaction |
| `500 Server Error` | Database connection errors or S3 panic | General backend issues |

---

## 10. File Upload Presign & Complete Architecture Flow

### Upload Workflow:
1. **Presign Request**: Call `POST /api/v1/users/me/profile-picture/presign` passing `"content_type": "image/jpeg"`.
2. **Direct PUT S3**: Perform an HTTP `PUT` request directly to the `upload_url` returned from step 1, attaching the binary image as the payload and locking `"Content-Type": "image/jpeg"` header.
3. **Complete Request**: Call `POST /api/v1/users/me/profile-picture/complete` passing `"object_key": "<returned_object_key>"` to verify the upload in S3 and finalize the database URL metadata update.
