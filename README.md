# go_finance
Insider şirketi bünyesinde, Go (Golang) programlama dili kullanılarak geliştirilecek bir finans uygulamasıdır.

Tamam, tek bir tam **README.md** dosyasını sana burada veriyorum.
İngilizce, lisans kısmı yok, verdiğin tüm servis bilgileri ve endpointler dahil:

---

````markdown
# Go Finance API

Go Finance is a financial API service that provides **user management**, **transaction handling**, and **balance tracking**.  
It comes with **PostgreSQL**, **Redis**, **Prometheus**, and **Grafana** integration for data storage, caching, and monitoring.

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

## 📡 API Endpoints

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
