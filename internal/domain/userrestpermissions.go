package domain

type UserRESTPermission struct {
	UserID     uint
	ResourceID uint
	MethodID   uint
}

// Map Key: Resource    string.
// Map Value: Methods []string.
type RESTPermissionsPathsMethods map[string][]string
