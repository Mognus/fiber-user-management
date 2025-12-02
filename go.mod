module github.com/yourcompany/auth-module

go 1.23

require (
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/yourcompany/backend-core v0.0.0
	golang.org/x/crypto v0.31.0
	gorm.io/gorm v1.31.1
)

// For local development
replace github.com/yourcompany/backend-core => ../backend-core
