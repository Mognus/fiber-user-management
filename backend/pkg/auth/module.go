package auth

import (
	"errors"
	"time"

	apperrors "template/modules/core/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Module implements the pluggable auth module interface
type Module struct {
	db            *gorm.DB
	jwtSecret     string
	jwtMiddleware fiber.Handler
	userProvider  *UserProvider
	roleProvider  *RoleProvider
}

// New creates a new Auth module instance
func New(db *gorm.DB, jwtSecret string) *Module {
	return &Module{
		db:            db,
		jwtSecret:     jwtSecret,
		jwtMiddleware: nil, // Will be set by SetJWTMiddleware
		userProvider:  NewUserProvider(db),
		roleProvider:  NewRoleProvider(db),
	}
}

// SetJWTMiddleware sets the JWT middleware for protected routes
func (m *Module) SetJWTMiddleware(middleware fiber.Handler) {
	m.jwtMiddleware = middleware
}

func (m *Module) UserProvider() *UserProvider {
	return m.userProvider
}

func (m *Module) RoleProvider() *RoleProvider {
	return m.roleProvider
}

// Name returns the module name
func (m *Module) Name() string {
	return "auth"
}

// RegisterRoutes registers authentication routes
func (m *Module) RegisterRoutes(router fiber.Router) {
	// Public auth routes
	auth := router.Group("/auth")
	auth.Post("/register", m.Register)
	auth.Post("/login", m.Login)
	auth.Post("/logout", m.Logout)

	// Protected auth routes
	if m.jwtMiddleware != nil {
		auth.Get("/me", m.jwtMiddleware, m.Me)
	} else {
		auth.Get("/me", m.Me) // Fallback without middleware
	}

	// CRUD routes for users
	users := router.Group("/users")
	users.Get("/", m.userProvider.ListHandler())
	users.Get("/schema", m.userProvider.SchemaHandler())
	users.Get("/:id", m.userProvider.GetHandler())
	users.Post("/", m.userProvider.CreateHandler())
	users.Put("/:id", m.userProvider.UpdateHandler())
	users.Delete("/:id", m.userProvider.DeleteHandler())

	roles := router.Group("/roles")
	roles.Get("/", m.roleProvider.ListHandler())
	roles.Get("/schema", m.roleProvider.SchemaHandler())
	roles.Get("/:id", m.roleProvider.GetHandler())
	roles.Post("/", m.roleProvider.CreateHandler())
	roles.Put("/:id", m.roleProvider.UpdateHandler())
	roles.Delete("/:id", m.roleProvider.DeleteHandler())
}

// Migrate runs database migrations
func (m *Module) Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Role{}); err != nil {
		return err
	}

	for _, role := range defaultRoles() {
		if err := db.FirstOrCreate(&role, Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents auth response with token
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// Register creates a new user account
// POST /api/auth/register
func (m *Module) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)

	if err := c.BodyParser(req); err != nil {
		appErr := apperrors.BadRequest("Invalid request body")
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Validation
	if req.Email == "" {
		appErr := apperrors.ValidationError(map[string]string{
			"email": "Email is required",
		})
		return c.Status(appErr.Code).JSON(appErr)
	}
	if req.Password == "" || len(req.Password) < 8 {
		appErr := apperrors.ValidationError(map[string]string{
			"password": "Password must be at least 8 characters",
		})
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Check if user already exists
	var existingUser User
	if err := m.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		appErr := apperrors.Conflict("User with this email already exists")
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Create new user
	user := User{
		Email:     req.Email,
		Password:  req.Password, // Will be hashed by BeforeCreate hook
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      RoleUser,
		Active:    true,
	}

	if err := m.db.Create(&user).Error; err != nil {
		appErr := apperrors.Internal(err)
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Generate JWT token
	token, err := m.generateToken(&user)
	if err != nil {
		appErr := apperrors.InternalWithMessage("Failed to generate token", err)
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Set HTTP-only cookie
	m.setAuthCookie(c, token)

	return c.Status(201).JSON(AuthResponse{
		Token: token,
		User:  &user,
	})
}

// Login authenticates a user
// POST /api/auth/login
func (m *Module) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)

	if err := c.BodyParser(req); err != nil {
		appErr := apperrors.BadRequest("Invalid request body")
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Validation
	if req.Email == "" || req.Password == "" {
		appErr := apperrors.ValidationError(map[string]string{
			"email":    "Email is required",
			"password": "Password is required",
		})
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Find user by email
	var user User
	if err := m.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appErr := apperrors.Unauthorized("Invalid email or password")
			return c.Status(appErr.Code).JSON(appErr)
		}
		appErr := apperrors.Internal(err)
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Check if user is active
	if !user.Active {
		appErr := apperrors.Forbidden("Account is deactivated")
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		appErr := apperrors.Unauthorized("Invalid email or password")
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Generate JWT token
	token, err := m.generateToken(&user)
	if err != nil {
		appErr := apperrors.InternalWithMessage("Failed to generate token", err)
		return c.Status(appErr.Code).JSON(appErr)
	}

	// Set HTTP-only cookie
	m.setAuthCookie(c, token)

	return c.JSON(AuthResponse{
		Token: token,
		User:  &user,
	})
}

// Logout logs out the user
// POST /api/auth/logout
func (m *Module) Logout(c *fiber.Ctx) error {
	// Clear the auth cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})

	return c.Status(200).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user
// GET /api/auth/me
func (m *Module) Me(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Fetch user from database
	var currentUser User
	if err := m.db.First(&currentUser, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appErr := apperrors.NotFound("User")
			return c.Status(appErr.Code).JSON(appErr)
		}
		appErr := apperrors.Internal(err)
		return c.Status(appErr.Code).JSON(appErr)
	}

	return c.JSON(currentUser)
}

// generateToken creates a JWT token for the user
func (m *Module) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.jwtSecret))
}

// setAuthCookie sets the JWT token as an HTTP-only cookie
func (m *Module) setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})
}
