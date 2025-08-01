# RSS Aggregator

A modern RSS feed aggregator built with Go, featuring real-time feed scraping, user authentication, and a RESTful API. This application allows users to subscribe to RSS feeds, automatically fetches new posts, and provides a clean API for accessing aggregated content.

## 🚀 Features

### Core Functionality

- **User Authentication**: Secure user registration and login with API key-based authentication
- **RSS Feed Management**: Add, view, and manage RSS feeds
- **Feed Following**: Subscribe to feeds and get personalized post streams
- **Automatic Scraping**: Background service that periodically fetches new posts from RSS feeds
- **RESTful API**: Clean, well-structured API endpoints for all operations
- **Database Integration**: PostgreSQL database with proper schema management using SQLC

### Technical Features

- **Concurrent Scraping**: Multi-threaded RSS feed scraping with configurable concurrency
- **Error Handling**: Robust error handling and logging throughout the application
- **CORS Support**: Cross-origin resource sharing enabled for web applications
- **Health Checks**: Built-in health check endpoint for monitoring
- **Docker Support**: Containerized deployment with Docker Compose

## 🏗️ Architecture

### Project Structure

```
rssagg/
├── main.go                 # Application entry point and routing
├── models.go               # Data models and conversion functions
├── scrapper.go             # RSS feed scraping logic
├── handler_*.go           # HTTP request handlers
├── middleware_auth.go      # Authentication middleware
├── json.go                 # JSON response utilities
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go module dependencies
├── sqlc.yaml              # SQLC configuration
├── internal/
│   ├── auth/              # Authentication utilities
│   └── database/          # Database models and queries (generated)
└── sql/
    ├── schema/            # Database migration files
    └── queries/           # SQL queries for SQLC
```

### Database Schema

- **Users**: User accounts with email, password, and API key
- **Feeds**: RSS feed information with URL and ownership
- **Feed Follows**: User subscriptions to feeds
- **Posts**: Individual RSS posts with metadata

## 🛠️ Technology Stack

- **Backend**: Go 1.24.5
- **Database**: PostgreSQL 15
- **ORM**: SQLC (type-safe SQL)
- **Web Framework**: Chi Router
- **Authentication**: API Key-based
- **Containerization**: Docker & Docker Compose
- **Environment**: Environment variables with godotenv
- **Database Migrations**: Goose

## 📋 Prerequisites

- Go 1.24.5 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd rssagg
   ```

2. **Create environment file**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the application**

   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   # If using goose for migrations
   goose -dir sql/schema postgres "host=localhost port=5432 user=rssagg_user password=rssagg_password dbname=rssagg sslmode=disable" up
   ```

### Local Development

1. **Install dependencies**

   ```bash
   go mod download
   ```

2. **Start PostgreSQL**

   ```bash
   docker-compose up postgres -d
   ```

3. **Set environment variables**

   ```bash
   export PORT=8080
   export DB_URL="postgres://rssagg_user:rssagg_password@localhost:5432/rssagg?sslmode=disable"
   ```

4. **Run the application**
   ```bash
   go run .
   ```

## 📚 API Documentation

### Authentication

All protected endpoints require an API key in the Authorization header:

```
Authorization: ApiKey <your-api-key>
```

### Endpoints

#### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user

#### Users

- `GET /users/me` - Get current user information

#### Feeds

- `POST /feeds/` - Create a new RSS feed
- `GET /feeds/` - Get all feeds
- `GET /feeds/me` - Get feeds created by current user
- `GET /feeds/{feedId}` - Get specific feed by ID

#### Feed Follows

- `POST /feed-follows/` - Follow a feed
- `GET /feed-follows/` - Get user's followed feeds
- `DELETE /feed-follows/{feedFollowID}` - Unfollow a feed

#### Posts

- `GET /posts/` - Get posts for current user (from followed feeds)

#### System

- `GET /` - Welcome message
- `GET /healthz` - Health check endpoint

## 🔧 Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_URL`: PostgreSQL connection string
- `SCRAPER_CONCURRENCY`: Number of concurrent scrapers (default: 10)
- `SCRAPER_INTERVAL`: Time between scraping runs (default: 10 minutes)

### Database Configuration

The application uses PostgreSQL with the following default configuration:

- Database: `rssagg`
- User: `rssagg_user`
- Password: `rssagg_password`
- Port: `5432`

## 🔄 RSS Scraping

The application includes a background scraper that:

- Runs every 10 minutes (configurable)
- Processes up to 10 feeds concurrently (configurable)
- Fetches new posts from RSS feeds
- Stores posts in the database
- Updates feed last fetched timestamp

## 🚀 Potential Enhancements

### 1. JWT Authentication

Replace the current API key authentication with JWT tokens for better security and user experience.

**Benefits:**

- Stateless authentication
- Token expiration and refresh
- Better security practices
- Easier integration with frontend applications

### 2. Pagination

Add pagination support to endpoints that return lists of data.

**Benefits:**

- Better performance with large datasets
- Reduced memory usage
- Better user experience

### 3. CLI Interface

Add a command-line interface for interacting with the application.

**Features:**

- User management commands
- Feed management
- Post viewing
- Configuration management

### 4. Real-time Notifications

Implement WebSocket support for real-time post notifications.

**Benefits:**

- Instant notifications of new posts
- Better user engagement
- Real-time updates

### 5. Feed Categories and Tags

Add support for organizing feeds into categories and tagging posts.

**Benefits:**

- Better content organization
- Improved search and filtering
- Personalized content discovery

### 6. Search and Filtering

Implement search functionality for posts and feeds.

**Features:**

- Full-text search
- Filter by date, feed, category
- Advanced search operators

### 7. Rate Limiting

Add rate limiting to prevent abuse and ensure fair usage.

### 8. Caching Layer

Implement Redis caching for frequently accessed data.

**Benefits:**

- Improved performance
- Reduced database load
- Better scalability

### 9. Metrics and Monitoring

Add comprehensive metrics and monitoring.

**Features:**

- Prometheus metrics
- Health checks
- Performance monitoring
- Error tracking

### 10. API Versioning

Implement API versioning for better backward compatibility.

### 11. Background Job Queue

Implement a proper job queue for RSS scraping.

**Benefits:**

- Better resource management
- Retry mechanisms
- Job prioritization
- Scalability

### 12. User Preferences

Add user preference management.

**Features:**

- Notification preferences
- Default feed settings
- UI preferences

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
