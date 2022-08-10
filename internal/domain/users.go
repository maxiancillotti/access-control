package domain

type User struct {
	ID           uint
	Username     string
	PasswordHash string // Salted and hashed password
	PasswordSalt string
	Enabled      bool
}

// PermissionCategory map[PermissionCategory]
// Resource    map[string]
// Permissions []string
type UserPermissions map[PermissionCategory]map[string][]string

type PermissionCategory string

var (
	RESTpermissionCategory PermissionCategory = "REST"
	GRPCpermissionCategory PermissionCategory = "gRPC"
)
