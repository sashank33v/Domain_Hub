# DomainHub 🚀

DomainHub is a full-stack Domain & Hosting Management application. A RESTful Go backend manages domain records in PostgreSQL, publishes domain events to Kafka, and serves a responsive vanilla HTML/CSS/JavaScript frontend. It also includes an AI assistant that answers questions about your stored domains.

The project demonstrates a layered Go backend architecture using raw PostgreSQL queries, database migrations, JWT authentication, middleware, DTOs, Kafka event streaming, an LLM-powered chatbot (Groq), Swagger/OpenAPI documentation, unit tests, and production deployment on Render.

## 🌐 Live Application

| Component | URL |
|-----------|-----|
| App link | https://domain-hub-b9o2.onrender.com |
| Swagger UI | https://domain-hub-b9o2.onrender.com/swagger/index.html |

The Go server serves both the API and the frontend. Opening the App link loads the login page.

---

## ✨ Features

### Domain Management
- Create, view, update, and delete domains
- View a domain by ID
- Filter domains by status
- Search domains by registrar
- Pagination with `page`, `limit`, `total`, and `total_pages`
- Dashboard with total, active, expiring-soon, and expired counts

### Authentication
- User registration with bcrypt password hashing
- Login with email and password
- JWT token generation
- JWT authentication middleware
- Protected write operations and chat endpoint
- Automatic redirect to login on `401 Unauthorized`

### AI Assistant 🤖
- Chat widget on the dashboard (💬 button)
- Answers **only** from the domains stored in DomainHub (names, registrars, expiry dates, status, counts)
- Uses live database data on every question
- `days_left` for each domain is calculated by the backend, so the model does not do date math
- Refuses questions that are not about your domains
- Remembers the last 6 messages for follow-up questions
- Powered by Groq (default model: `openai/gpt-oss-20b`)

### Event Streaming
- Domain create, update, and delete publish events to Kafka (`domain.created`, `domain.updated`, `domain.deleted`)
- Producer with retry and backoff
- Consumer group that reads the events in the same application
- Kafka broker for local development via `docker-compose.yaml`

### Backend Engineering
- Layered architecture (handler → service → repository)
- Raw PostgreSQL queries (no ORM)
- DTO layer for request/response separation
- Logger and Recovery middleware
- Graceful server shutdown
- Database migrations
- Swagger/OpenAPI documentation
- Unit tests for handler, service, and repository layers
- CLI tool to reset a user password

---

## 🛠️ Tech Stack

### Backend
- Go (Golang)
- Chi Router
- PostgreSQL with `database/sql` and pgx
- JWT (`github.com/golang-jwt/jwt/v5`)
- bcrypt (`golang.org/x/crypto/bcrypt`)
- Kafka (`github.com/segmentio/kafka-go`)
- Groq API (OpenAI-compatible chat completions)
- Swagger / Swaggo
- Go testing + sqlmock

### Frontend
- HTML5, CSS3, Vanilla JavaScript
- Tailwind CSS (CDN)

### Deployment
- Render Web Service (Go backend + frontend)
- Render PostgreSQL
- GitHub

---

## 🏗️ Architecture

```
                         DomainHub
                             │
                  Browser (HTML / CSS / JS)
                             │  REST API (JWT)
                             ▼
                  Go server (Chi Router)
                             │
                 Logger → Recovery → Auth
                             │
                         Handlers
                             │
                         Services ──────────► Groq API (chatbot)
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
         Repository                    Kafka producer
              │                             │
        Raw SQL queries               Kafka (domain-events)
              │                             │
          PostgreSQL                  Kafka consumer (logs events)
```

### AI chat flow

```
Question ─► POST /chat (JWT) ─► validate (max 500 chars)
        ─► load domains + counts from PostgreSQL
        ─► system prompt (rules + today's date + DATA)
        ─► Groq ─► reply ─► chat widget
```

---

## 📁 Project Structure

```
domainhub/
├── cmd/
│   ├── server/main.go              # application entrypoint
│   └── resetpassword/main.go       # CLI: reset a user's password
│
├── internal/
│   ├── config/config.go
│   ├── database/database.go
│   ├── dto/                        # auth, domain request/response/mapper
│   ├── handlers/                   # auth, domain, chat, health (+ tests)
│   ├── messaging/                  # kafka_producer.go, kafka_consumer.go
│   ├── middleware/                 # auth.go, cors.go, logger.go (Logger + Recovery)
│   ├── models/                     # domain, user, response, pagination
│   ├── repository/                 # domain, user (+ tests)
│   ├── routes/routes.go
│   ├── service/                    # auth, domain, ai_service (+ tests)
│   └── utils/response.go
│
├── migrations/                     # 001 domains, 002 users
├── docs/                           # generated Swagger files
├── frontend/
│   ├── login.html / register.html
│   ├── dashboard.html              # stats + AI chat widget
│   ├── domains.html / manage-domains.html
│   ├── api.js                      # API client + auth helpers
│   ├── chat.js                     # chat widget logic
│   └── config.js
│
├── docker-compose.yaml             # local Kafka broker
├── .env                            # local only, never commit
├── go.mod / go.sum
└── README.md
```

---

## 🗄️ Database Schema

### domains

| Column | Type | Description |
|--------|------|-------------|
| id | serial | Primary key |
| domain_name | varchar(200) | Unique domain name |
| registrar | varchar(100) | Domain registrar |
| expiry_date | date | Domain expiry date |
| status | varchar(20) | `ACTIVE` or `EXPIRED` |

### users

| Column | Type | Description |
|--------|------|-------------|
| id | serial | Primary key |
| email | varchar(150) | Unique user email |
| password_hash | varchar(255) | bcrypt password hash |
| role | varchar(20) | Stored user role (default `USER`) |
| created_at | timestamp | User creation time |

The `role` field is stored and included in the JWT, but role-based authorization is not implemented. Authentication is the focus.

---

## 🔐 Authentication Flow

```
POST /auth/register ─► validate ─► bcrypt hash ─► store user

POST /auth/login ─► find user ─► bcrypt compare ─► generate JWT ─► return token

Authorization: Bearer <JWT> ─► Auth middleware ─┬─ invalid → 401
                                                └─ valid   → handler
```

### Protected APIs

```
POST   /domains
PUT    /domains/{id}
DELETE /domains/{id}
POST   /chat
```

Public: `GET /health`, `GET /domains`, `GET /domains/{id}`, `/auth/*`, Swagger, and the frontend files.

---

## 📡 API Endpoints

### Health
```http
GET /health
```

### Authentication
```http
POST /auth/register
POST /auth/login
Content-Type: application/json

{ "email": "user@example.com", "password": "StrongPassword123" }
```
A successful login returns a JWT under `data.token`.

### Domains
```http
GET    /domains                  # page, limit, status, registrar
GET    /domains/{id}
POST   /domains                  # JWT
PUT    /domains/{id}             # JWT
DELETE /domains/{id}             # JWT
```

Examples:
```
/domains?page=1&limit=10
/domains?status=ACTIVE
/domains?registrar=Cloudflare
/domains?page=2&limit=5&status=ACTIVE&registrar=Cloudflare
```

Create body:
```json
{
  "domain_name": "example.com",
  "registrar": "Cloudflare",
  "expiry_date": "2030-01-01T00:00:00Z",
  "status": "ACTIVE"
}
```

### AI Chat
```http
POST /chat
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "message": "Which domains expire in the next 30 days?",
  "history": [
    { "role": "user", "content": "List domains registered with GoDaddy" },
    { "role": "assistant", "content": "..." }
  ]
}
```

Response:
```json
{
  "success": true,
  "message": "Chat response generated",
  "data": { "reply": "..." }
}
```

Rules:
- `message` is required, maximum 500 characters
- `history` is optional; only the last 6 `user`/`assistant` messages are used (other roles are ignored)
- The request body is limited to 32 KB
- Missing `GROQ_API_KEY` or a Groq failure returns an error response

### Standard API Response
```json
{
  "success": true,
  "message": "Domains fetched successfully",
  "data": [],
  "pagination": { "page": 1, "limit": 10, "total": 0, "total_pages": 0 }
}
```

---

## 🧩 Middleware

- **Logger**: method, path, status, and duration of each request
- **Recovery**: recovers from panics so one request does not crash the server
- **Auth**: validates `Authorization: Bearer <token>`
- **CORS**: `middleware/cors.go` exists but is not enabled, because the frontend is served by the same server (same origin)

---

## 🧪 Testing

```bash
go test ./...
go test ./internal/repository -v
go test ./internal/handlers
go test ./internal/service
```

- Repository tests use sqlmock to verify raw SQL without a real database
- Handler tests use a fake service
- Chat tests cover context building (including `days_left`), message validation, and bad requests

---

## 🗃️ Database Migrations

```
001 → create_domains_table
002 → create_users_table
```

```bash
# check version
migrate -path ./migrations -database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" version

# apply
migrate -path ./migrations -database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" up

# rollback latest
migrate -path ./migrations -database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" steps -1
```

---

## 📚 Swagger / OpenAPI

```bash
swag init -g cmd/server/main.go
```

Generated files are stored in `docs/`. Swagger UI: https://domain-hub-b9o2.onrender.com/swagger/index.html

---

## ▶️ Run Locally

### 1. Clone
```bash
git clone <your-github-repository-url>
cd domainhub
```

### 2. Environment variables
Create `.env` in the project root:

```env
APP_NAME=DomainHub
PORT=8000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=domainhub_db
DB_SSLMODE=disable

JWT_SECRET=your-local-jwt-secret

# Kafka
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=domain-events

# AI assistant
GROQ_API_KEY=your-groq-api-key
GROQ_MODEL=openai/gpt-oss-20b
```

| Variable | Purpose |
|----------|---------|
| `PORT`, `APP_NAME` | Server settings |
| `DB_*` | PostgreSQL connection |
| `JWT_SECRET` | JWT signing key |
| `KAFKA_BROKER`, `KAFKA_TOPIC` | Kafka connection and topic |
| `GROQ_API_KEY` | Groq API key for the AI assistant |
| `GROQ_MODEL` | Optional; defaults to `openai/gpt-oss-20b` |

### 3. Install dependencies
```bash
go mod download
```

### 4. Start Kafka
```bash
docker compose up -d
```

### 5. Apply migrations
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/domainhub_db?sslmode=disable" up
```

### 6. Start the app
Run from the project root (the server reads `.env` and `./frontend` relative to it):

```bash
go run ./cmd/server
```

Open `http://localhost:8000` (the frontend and API are served together).

### Reset a user password
```bash
go run ./cmd/resetpassword
```

---

## ⚙️ Frontend Configuration

`frontend/config.js` selects the API base URL:

```js
const API_BASE_URL =
    window.location.hostname === "localhost"
        ? window.location.origin
        : "https://domain-hub-b9o2.onrender.com";
```

- On `localhost` the frontend calls the local server
- Elsewhere it calls the production server
- JWT tokens are stored in browser local storage and attached to protected requests (all non-GET requests except `/auth/*`)

---

## 🚀 Deployment

- **Web Service**: the Go app on Render serves the API and the frontend
- **Database**: Render PostgreSQL; migrations are applied before running the API
- **Environment variables** (set in Render): `APP_NAME`, `PORT`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`, `JWT_SECRET`, `GROQ_API_KEY`, `GROQ_MODEL`, and `KAFKA_BROKER` / `KAFKA_TOPIC` if a Kafka broker is available

```
User Browser ─► Render Web Service (Go: API + frontend) ─┬─► Render PostgreSQL
                                                          ├─► Groq API
                                                          └─► Kafka (if configured)
```

Note: domain create, update, and delete publish a Kafka event after the database change. If the broker is unreachable, the API returns an error even though the database change was saved.

---

## 🔒 Security Notes

- Passwords are stored as bcrypt hashes and excluded from JSON responses
- `.env` must not be committed; `JWT_SECRET` and `GROQ_API_KEY` must be deployment secrets
- Write APIs and `/chat` require a valid JWT
- The AI assistant only receives domain data, never user accounts
- User-entered domain fields are sanitized before being sent to the model, and chat history roles other than `user`/`assistant` are ignored
- Chat replies are HTML-escaped in the browser before formatting

---

## 🔄 Graceful Shutdown

```
Ctrl+C / SIGTERM ─► stop accepting new work ─► finish active requests
                 ─► shut down HTTP server ─► stop Kafka consumer ─► close database
```

---

## 📌 Current Project Status

- [x] REST API CRUD
- [x] Raw PostgreSQL integration
- [x] Database migrations
- [x] Pagination, status and registrar filtering
- [x] DTO layer and standard API responses
- [x] Logger and Recovery middleware
- [x] Swagger/OpenAPI documentation
- [x] Unit tests
- [x] Graceful shutdown
- [x] User registration, bcrypt hashing, JWT login
- [x] JWT middleware and protected write APIs
- [x] Kafka producer and consumer
- [x] AI assistant (`POST /chat`) with live domain data and chat history
- [x] Dashboard chat widget
- [x] Production PostgreSQL and deployed Go backend serving the frontend

---

## 👨‍💻 Project Focus

DomainHub was built as a practical backend-focused project to demonstrate:

- Go REST API development and layered architecture
- Raw SQL and PostgreSQL
- JWT and bcrypt authentication
- Middleware design
- Kafka event streaming
- Grounded LLM integration (Groq) with prompt scope control
- Database migrations, automated testing, and API documentation
- Production deployment

## 📄 License

This project is for learning, portfolio, and demonstration purposes.