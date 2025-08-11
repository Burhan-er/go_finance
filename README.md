# Go Finance API

Go Finance is a financial API service that provides user management, transaction handling, and balance tracking.
It includes integration with PostgreSQL, Redis, Prometheus, and Grafana for data storage, caching, and monitoring.

This project is developed within Insider using the Go (Golang) programming language.
---

## 🚀 Setup & Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/username/go-finance.git
   cd go-finance
````

2. **Start all services using Docker Compose:**

   ```bash
   docker-compose up --build
   ```

3. **Access the services once they are running:**

   * **API**: [http://localhost:8080](http://localhost:8080)
   * **PostgreSQL**: `localhost:5432`

     * **Database**: `mydb`
     * **User**: `myuser`
     * **Password**: `mypassword`
   * **Redis**: `localhost:6379`
   * **Prometheus**: [http://localhost:9090](http://localhost:9090)
   * **Grafana**: [http://localhost:3000](http://localhost:3000)

     * **Username**: `admin`
     * **Password**: `admin`

---

## 📡 API Endpoints & Examples

### **Authentication**

| Method | Endpoint                | Description          |
| ------ | ----------------------- | -------------------- |
| POST   | `/api/v1/auth/register` | Register a new user  |
| POST   | `/api/v1/auth/login`    | Authenticate user    |
| POST   | `/api/v1/auth/refresh`  | Refresh access token |

### **User Management**

| Method | Endpoint             | Description               |
| ------ | -------------------- | ------------------------- |
| GET    | `/api/v1/users`      | Retrieve all users        |
| GET    | `/api/v1/users/{id}` | Get a specific user by ID |
| PUT    | `/api/v1/users/{id}` | Update user details       |
| DELETE | `/api/v1/users/{id}` | Remove a user             |

### **Transaction**

| Method | Endpoint                        | Description                           |
| ------ | ------------------------------- | ------------------------------------- |
| POST   | `/api/v1/transactions/credit`   | Add funds to a user's account         |
| POST   | `/api/v1/transactions/debit`    | Deduct funds from a user's account    |
| POST   | `/api/v1/transactions/transfer` | Transfer funds between users          |
| GET    | `/api/v1/transactions/history`  | Get transaction history               |
| GET    | `/api/v1/transactions/{id}`     | Get details of a specific transaction |

### **Balance**

| Method | Endpoint                      | Description                         |
| ------ | ----------------------------- | ----------------------------------- |
| GET    | `/api/v1/balances/current`    | Get current balance                 |
| GET    | `/api/v1/balances/historical` | Retrieve historical balances        |
| GET    | `/api/v1/balances/at-time`    | Get balance at a specific timestamp |


#### Register

**Request**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "123456"
}
```

**Response**

```json
{
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-08-11T12:34:56Z",
    "updated_at": "2025-08-11T12:34:56Z"
  },
  "message": "User registered successfully"
}
```

#### Login

**Request**

```json
{
  "email": "john@example.com",
  "password": "123456"
}
```

**Response**

```json
{
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-08-11T12:34:56Z",
    "updated_at": "2025-08-11T12:34:56Z"
  },
  "access_token": "jwt_token_here",
  "message": "Login successful"
}
```

---

### **User Management**

#### Get All Users — `GET /api/v1/users`

**Response**

```json
{
  "users": [
    {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-08-11T12:34:56Z",
      "updated_at": "2025-08-11T12:34:56Z"
    }
  ]
}
```

#### Get User by ID — `GET /api/v1/users/{id}`

**Response**

```json
{
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-08-11T12:34:56Z",
    "updated_at": "2025-08-11T12:34:56Z"
  }
}
```

#### Update User — `PUT /api/v1/users/{id}`

**Request**

```json
{
  "username": "new_username",
  "email": "new_email@example.com"
}
```

**Response**

```json
{
  "user": {
    "id": "uuid",
    "username": "new_username",
    "email": "new_email@example.com",
    "role": "user",
    "created_at": "2025-08-11T12:34:56Z",
    "updated_at": "2025-08-11T12:40:00Z"
  },
  "message": "User updated successfully"
}
```

#### Delete User — `DELETE /api/v1/users/{id}`

**Response**

```json
{
  "message": "User deleted successfully"
}
```

---

### **Transactions**

#### Credit Transaction — `POST /api/v1/transactions/credit`

**Request**

```json
{
  "to_user_id": "uuid",
  "from_user_id": "uuid",
  "type": "credit",
  "amount": "100.00",
  "description": "Salary"
}
```

**Response**

```json
{
  "transaction": {
    "id": "uuid",
    "user_id": "uuid",
    "to_user_id": "uuid",
    "type": "credit",
    "status": "completed",
    "amount": "100.00",
    "created_at": "2025-08-11T12:45:00Z"
  },
  "message": "Transaction successful"
}
```

#### Debit Transaction — `POST /api/v1/transactions/debit`

**Request**

```json
{
  "to_user_id": "uuid",
  "from_user_id": "uuid",
  "type": "debit",
  "amount": "50.00",
  "description": "Shopping"
}
```

**Response**

```json
{
  "transaction": {
    "id": "uuid",
    "user_id": "uuid",
    "to_user_id": "uuid",
    "type": "debit",
    "status": "completed",
    "amount": "50.00",
    "created_at": "2025-08-11T12:50:00Z"
  },
  "message": "Transaction successful"
}
```

#### Transfer Transaction — `POST /api/v1/transactions/transfer`

**Request**

```json
{
  "from_user_id": "uuid",
  "to_user_id": "uuid",
  "amount": "25.00",
  "description": "Friend payment"
}
```

**Response**

```json
{
  "transaction": {
    "id": "uuid",
    "user_id": "uuid",
    "to_user_id": "uuid",
    "type": "transfer",
    "status": "completed",
    "amount": "25.00",
    "created_at": "2025-08-11T13:00:00Z"
  },
  "message": "Transfer completed successfully"
}
```

---

### **Balances**

#### Current Balance — `GET /api/v1/balances/current`

**Response**

```json
{
  "balance": {
    "user_id": "uuid",
    "amount": "200.00",
    "last_updated_at": "2025-08-11T13:05:00Z"
  }
}
```

#### Historical Balance — `GET /api/v1/balances/historical`

**Request**

```json
{
  "user_id": "uuid",
  "start_date": "2025-08-01T00:00:00Z",
  "end_date": "2025-08-11T00:00:00Z"
}
```

**Response**

```json
{
  "history": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "amount": "150.00",
      "recorded_at": "2025-08-05T00:00:00Z"
    },
    {
      "id": "uuid",
      "user_id": "uuid",
      "amount": "200.00",
      "recorded_at": "2025-08-10T00:00:00Z"
    }
  ]
}
```

#### Balance at Specific Time — `GET /api/v1/balances/at-time`

**Request**

```json
{
  "user_id": "uuid",
  "timestamp": "2025-08-09T12:00:00Z"
}
```

**Response**

```json
{
  "balance": {
    "user_id": "uuid",
    "amount": "175.00",
    "last_updated_at": "2025-08-09T12:00:00Z"
  }
}
```

---

## 🛠 Technologies Used

* **Go** — Backend API service
* **PostgreSQL** — Relational database
* **Redis** — In-memory cache and message broker
* **Prometheus** — Metrics collection
* **Grafana** — Metrics visualization and dashboard
* **Docker & Docker Compose** — Containerized deployment

---

## 📊 Monitoring

* **Prometheus** is available at: [http://localhost:9090](http://localhost:9090)
* **Grafana** is available at: [http://localhost:3000](http://localhost:3000)

  ```
  Username: admin
  Password: admin
  ```
