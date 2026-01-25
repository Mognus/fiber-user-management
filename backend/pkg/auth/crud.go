package auth

import (
	"template/modules/core/pkg/crud"
)

// GetModelName implements crud.CRUDProvider
func (m *Module) GetModelName() string {
	return "users"
}

// GetSchema implements crud.CRUDProvider
func (m *Module) GetSchema() crud.Schema {
	return crud.Schema{
		Name:        "users",
		DisplayName: "Users",
		Fields: []crud.Field{
			{Name: "id", Type: "number", Label: "ID", Readonly: true, Editable: true},
			{Name: "email", Type: "string", Label: "Email", Required: true, Editable: true},
			{Name: "password", Type: "string", Label: "Password", Required: true, Editable: false},
			{Name: "first_name", Type: "string", Label: "First Name", Editable: true},
			{Name: "last_name", Type: "string", Label: "Last Name", Editable: true},
			{
				Name:       "role",
				Type:       "enum",
				Label:      "Role",
				EnumValues: []string{"admin", "user", "guest"},
				Editable:   true,
			},
			{Name: "active", Type: "boolean", Label: "Active", Editable: true},
			{Name: "created_at", Type: "date", Label: "Created", Readonly: true, Editable: true},
			{Name: "updated_at", Type: "date", Label: "Updated", Readonly: true, Editable: true},
		},
		Filterable: []string{"role", "active"},
		Searchable: []string{"email", "first_name", "last_name"},
	}
}

// CRUD Operations using default implementations

func (m *Module) List(filters map[string]string, page, limit int) (crud.ListResponse, error) {
	return crud.DefaultList(m.db, &User{}, filters, page, limit)
}

func (m *Module) Get(id string) (any, error) {
	return crud.DefaultGet(m.db, &User{}, id)
}

func (m *Module) Create(data map[string]any) (any, error) {
	return crud.DefaultCreate(m.db, &User{}, data)
}

func (m *Module) Update(id string, data map[string]any) (any, error) {
	return crud.DefaultUpdate(m.db, &User{}, id, data)
}

func (m *Module) Delete(id string) error {
	return crud.DefaultDelete(m.db, &User{}, id)
}
