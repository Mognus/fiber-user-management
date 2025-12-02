# Auth Module

JWT-based authentication with user registration, login, and role-based access control.

## Features

- ✅ User registration with email/password
- ✅ Login with JWT token generation
- ✅ Logout (cookie clearing)
- ✅ Get current user endpoint
- ✅ Password hashing with bcrypt
- ✅ Role-based access control (admin, user, guest)
- ✅ HTTP-only cookies for security
- ✅ Structured error handling

## Installation

```bash
go get github.com/yourcompany/auth-module
```

## Usage

```go
import (
    "github.com/yourcompany/auth-module/pkg/auth"
    "yourproject/pkg/module"
)

func main() {
    db := connectDatabase()
    jwtSecret := os.Getenv("JWT_SECRET")
    
    // Create module registry
    registry := module.NewRegistry()
    
    // Register auth module
    registry.Register(auth.New(db, jwtSecret))
    
    // Run migrations
    registry.MigrateAll(db)
    
    // Register routes
    api := app.Group("/api")
    registry.RegisterAll(api)
}
```

## Environment Variables

```bash
JWT_SECRET=your-secret-key-change-in-production
```

## API Endpoints

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/me` - Get current authenticated user

## Example Requests

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Get Current User
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Logout
```bash
curl -X POST http://localhost:8080/api/auth/logout
```

## User Model

```go
type User struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      UserRole  `json:"role"` // admin, user, guest
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## User Roles

- `admin` - Full system access
- `user` - Standard user access
- `guest` - Limited access

## Security

- Passwords are hashed with bcrypt before storage
- JWT tokens expire after 24 hours
- Tokens stored in HTTP-only cookies
- Set `Secure: true` in production with HTTPS

## Dependencies

- `github.com/yourcompany/backend-core` - Error handling
- `github.com/golang-jwt/jwt/v5` - JWT handling
- `golang.org/x/crypto` - Password hashing
- `gorm.io/gorm` - ORM
