
# Library API
A **Go + Gin** REST API for managing books.

- **PostgreSQL** (via `pgx`)
- **Docker**
- **GitHub Actions CI/CD**
- **Google Cloud Run**
- **Free-tier Cloud SQL**

---

## Features

| Endpoint | Method | Description |
|--------|--------|-----------|
| `/` | GET | Home page |
| `/register` | POST | Create user |
| `/login` | POST | Login (no JWT yet) |
| `/books` | GET | List books |
| `/books` | POST | Add book |
| `/books/:id` | PUT | Update book |
| `/books/:id` | DELETE | Delete book |
| `/health` | GET | Health check |


## Local Testing (Docker)

### 1. Start PostgreSQL + API

```bash
# Create network
docker network create library-net

# Start DB
docker run -d \
  --name library-db \
  --network library-net \
  -e POSTGRES_PASSWORD=Nel01@cor10 \
  -e POSTGRES_DB=library \
  -p 5433:5432 \
  postgres:15-alpine

# Wait
sleep 20

# Run API
docker run -p 8080:8080 \
  --network library-net \
  -e DATABASE_URL="postgres://postgres:Nel01%40cor10@library-db:5432/library?sslmode=disable" \
  library-api
```

### 2. Testing the Endpoints

```bash
# Home
curl http://localhost:8080

# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"123456"}'

# Add book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"1984","author":"George Orwell"}'

# List books
curl http://localhost:8080/books
```

## Testing on Cloud Run

```bash
URL=https://library-api-72471622130.us-central1.run.app

# 1. Health Check
curl -s $URL/health

# 2. Home
curl -s $URL

# 3. Register
curl -s -X POST $URL/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# 4. Login
curl -s -X POST $URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# 5. Get Books (empty)
curl -s $URL/books

# 6. Create Book
curl -s -X POST $URL/books \
  -H "Content-Type: application/json" \
  -d '{"title":"newBook","author":"newAuthor"}'

# 7. Get Books
curl -s $URL/books

# 8. Update Book with id = x
curl -s -X PUT $URL/books/x \
  -H "Content-Type: application/json" \
  -d '{"title":"newBookUpdated","author":"Orwell, George"}'

# 9. Delete Book with id = x
curl -s -X DELETE $URL/books/x
```

