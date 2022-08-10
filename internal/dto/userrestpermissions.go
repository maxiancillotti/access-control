package dto

import "github.com/pkg/errors"

type UserRESTPermission struct {
	UserID     *uint       `json:"user_id"`
	Permission *Permission `json:"permission"`
}

func (up *UserRESTPermission) ValidateEmpty() error {
	if up.UserID == nil || up.Permission == nil {
		return errors.New("fields cannot be empty")
	}
	return up.Permission.ValidateEmpty()
}

type Permission struct {
	ResourceID *uint `json:"resource_id"`
	MethodID   *uint `json:"method_id"`
}

func (p *Permission) ValidateEmpty() error {
	if p.ResourceID == nil || p.MethodID == nil {
		return errors.New("fields cannot be empty")
	}
	return nil
}

/********************************************************/

// Intersection transfer
type UserRESTPermissionsCollection struct {
	UserID      uint             `json:"user_id"`
	Permissions []PermissionsIDs `json:"permissions_ids"`
}

type PermissionsIDs struct {
	ResourceID uint   `json:"resource_id"`
	MethodsIDs []uint `json:"method_ids"`
}

// JSON Encoding Test
/*
func getURP() *UserRESTPermissionsCollection {

	collection := &UserRESTPermissionsCollection{
		UserID: 1,
		Permissions: []ResourceIDsMethodsIDs{
			ResourceIDsMethodsIDs{
				ResourceID: 15,
				MethodsIDs: []uint{4, 5, 6},
			},
			ResourceIDsMethodsIDs{
				ResourceID: 27,
				MethodsIDs: []uint{8, 9, 10},
			},
		},
	}

	return collection
}

func jsonEnc() ([]byte, error) {

	urp := f()
	return json.Marshal(urp)

}
*/
/********************************************************/

// Intersection transfer with permissions's attributes names
type UserRESTPermissionsDescriptionsCollection struct {
	UserID      uint
	Permissions []PermissionsWithDescriptions
}

type PermissionsWithDescriptions struct {
	Resource `json:"resource"`
	Methods  []HttpMethod `json:"methods"`
}
