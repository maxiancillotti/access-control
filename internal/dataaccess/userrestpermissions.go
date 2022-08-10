package dataaccess

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/pkg/errors"
)

type usersRESTPermissionsDAL struct {
	sqlDB *sql.DB
}

func NewUserRESTPermissionsDAL(sqlDB *sql.DB) dal.UsersRESTPermissionsDAL {
	return &usersRESTPermissionsDAL{
		sqlDB: sqlDB,
	}
}

// Inserts a User and returns the created ID
func (d *usersRESTPermissionsDAL) Insert(permission *domain.UserRESTPermission) error {

	_, err := d.sqlDB.Exec("UsersRESTPermissions_Insert",
		sql.Named("userID", permission.UserID),
		sql.Named("resourceID", permission.ResourceID),
		sql.Named("methodID", permission.MethodID),
	)
	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *usersRESTPermissionsDAL) Delete(permission *domain.UserRESTPermission) error {

	_, err := d.sqlDB.Exec("UsersRESTPermissions_Delete",
		sql.Named("userID", permission.UserID),
		sql.Named("resourceID", permission.ResourceID),
		sql.Named("methodID", permission.MethodID),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *usersRESTPermissionsDAL) SelectAllByUserID(userID uint) ([]dto.PermissionsIDs, error) {

	// The query result should be ordered by: UserID, ResourceID, MethodID
	// This method depends on this to traverse the data structure correctly
	// and effeciently.
	rows, err := d.sqlDB.Query("UsersRESTPermissions_SelectAllByUserID", sql.Named("userID", userID))

	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var permission dto.PermissionsIDs
	var permissions = make([]dto.PermissionsIDs, 0)

	var rowResourceID, rowMethodID uint
	var currentResource uint
	var rowsHasBeenRead bool

	for rows.Next() {
		err = rows.Scan(&rowResourceID, &rowMethodID)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}

		// new resource in the result
		if currentResource != rowResourceID {

			// flush previous permission into slice when index doesn't equal initial state
			if currentResource > 0 {
				permissions = append(permissions, permission)
			} else {
				// when index equals initial state
				rowsHasBeenRead = true
			}

			permission = dto.PermissionsIDs{
				ResourceID: rowResourceID,
				MethodsIDs: make([]uint, 0),
			}
		}

		permission.MethodsIDs = append(permission.MethodsIDs, rowMethodID)
	}
	if rowsHasBeenRead {
		// last pending append
		permissions = append(permissions, permission)
	} else {
		return nil, dal.ErrorDataAccessEmptyResult
		/*
			if len(permissions) == 0 {
				return nil, dal.ErrorDataAccessEmptyResult
			}
		*/
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	return permissions, nil
}

func (d *usersRESTPermissionsDAL) SelectAllWithDescriptionsByUserID(userID uint) ([]dto.PermissionsWithDescriptions, error) {

	// Returns:
	/*
		URP.ResourceID,
		R.[Path] 'ResourcePath',
		URP.HttpMethodID,
		M.[Name] 'HttpMethodName'
	*/
	// The query result should be ordered by: ResourceID, MethodID
	// This method depends on this to traverse the data structure correctly
	// and effectively.
	rows, err := d.sqlDB.Query("UsersRESTPermissions_SelectAllWithDescriptionsByUserID", sql.Named("userID", userID))
	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var permissions []dto.PermissionsWithDescriptions
	var permission dto.PermissionsWithDescriptions

	var currentResource uint
	var rowsHasBeenRead bool

	for rows.Next() {

		var rowResourceID, rowMethodID uint
		var rowResourcePath, rowMethodName string

		err = rows.Scan(&rowResourceID, &rowResourcePath, &rowMethodID, &rowMethodName)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}

		// new resource in the result
		if currentResource != rowResourceID {

			// flush previous permission into slice when index doesn't equal initial state
			if currentResource > 0 {
				permissions = append(permissions, permission)
			} else {
				// when index equals initial state
				rowsHasBeenRead = true
			}

			// creating new resource permission
			permission = dto.PermissionsWithDescriptions{
				Resource: dto.Resource{
					ID:   &rowResourceID,
					Path: &rowResourcePath,
				},
				Methods: make([]dto.HttpMethod, 0),
			}

			// updating index
			currentResource = rowResourceID
		}

		permission.Methods = append(permission.Methods, dto.HttpMethod{ID: &rowMethodID, Name: &rowMethodName})
	}
	if rowsHasBeenRead {
		// last pending append
		permissions = append(permissions, permission)
	} else {
		return nil, dal.ErrorDataAccessEmptyResult
		/*
			if len(permissions) == 0 {
				return nil, dal.ErrorDataAccessEmptyResult
			}
		*/
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	return permissions, nil
}

func (d *usersRESTPermissionsDAL) SelectAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error) {

	permissions := make(domain.RESTPermissionsPathsMethods)

	// Returns:
	/*
		R.[Path] 'Resource',
		M.[Name] 'Method'
	*/
	// The query result should be ordered by: Resource, Method.
	// This method depends on this to traverse the data structure correctly
	// and effectively.
	rows, err := d.sqlDB.Query("UsersRESTPermissions_SelectAllPathsMethodsByUserID", sql.Named("userID", userID))
	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var rowResource, rowMethod string
	for rows.Next() {
		err = rows.Scan(&rowResource, &rowMethod)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		permissions[rowResource] = append(permissions[rowResource], rowMethod)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(permissions) == 0 {
		return nil, dal.ErrorDataAccessEmptyResult
	}

	return permissions, nil
}

func (d *usersRESTPermissionsDAL) Exists(permission *domain.UserRESTPermission) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("UsersRESTPermissions_Exists",
		sql.Named("userID", permission.UserID),
		sql.Named("resourceID", permission.ResourceID),
		sql.Named("methodID", permission.MethodID),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}

func (d *usersRESTPermissionsDAL) RelationshipsExists(permission *domain.UserRESTPermission) (int, error) {

	var exists int

	// According to the DAL definition, this SP is expected
	// to return the following values:
	// 1 == All relationships exist.
	// -1 == UserID does not exist.
	// -2 == ResourceID does not exist.
	// -3 == HttpMethodID does not exist.
	_, err := d.sqlDB.Exec("UsersRESTPermissions_RelationshipsExists",
		sql.Named("userID", permission.UserID),
		sql.Named("resourceID", permission.ResourceID),
		sql.Named("methodID", permission.MethodID),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return 0, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}
