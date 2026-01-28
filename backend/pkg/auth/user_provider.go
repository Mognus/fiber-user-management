package auth

import (
	"template/modules/core/pkg/crud"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserProvider exposes CRUD operations for users.
type UserProvider struct {
	db *gorm.DB
}

func NewUserProvider(db *gorm.DB) *UserProvider {
	return &UserProvider{db: db}
}

// GetModelName implements crud.CRUDProvider.
func (p *UserProvider) GetModelName() string {
	return "users"
}

// GetSchema implements crud.CRUDProvider.
func (p *UserProvider) GetSchema() crud.Schema {
	// Load roles for relation options
	var roles []Role
	p.db.Find(&roles)
	roleOptions := make([]crud.SelectOption, len(roles))
	for i, r := range roles {
		roleOptions[i] = crud.SelectOption{Value: r.ID, Label: r.Name}
	}

	return crud.Schema{
		Name:        "users",
		DisplayName: "Users",
		Fields: []crud.Field{
			{Name: "id", Type: "number", Label: "ID", Readonly: true, Editable: true, Width: "80px"},
			{Name: "email", Type: "string", Label: "Email", Required: true, Editable: true, Width: "200px"},
			{Name: "password", Type: "string", Label: "Password", Required: true, Editable: false, Width: "180px"},
			{Name: "first_name", Type: "string", Label: "First Name", Editable: true, Width: "140px"},
			{Name: "last_name", Type: "string", Label: "Last Name", Editable: true, Width: "140px"},
			{Name: "role_id", Type: "relation", Label: "Role", Required: true, Editable: true, Hidden: true, Options: roleOptions},
			{Name: "role", Type: "object", Label: "Role", Readonly: true, Editable: false, Width: "120px"},
			{Name: "active", Type: "boolean", Label: "Active", Editable: true, Width: "100px"},
			{Name: "created_at", Type: "date", Label: "Created", Readonly: true, Editable: true, Width: "160px"},
			{Name: "updated_at", Type: "date", Label: "Updated", Readonly: true, Editable: true, Width: "160px"},
		},
		Filterable: []string{"role_id", "active"},
		Searchable: []string{"email", "first_name", "last_name"},
	}
}

// CRUD Operations using default implementations.

func (p *UserProvider) List(filters map[string]string, page, limit int) (crud.ListResponse, error) {
	return crud.DefaultList(p.db, &User{}, filters, page, limit, "Role")
}

func (p *UserProvider) Get(id string) (any, error) {
	return crud.DefaultGet(p.db, &User{}, id, "Role")
}

func (p *UserProvider) Create(data map[string]any) (any, error) {
	return crud.DefaultCreate(p.db, &User{}, data)
}

func (p *UserProvider) Update(id string, data map[string]any) (any, error) {
	return crud.DefaultUpdate(p.db, &User{}, id, data)
}

func (p *UserProvider) Delete(id string) error {
	return crud.DefaultDelete(p.db, &User{}, id)
}

// HTTP Handlers.

func (p *UserProvider) ListHandler() fiber.Handler {
	return crud.DefaultListHandler(p)
}

func (p *UserProvider) SchemaHandler() fiber.Handler {
	return crud.DefaultSchemaHandler(p)
}

func (p *UserProvider) GetHandler() fiber.Handler {
	return crud.DefaultGetHandler(p)
}

func (p *UserProvider) CreateHandler() fiber.Handler {
	return crud.DefaultCreateHandler(p)
}

func (p *UserProvider) UpdateHandler() fiber.Handler {
	return crud.DefaultUpdateHandler(p)
}

func (p *UserProvider) DeleteHandler() fiber.Handler {
	return crud.DefaultDeleteHandler(p)
}
