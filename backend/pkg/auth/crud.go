package auth

import (
	"fmt"
	"strconv"

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
			{Name: "id", Type: "number", Label: "ID", Readonly: true},
			{Name: "email", Type: "string", Label: "Email", Required: true},
			{Name: "password", Type: "string", Label: "Password", Required: true, Editable: false},
			{Name: "first_name", Type: "string", Label: "First Name"},
			{Name: "last_name", Type: "string", Label: "Last Name"},
			{
				Name:       "role",
				Type:       "enum",
				Label:      "Role",
				EnumValues: []string{"admin", "user", "guest"},
			},
			{Name: "active", Type: "boolean", Label: "Active"},
			{Name: "created_at", Type: "date", Label: "Created", Readonly: true},
			{Name: "updated_at", Type: "date", Label: "Updated", Readonly: true},
		},
		Filterable: []string{"role", "active"},
		Searchable: []string{"email", "first_name", "last_name"},
	}
}

// List implements crud.CRUDProvider
func (m *Module) List(filters map[string]string, page, limit int) (crud.ListResponse, error) {
	var users []User
	query := m.db.Model(&User{})

	// Apply filters
	if role, ok := filters["role"]; ok {
		query = query.Where("role = ?", role)
	}
	if active, ok := filters["active"]; ok {
		isActive := active == "true"
		query = query.Where("active = ?", isActive)
	}

	// Apply search if provided
	if search, ok := filters["search"]; ok && search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return crud.ListResponse{}, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return crud.ListResponse{}, err
	}

	// Convert to []any
	items := make([]any, len(users))
	for i, user := range users {
		items[i] = user
	}

	return crud.ListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// Get implements crud.CRUDProvider
func (m *Module) Get(id string) (any, error) {
	var user User
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	if err := m.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Create implements crud.CRUDProvider
func (m *Module) Create(data map[string]any) (any, error) {
	user := User{
		Active: true, // Default to active
		Role:   RoleUser, // Default role
	}

	// Map data to user fields
	if email, ok := data["email"].(string); ok {
		user.Email = email
	}
	if password, ok := data["password"].(string); ok {
		user.Password = password // Will be hashed by BeforeCreate hook
	}
	if firstName, ok := data["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := data["last_name"].(string); ok {
		user.LastName = lastName
	}
	if role, ok := data["role"].(string); ok {
		user.Role = UserRole(role)
	}
	if active, ok := data["active"].(bool); ok {
		user.Active = active
	}

	if err := m.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Update implements crud.CRUDProvider
func (m *Module) Update(id string, data map[string]any) (any, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	var user User
	if err := m.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Update fields that are provided
	updates := make(map[string]interface{})

	if email, ok := data["email"].(string); ok {
		updates["email"] = email
	}
	if firstName, ok := data["first_name"].(string); ok {
		updates["first_name"] = firstName
	}
	if lastName, ok := data["last_name"].(string); ok {
		updates["last_name"] = lastName
	}
	if role, ok := data["role"].(string); ok {
		updates["role"] = role
	}
	if active, ok := data["active"].(bool); ok {
		updates["active"] = active
	}

	// Handle password separately if provided
	if password, ok := data["password"].(string); ok && password != "" {
		hashedPassword, err := user.HashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		updates["password"] = hashedPassword
	}

	if err := m.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload user to get updated data
	if err := m.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Delete implements crud.CRUDProvider
func (m *Module) Delete(id string) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	// Soft delete
	if err := m.db.Delete(&User{}, userID).Error; err != nil {
		return err
	}

	return nil
}
