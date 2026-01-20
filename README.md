# RSS Aggregator

A Go-based REST API service that aggregates RSS feeds from multiple sources. Users can subscribe to feeds, and the service automatically fetches and stores content in the background, providing a personalized feed of posts.

## Features

- **User Management** - Create accounts with auto-generated API keys for authentication
- **Feed Management** - Subscribe to RSS feeds by providing a name and URL
- **Automatic Fetching** - Background service continuously fetches and stores feed content
- **Personalized Feed** - View posts from all feeds you follow in one place
- **Subscription Control** - Follow/unfollow feeds easily
- **Post Deduplication** - Automatic handling of duplicate posts

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21.3+ |
| Web Framework | Chi (github.com/go-chi/chi) |
| CORS | Chi CORS |
| Database | PostgreSQL |
| ORM/Query Generation | sqlc |
| Migrations | goose |

## Getting Started

### Prerequisites

- Go 1.21.3 or later
- PostgreSQL database
- sqlc (for code generation)
- goose (for database migrations)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd RSS-Aggregator
   ```

2. **Configure environment variables**
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/rss_aggregator?sslmode=disable
   ```

3. **Setup the database**
   ```bash
   # Create the database
   createdb rss_aggregator

   # Run migrations
   goose postgres "postgresql://user:password@localhost:5432/rss_aggregator?sslmode=disable" up
   ```

4. **Generate database code** (if modifying SQL queries)
   ```bash
   sqlc generate
   ```

5. **Build and run**
   ```bash
   go build -o RSS-Aggregator
   ./RSS-Aggregator
   ```

## API Reference

All endpoints are prefixed with `/v1`.

### Authentication

Most endpoints require authentication using an API key. Include the key in the request header:

```
Authorization: ApiKey <your-api-key>
```

### Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/healthz` | No | Health check |
| GET | `/err` | No | Error test endpoint |
| POST | `/users` | No | Create a new user |
| GET | `/users` | Yes | Get current user info |
| POST | `/feeds` | Yes | Create a new feed |
| GET | `/feeds` | No | Get all available feeds |
| POST | `/feed_follows` | Yes | Follow a feed |
| GET | `/feed_follows` | Yes | Get user's feed follows |
| DELETE | `/feed_follows/{id}` | Yes | Unfollow a feed |
| GET | `/posts` | Yes | Get posts from followed feeds |

### Usage Examples

#### Create a User

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "My Username"}'
```

Response:
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "name": "My Username",
  "api_key": "your-api-key-here"
}
```

**Important:** Save your `api_key` - you'll need it for authentication.

#### Create a Feed

```bash
curl -X POST http://localhost:8080/v1/feeds \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <YOUR_API_KEY>" \
  -d '{
    "name": "Example Feed",
    "url": "https://example.com/feed.xml"
  }'
```

#### Follow a Feed

```bash
curl -X POST http://localhost:8080/v1/feed_follows \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <YOUR_API_KEY>" \
  -d '{"feed_id": "<FEED_ID>"}'
```

#### Get Your Posts

```bash
curl http://localhost:8080/v1/posts \
  -H "Authorization: ApiKey <YOUR_API_KEY>"
```

## Project Structure

```
RSS-Aggregator/
├── main.go                           # Application entry point
├── models.go                         # Domain models and converters
├── json.go                           # JSON response helpers
├── rss.go                            # RSS feed parsing
├── scraper.go                        # Background RSS scraping service
├── middleware_auth.go                # API key authentication middleware
├── handler_*.go                      # HTTP handlers
│   ├── handler_user.go               # User operations
│   ├── handler_feed.go               # Feed operations
│   ├── handler_feed_follows.go       # Feed subscription operations
│   ├── handler_readiness.go          # Health check
│   └── handler_err.go                # Error testing
├── internal/
│   ├── auth/
│   │   └── auth.go                   # API key extraction
│   └── database/
│       ├── db.go                     # Database connection
│       ├── models.go                 # Database models
│       ├── users.sql.go              # User queries
│       ├── feeds.sql.go              # Feed queries
│       ├── posts.sql.go              # Post queries
│       └── feed_follows.sql.go       # Feed follow queries
├── sql/
│   ├── schema/                       # Database migrations
│   │   ├── 001_users.sql
│   │   ├── 002_users_apikey.sql
│   │   ├── 003_feeds.sql
│   │   ├── 004_feed_follows.sql
│   │   ├── 005_feeds_lastfetchedat.sql
│   │   └── 006_posts.sql
│   └── queries/                      # SQL queries for sqlc
│       ├── users.sql
│       ├── feeds.sql
│       ├── posts.sql
│       └── feed_follows.sql
├── sqlc.yaml                         # sqlc configuration
├── go.mod                            # Go module
├── go.sum                            # Checksums
└── .env                              # Environment variables
```

## Database Schema

### Users
- `id` (UUID, Primary Key)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `name` (TEXT)
- `api_key` (VARCHAR(64), UNIQUE) - SHA256 hashed

### Feeds
- `id` (UUID, Primary Key)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `name` (TEXT)
- `url` (TEXT, UNIQUE)
- `user_id` (UUID, Foreign Key → users.id)
- `last_fetched_at` (TIMESTAMP, nullable)

### FeedFollows
- `id` (UUID, Primary Key)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `user_id` (UUID, Foreign Key → users.id)
- `feed_id` (UUID, Foreign Key → feeds.id)
- Unique constraint on (user_id, feed_id)

### Posts
- `id` (UUID, Primary Key)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `title` (TEXT)
- `description` (TEXT, nullable)
- `published_at` (TIMESTAMP)
- `url` (TEXT, UNIQUE)
- `feed_id` (UUID, Foreign Key → feeds.id)

## Background Scraper

The RSS aggregator runs a background scraper that:

- Continuously fetches RSS feeds at configurable intervals (default: 1 minute)
- Uses concurrent workers for efficient fetching (default: 10 parallel workers)
- Automatically deduplicates posts
- Tracks the last fetch time for each feed

## Configuration

The application can be configured through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP server port | 8080 |
| DB_URL | PostgreSQL connection string | - |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on how to contribute to this project.
