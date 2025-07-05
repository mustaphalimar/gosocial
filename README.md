# Go Social API

A modern social media API built with Go, featuring user authentication, posts, comments, following system, and real-time feeds.

## 🚀 Features

- **User Management**: Registration, authentication, and profile management
- **Post System**: Create, read, update, and delete posts with tags
- **Social Features**: Follow/unfollow users and personalized feeds
- **Comment System**: Interactive commenting on posts
- **Role-Based Access Control**: Admin, moderator, and user roles
- **Email Integration**: User activation via SendGrid
- **Real-time Feed**: Get posts from followed users with search and filtering
- **API Documentation**: Auto-generated Swagger documentation

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Chi Router
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Email Service**: SendGrid
- **Documentation**: Swagger/OpenAPI
- **Development**: Air (hot reloading)
- **Containerization**: Docker & Docker Compose
- **Database Migrations**: golang-migrate

## 📁 Project Structure

```
go-social/
├── cmd/
│   ├── api/           # API server and handlers
│   └── migrate/       # Database migrations and seeding
├── internal/
│   ├── auth/          # Authentication logic
│   ├── db/            # Database connection
│   ├── env/           # Environment configuration
│   ├── mailer/        # Email service
│   └── store/         # Data layer (models and repositories)
├── docs/              # Auto-generated Swagger documentation
├── scripts/           # Utility scripts
├── web/               # React.JS Front-end app.
└── bin/               # Compiled binaries
```

## 🚦 Getting Started

### Prerequisites

- Go 1.24 or later
- PostgreSQL 16.3 or later
- Docker & Docker Compose (optional)
- SendGrid API key (for email functionality)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/mustaphalimar/go-social.git
cd go-social
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up environment variables**

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server
ADDR=:8080
EXTERNAL_URL=localhost:8080
CLIENT_URL=http://localhost:5173
ENV=development

# Database
DATABASE_URL=postgresql://admin:admin@localhost:5432/social?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

# Authentication
JWT_SECRET=your-jwt-secret-here
BASIC_AUTH_USERNAME=admin
BASIC_AUTH_PASSWORD=admin

# Email
SENDGRID_API_KEY=your-sendgrid-api-key
FROM_EMAIL=noreply@yourdomain.com
```

4. **Start the database**

```bash
docker-compose up -d
```

5. **Run database migrations**

```bash
make migrate-up
```

6. **Seed the database (optional)**

```bash
make seed
```

7. **Start the development server**

```bash
# With hot reloading
air

# Or build and run directly
go run cmd/api/*.go
```

The API will be available at `http://localhost:8080`

## 📖 API Documentation

### Swagger UI

Access the interactive API documentation at: `http://localhost:8080/v1/swagger/index.html`

### Health Check

```bash
curl -X GET http://localhost:8080/v1/health \
  -H "Authorization: Basic YWRtaW46YWRtaW4="
```

### Main Endpoints

#### Authentication

- `POST /v1/auth/register` - Register a new user
- `POST /v1/auth/token` - Login and get JWT token

#### Users

- `GET /v1/users/{userId}` - Get user profile
- `PUT /v1/users/activate/{token}` - Activate user account
- `PUT /v1/users/{userId}/follow` - Follow a user
- `PUT /v1/users/{userId}/unfollow` - Unfollow a user
- `GET /v1/users/feed` - Get personalized feed

#### Posts

- `POST /v1/posts` - Create a new post
- `GET /v1/posts/{postId}` - Get a specific post
- `PATCH /v1/posts/{postId}` - Update a post
- `DELETE /v1/posts/{postId}` - Delete a post

## 🔧 Development

### Available Make Commands

```bash
# Database migrations
make migrate-up          # Apply all migrations
make migrate-down        # Rollback migrations
make migration name      # Create new migration

# Development
make seed               # Seed database with sample data
make gen-docs          # Generate Swagger documentation
```

### Hot Reloading

The project uses Air for hot reloading during development. Configuration is in `.air.toml`.

```bash
air
```

### Database Migrations

Create a new migration:

```bash
make migration create_new_table
```

This creates two files:

- `000XXX_create_new_table.up.sql` - Forward migration
- `000XXX_create_new_table.down.sql` - Rollback migration

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/store
```

## 📊 Database Schema

### Core Tables

- **users**: User accounts with roles and activation status
- **posts**: User posts with tags and versioning
- **comments**: Comments on posts
- **followers**: User following relationships
- **roles**: User roles (user, moderator, admin)
- **user_invitations**: Email activation tokens

### Key Features

- **Optimistic Locking**: Posts use versioning to prevent concurrent updates
- **Soft Relationships**: Foreign key constraints maintain data integrity
- **Indexing**: Optimized queries for feeds and searches
- **CITEXT**: Case-insensitive email handling

## 🔐 Authentication & Authorization

### JWT Authentication

- Access tokens expire in 3 days
- Tokens include user ID and role information
- Bearer token format: `Authorization: Bearer <token>`

### Role-Based Access Control

- **User**: Basic post creation and social features
- **Moderator**: Can edit any post
- **Admin**: Full access including user management

### User Activation

- New users receive activation emails
- Accounts are inactive until email verification
- Activation tokens expire in 3 days

## 📧 Email Integration

Uses SendGrid for transactional emails:

- User registration confirmation
- Password reset (if implemented)
- Account notifications

## 🐳 Docker Deployment

### Development

```bash
docker-compose up -d
```

### Production

```bash
# Build the application
docker build -t go-social .

# Run with environment variables
docker run -e DATABASE_URL=... -e JWT_SECRET=... -p 8080:8080 go-social
```

## 🔍 API Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Create a post

```bash
curl -X POST http://localhost:8080/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "My First Post",
    "content": "This is my first post on Go Social!",
    "tags": ["golang", "social", "api"]
  }'
```

### Get user feed

```bash
curl -X GET "http://localhost:8080/v1/users/feed?limit=10&offset=0&sort=desc" \
  -H "Authorization: Bearer <your-jwt-token>"
```

## 📈 Performance Considerations

- **Connection Pooling**: Configurable database connection limits
- **Context Timeouts**: 5-second query timeout for database operations
- **Pagination**: Feed endpoints support limit/offset pagination
- **Indexing**: Database indexes on frequently queried columns
- **Caching**: Ready for Redis integration for session management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Chi Router for the excellent HTTP router
- golang-migrate for database migrations
- Swagger for API documentation
- SendGrid for email services
- The Go community for amazing libraries and tools

**Built with ❤️ using Go**
