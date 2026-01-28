package auth

import (
	"template/modules/core/pkg/crud"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RoleProvider exposes CRUD operations for roles.
type RoleProvider struct {
	db *gorm.DB
}

func NewRoleProvider(db *gorm.DB) *RoleProvider {
	return &RoleProvider{db: db}
}

// GetModelName implements crud.CRUDProvider.
func (p *RoleProvider) GetModelName() string {
	return "roles"
}

// GetSchema implements crud.CRUDProvider.
func (p *RoleProvider) GetSchema() crud.Schema {
	return crud.Schema{
		Name:        "roles",
		DisplayName: "Roles",
		Fields: []crud.Field{
			{Name: "id", Type: "number", Label: "ID", Readonly: true, Editable: true, Width: "80px"},
			{Name: "name", Type: "string", Label: "Name", Required: true, Editable: true, Width: "250px"},
			{Name: "created_at", Type: "date", Label: "Created", Readonly: true, Editable: true, Width: "250px"},
			{Name: "updated_at", Type: "date", Label: "Updated", Readonly: true, Editable: true, Width: "250px"},
		},
		Searchable: []string{"name"},
	}
}

// CRUD Operations using default implementations.

func (p *RoleProvider) List(filters map[string]string, page, limit int) (crud.ListResponse, error) {
	return crud.DefaultList(p.db, &Role{}, filters, page, limit)
}

func (p *RoleProvider) Get(id string) (any, error) {
	return crud.DefaultGet(p.db, &Role{}, id)
}

func (p *RoleProvider) Create(data map[string]any) (any, error) {
	return crud.DefaultCreate(p.db, &Role{}, data)
}

func (p *RoleProvider) Update(id string, data map[string]any) (any, error) {
	return crud.DefaultUpdate(p.db, &Role{}, id, data)
}

func (p *RoleProvider) Delete(id string) error {
	return crud.DefaultDelete(p.db, &Role{}, id)
}

// HTTP Handlers.

func (p *RoleProvider) ListHandler() fiber.Handler {
	return crud.DefaultListHandler(p)
}

func (p *RoleProvider) SchemaHandler() fiber.Handler {
	return crud.DefaultSchemaHandler(p)
}

func (p *RoleProvider) GetHandler() fiber.Handler {
	return crud.DefaultGetHandler(p)
}

func (p *RoleProvider) CreateHandler() fiber.Handler {
	return crud.DefaultCreateHandler(p)
}

func (p *RoleProvider) UpdateHandler() fiber.Handler {
	return crud.DefaultUpdateHandler(p)
}

func (p *RoleProvider) DeleteHandler() fiber.Handler {
	return crud.DefaultDeleteHandler(p)
}
