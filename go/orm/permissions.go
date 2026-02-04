package orm

import (
	"slices"

	"gorm.io/gorm"
)

// check if user has a specific permission in the server
func (m *Membership) CanUser(perm string, db *gorm.DB) bool {
	// banned users have no permissions
	if m.Banned {
		return false
	}

	// ensure server is loaded
	if m.Server.ID == "" {
		if err := db.First(&m.Server, "id = ?", m.ServerID).Error; err != nil {
			return false
		}
	}

	// server owner has all permissions
	if m.Server.OwnerID == m.UserID {
		return true
	}

	// get all user permissions for this server
	userPerms := m.GetAllPermissions(db)

	// check if user has the required permission
	return slices.Contains(userPerms, perm)
}

// check if user can assign a specific role
func (m *Membership) CanAssignRole(role *Role, db *gorm.DB) bool {
	// must have manage_roles permission
	if !m.CanUser("manage_roles", db) {
		return false
	}

	// ensure server is loaded
	if m.Server.ID == "" {
		if err := db.First(&m.Server, "id = ?", m.ServerID).Error; err != nil {
			return false
		}
	}

	// server owner can assign any role
	if m.Server.OwnerID == m.UserID {
		return true
	}

	// can only assign roles with permissions the user has
	userPerms := m.GetAllPermissions(db)

	// check if all role permissions are in user permissions
	for _, rolePerm := range role.Permissions {
		found := slices.Contains(userPerms, rolePerm)
		if !found {
			return false
		}
	}
	return true
}

// get all permissions for a user in this server
func (m *Membership) GetAllPermissions(db *gorm.DB) []string {
	// banned users have no permissions
	if m.Banned {
		return []string{}
	}

	// ensure server is loaded
	if m.Server.ID == "" {
		if err := db.First(&m.Server, "id = ?", m.ServerID).Error; err != nil {
			return []string{}
		}
	}

	// server owner has all permissions
	if m.Server.OwnerID == m.UserID {
		return []string{
			"administrator",
			"manage_server",
			"manage_roles",
			"manage_channels",
			"kick_users",
			"ban_users",
			"create_invites",
			"manage_invites",
			"view_channels",
			"send_messages",
			"manage_messages",
			"embed_links",
			"attach_files",
			"mention_everyone",
		}
	}

	// get all user roles for this server
	var userRoles []UserRole
	db.Where("user_id = ? AND server_id = ?", m.UserID, m.ServerID).Preload("Role").Find(&userRoles)

	// union all permissions from roles
	permSet := make(map[string]bool)
	for _, ur := range userRoles {
		for _, perm := range ur.Role.Permissions {
			permSet[perm] = true
		}
	}

	// convert set to slice
	var userPerms []string
	for perm := range permSet {
		userPerms = append(userPerms, perm)
	}

	return userPerms
}
