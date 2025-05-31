# Targeting Engine

A high-performance microservice for personalized advertisement targeting built with Go, PostgreSQL, and Redis.

## 🚀 Features

- **Campaign Targeting**: Intelligent campaign delivery based on user attributes
- **Rule-based Engine**: Flexible targeting rules with INCLUDE/EXCLUDE logic
- **High Performance**: Redis caching for sub-millisecond response times
- **RESTful API**: Clean HTTP endpoints for campaign retrieval
- **Real-time Updates**: Dynamic targeting rule refresh without downtime

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Git

## 🛠️ Local Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd targeting-engine
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup PostgreSQL Database

#### Install PostgreSQL (if not already installed)

**macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

#### Create Database and User

```bash
# Connect to PostgreSQL
psql postgres

# Create database and user
CREATE DATABASE targeting_engine;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE targeting_engine TO postgres;
\q
```

### 4. Setup Redis

#### Install Redis (if not already installed)

**macOS:**
```bash
brew install redis
brew services start redis
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
```

### 5. Environment Configuration

Create a `.env` file in the project root:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=targeting_engine

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# Server Configuration
SERVER_PORT=7000
```

### 6. Initialize Database Schema

Run the database setup script:

```bash
go run scripts/setup_db.go
```

This will create the necessary tables and insert sample data.

### 7. Build and Run the Application

```bash
# Build the application
go build -o targeting-engine cmd/server/main.go

# Set environment variables and run
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=targeting_engine
export REDIS_HOST=localhost
export REDIS_PORT=6379
export SERVER_PORT=7000

# Run the application
./targeting-engine
```

The server will start on `http://localhost:7000`

## 📚 API Documentation

### Base URL
```
http://localhost:7000
```

### Endpoints

#### 1. Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "targeting-engine"
}
```

#### 2. Get Campaigns (GET)
```http
GET /api/v1/campaigns?app={app}&country={country}&os={os}&limit={limit}
```

**Parameters:**
- `app` (required): Application name
- `country` (required): User's country
- `os` (required): Operating system (android/ios)
- `limit` (optional): Maximum number of campaigns to return

**Example:**
```bash
curl "http://localhost:7000/api/v1/campaigns?app=myapp&country=india&os=android&limit=3"
```

#### 3. Get Campaigns (POST)
```http
POST /api/v1/campaigns
Content-Type: application/json
```

**Request Body:**
```json
{
  "app": "myapp",
  "country": "india",
  "os": "android",
  "limit": 3
}
```

**Example:**
```bash
curl -X POST http://localhost:7000/api/v1/campaigns \
  -H "Content-Type: application/json" \
  -d '{
    "app": "myapp",
    "country": "india",
    "os": "android",
    "limit": 3
  }'
```

#### 4. Refresh Targeting Data
```http
POST /api/v1/refresh-targeting
```

**Example:**
```bash
curl -X POST http://localhost:7000/api/v1/refresh-targeting
```

### Response Format

**Success Response:**
```json
{
  "campaigns": [
    {
      "id": "1",
      "name": "Summer Sale Campaign",
      "image_url": "https://example.com/summer-sale.jpg",
      "cta": "Shop Now",
      "status": "ACTIVE",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1
}
```

**Error Response:**
```json
{
  "error": "Failed to get campaigns",
  "message": "missing app param"
}
```

## 🧪 Testing

### Automated Testing

Run the comprehensive test script:

```bash
chmod +x test_api.sh
./test_api.sh
```

### Manual Testing

#### Test Health Check
```bash
curl http://localhost:7000/health
```

#### Test Campaign Retrieval
```bash
# GET request
curl "http://localhost:7000/api/v1/campaigns?app=testapp&country=india&os=android"

# POST request
curl -X POST http://localhost:7000/api/v1/campaigns \
  -H "Content-Type: application/json" \
  -d '{"app": "testapp", "country": "india", "os": "android", "limit": 2}'
```

## 📊 Database Schema

### Campaigns Table
```sql
CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url TEXT,
    cta VARCHAR(255),
    status VARCHAR(50) DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Targeting Rules Table
```sql
CREATE TABLE targeting_rules (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id),
    rule_type VARCHAR(50) NOT NULL, -- 'INCLUDE' or 'EXCLUDE'
    dimension VARCHAR(50) NOT NULL, -- 'APP', 'COUNTRY', 'OS'
    values TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_NAME` | Database name | `targeting_engine` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `SERVER_PORT` | HTTP server port | `7000` |

## 🏗️ Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   ├── campaign.go          # Campaign models
│   │   └── targeting.go         # Targeting models
│   ├── repository/
│   │   └── campaign_repository.go # Data access layer
│   ├── service/
│   │   └── delivery_service.go   # Business logic
│   └── transport/
│       └── http/
│           ├── handlers.go       # HTTP handlers
│           └── routes.go         # Route definitions
├── pkg/
│   ├── cache/
│   │   └── redis.go             # Redis cache implementation
│   └── database/
│       └── postgres.go          # PostgreSQL connection
├── scripts/
│   └── setup_db.go              # Database initialization
├── test_api.sh                  # API testing script
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── README.md                    # This file
```

## 🚨 Troubleshooting

### Common Issues

1. **"missing app param" error**
   - Ensure you're providing all required parameters: `app`, `country`, `os`

2. **Database connection failed**
   - Check if PostgreSQL is running: `brew services list | grep postgresql`
   - Verify database credentials and connection details

3. **Redis connection failed**
   - Check if Redis is running: `brew services list | grep redis`
   - Test Redis connection: `redis-cli ping`

4. **Port already in use**
   - Change the `SERVER_PORT` environment variable
   - Kill existing process: `lsof -ti:7000 | xargs kill -9`

### Logs and Debugging

The application logs important events to stdout. To see detailed logs:

```bash
# Run with verbose logging
./targeting-engine 2>&1 | tee app.log
```

## 🤝 Development

### Adding New Features

1. **New Targeting Dimensions**: Add to `internal/models/targeting.go`
2. **New API Endpoints**: Add handlers in `internal/transport/http/handlers.go`
3. **Database Changes**: Update schema and add migration scripts

### Code Style

- Follow Go conventions and use `gofmt`
- Add tests for new functionality
- Update documentation for API changes

## 📝 Sample Data

The setup script creates sample campaigns and targeting rules:

**Campaigns:**
- Summer Sale Campaign
- Winter Collection
- Mobile App Promotion

**Targeting Rules:**
- App-based targeting (myapp, otherapp)
- Country-based targeting (india, usa)
- OS-based targeting (android, ios)

## 🔒 Security Considerations

- Input validation on all API endpoints
- SQL injection prevention using parameterized queries
- Rate limiting (can be added with middleware)
- CORS configuration for production deployment

---

**Happy Coding! 🎯** 