
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

bash:

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
  
# TESTING THE ENDPOINTS 
 
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
  

### Testing on Cloud Run



bash: 

URL=https://library-api-72471622130.us-central1.run.app


### test wrapped on code : 
## download test_api_coll.py as test.py and run it to get the token and test the endpoints: 
curl -sSL https://raw.githubusercontent.com/kornnellio/library-api/main/test_api_coll.py -o test.py && python3 test.py

### Or :    

### test with curl 

# === CONFIG ===
URL="https://library-api-72471622130.us-central1.run.app"

# === 1. Get your Firebase token (ONE TIME) ===
echo "Open this → click Login → copy token:"
echo "https://kornnellio.github.io/library-api/test-login.html"
read -p "Paste your Firebase ID Token: " TOKEN
echo


# ===  Public endpoints (no token) ===
echo "Health check:"
curl -s $URL/health
echo -e "\n"

echo "Home page:"
curl -s $URL | head -c 100
echo -e "\n"

# === 3. Register a test user (optional) ===
echo "Register test user:"
curl -s -X POST $URL/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@kornnellio.ro","password":"123456"}' | jq
echo

# === 4. PROTECTED ENDPOINTS (with token) ===
echo "List your books:"
curl -s -H "Authorization: Bearer $TOKEN" $URL/books | jq
echo

echo "Add a book:"
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune","author":"Frank Herbert"}' \
  $URL/books | jq
echo

echo "List again (should show Dune):"
curl -s -H "Authorization: Bearer $TOKEN" $URL/books | jq
echo

# === 5. Update & Delete (get ID from previous output) ===
echo "Replace 123 with real book ID from above"
read -p "Book ID to update/delete: " BOOK_ID

curl -s -X PUT -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune (Updated)","author":"Herbert"}' \
  $URL/books/$BOOK_ID | jq

curl -s -X DELETE -H "Authorization: Bearer $TOKEN" \
  $URL/books/$BOOK_ID | jq

echo "Final list (should be empty):"
curl -s -H "Authorization: Bearer $TOKEN" $URL/books | jq
