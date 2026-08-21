DomainHub 🚀

DomainHub is a full-stack Domain & Hosting Management application for managing domain records through a RESTful Go backend and a responsive vanilla HTML/CSS/JavaScript frontend.

The project demonstrates a layered Go backend architecture using raw PostgreSQL queries, database migrations, JWT authentication, middleware, DTOs, Swagger/OpenAPI documentation, unit tests, and production deployment on Render.

🌐 Live Application

Component

URL

Frontend

https://domain-hub-frontend.onrender.com/

Backend API

https://domain-hub-b9o2.onrender.com

Swagger UI

https://domain-hub-b9o2.onrender.com/swagger/index.html

The backend root path / is not an application page. Use /health, /domains, /auth/login, or Swagger for API access.

✨ Features

Domain Management

Create a domain

View all domains

View a domain by ID

Update a domain

Delete a domain

Filter domains by status

Search domains by registrar

Pagination with page, limit, total, and total pages

Authentication

User registration

Secure password hashing with bcrypt

Login with email and password

JWT token generation

JWT authentication middleware

Protected domain write operations

Automatic 401 Unauthorized handling for invalid/missing tokens

Backend Engineering

Layered architecture

Raw PostgreSQL queries — no GORM or sqlboiler

DTO layer for request/response separation

Logger middleware

Recovery middleware

CORS middleware

Graceful server shutdown

Database migrations

Swagger/OpenAPI documentation

Unit tests for Handler, Service, and Repository layers

🛠️ Tech Stack

Backend

**Go (Golang)"

Chi Router

PostgreSQL

database/sql with raw SQL queries

JWT (github.com/golang-jwt/jwt/v5)

bcrypt (golang.org/x/crypto/bcrypt)

Swagger / Swaggo

Go testing + sqlmock

Database Migrations

Frontend

HTML5

CSS3

Vanilla JavaScript

Tailwind CSS CDN for the UI styling

Deployment

Render Web Service — Go backend

Render PostgreSQL — production database

Render Static Site — frontend

GitHub — source control

🏗️ Architecture

                         DomainHub
                            │
              ┌─────────────┴─────────────┐
              │                           │
         Frontend                     Backend
              │                           │
   HTML / CSS / JavaScript          Chi Router
              │                           │
              │                     Middleware
              │                  ┌────────┼────────┐
              │                  │        │        │
              │                Logger  Recovery   CORS
              │                           │
              │                     Auth Middleware
              │                           │
              │                       Handlers
              │                           │
              │                        Services
              │                           │
              │                       Repository
              │                           │
              │                    Raw SQL Queries
              │                           │
              └──────────── REST API ─────┤
                                          │
                                      PostgreSQL

📁 Project Structure

domainhub/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── database.go
│   │
│   ├── dto/
│   │   ├── auth.go
│   │   ├── domain.go
│   │   └── response.go
│   │
│   ├── handlers/
│   │   ├── auth_handlers.go
│   │   ├── domain_handler.go
│   │   └── domain_handler_test.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   │
│   ├── models/
│   │   ├── domain.go
│   │   ├── response.go
│   │   └── user.go
│   │
│   ├── repository/
│   │   ├── domain_repository.go
│   │   ├── domain_repository_test.go
│   │   └── user_repository.go
│   │
│   ├── routes/
│   │   └── routes.go
│   │
│   ├── service/
│   │   ├── domain_service.go
│   │   ├── domain_service_test.go
│   │   └── auth_service.go
│   │
│   └── utils/
│       └── response.go
│
├── migrations/
│   ├── 001_create_domains_table.up.sql
│   ├── 001_create_domains_table.down.sql
│   ├── 002_create_users_table.up.sql
│   └── 002_create_users_table.down.sql
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── frontend/
│   ├── login.html
│   ├── register.html
│   ├── dashboard.html
│   ├── domains.html
│   ├── manage-domains.html
│   ├── api.js
│   └── config.js
│
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md

.env is for local development only and must not be committed to GitHub.

🗄️ Database Schema

domains

Column

Type

Description

id

serial

Primary key

domain_name

varchar(200)

Unique domain name

registrar

varchar(100)

Domain registrar

expiry_date

date

Domain expiry date

status

varchar(20)

ACTIVE or EXPIRED

users

Column

Type

Description

id

serial

Primary key

email

varchar(150)

Unique user email

password_hash

varchar(255)

bcrypt password hash

role

varchar(20)

Stored user role

created_at

timestamp

User creation time

The project currently uses the role field as part of the user model/JWT, but authorization rules are intentionally not implemented. Authentication is the focus.

🔐 Authentication Flow

POST /auth/register
        │
        ▼
Validate input
        │
        ▼
bcrypt hash password
        │
        ▼
Store user in PostgreSQL

Login:

POST /auth/login
        │
        ▼
Find user by email
        │
        ▼
bcrypt.CompareHashAndPassword()
        │
        ▼
Generate JWT
        │
        ▼
Return token

Protected request:

Authorization: Bearer <JWT>
        │
        ▼
JWT Middleware
        │
        ├── Invalid → 401 Unauthorized
        │
        └── Valid → Continue to handler

Protected APIs

POST   /domains
PUT    /domains/{id}
DELETE /domains/{id}

The read APIs remain public:

GET /domains
GET /domains/{id}

📡 API Endpoints

Health

GET /health

Authentication

Register

POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "StrongPassword123"
}

Login

POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "StrongPassword123"
}

Successful login returns a JWT token under data.token.

Domains

Get all domains

GET /domains

Supports:

page
limit
status
registrar

Examples:

/domains?page=1&limit=10
/domains?status=ACTIVE
/domains?registrar=Cloudflare
/domains?page=2&limit=5&status=ACTIVE&registrar=Cloudflare

Get domain by ID

GET /domains/{id}

Create domain

POST /domains
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "domain_name": "example.com",
  "registrar": "Cloudflare",
  "expiry_date": "2030-01-01T00:00:00Z",
  "status": "ACTIVE"
}

Update domain

PUT /domains/{id}
Authorization: Bearer <JWT>
Content-Type: application/json

Delete domain

DELETE /domains/{id}
Authorization: Bearer <JWT>

📄 Standard API Response

The backend uses a standardized response structure:

{
  "success": true,
  "message": "Domains fetched successfully",
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 0,
    "total_pages": 0
  }
}

🧩 Middleware

Logger Middleware

Records:

HTTP method

Request path

Response status

Request duration

Example:

GET /domains/4 200 27.115194ms

Recovery Middleware

Recovers from panics so that a single request does not crash the whole server.

CORS Middleware

Allows the deployed frontend to communicate with the Go backend and supports the Authorization header used by JWT authentication.

Auth Middleware

Validates the JWT sent through:

Authorization: Bearer <token>

🧪 Testing

The project contains unit tests for the important application layers.

Run all tests:

go test ./...

Run repository tests:

go test ./internal/repository -v

Run handler tests:

go test ./internal/handlers

Run service tests:

go test ./internal/service

The repository tests use SQL mocking to verify raw SQL operations without requiring the real database for every test.

🗃️ Database Migrations

The migration history is:

001 → create_domains_table
002 → create_users_table

Check migration version:

migrate -path ./migrations \
-database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" \
version

Apply migrations:

migrate -path ./migrations \
-database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" \
up

Rollback the latest migration:

migrate -path ./migrations \
-database "postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require" \
steps -1

📚 Swagger / OpenAPI

Swagger documentation is generated using Swaggo.

Generate documentation:

swag init -g cmd/server/main.go

Generated files are stored in:

docs/
├── docs.go
├── swagger.json
└── swagger.yaml

Swagger UI:

https://domain-hub-b9o2.onrender.com/swagger/index.html

▶️ Run Locally

1. Clone the repository

git clone <your-github-repository-url>
cd domainhub