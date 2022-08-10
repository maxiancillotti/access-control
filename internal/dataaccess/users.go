package dataaccess

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type userDAL struct {
	sqlDB *sql.DB
}

func NewUserDAL(sqlDB *sql.DB) dal.UserDAL {
	return &userDAL{
		sqlDB: sqlDB,
	}
}

func (d *userDAL) Insert(user *domain.User) (uint, error) {

	var id uint

	_, err := d.sqlDB.Exec("UsersInsert",
		sql.Named("username", user.Username),
		sql.Named("passwordHash", user.PasswordHash),
		sql.Named("passwordSalt", user.PasswordSalt),
		sql.Named("enabled", user.Enabled),
		sql.Named("id", sql.Out{Dest: &id}),
	)

	if err != nil {
		return 0, errors.Wrap(err, errMsgSPExecFailed)
	}
	if id == 0 {
		return 0, errors.Wrap(err, errMsgSPReturnedAnInvalidOutput)
	}
	return id, nil
}

func (d *userDAL) UpdatePassword(user *domain.User) error {

	_, err := d.sqlDB.Exec("UsersUpdatePassword",
		sql.Named("id", user.ID),
		sql.Named("passwordHash", user.PasswordHash),
		sql.Named("passwordSalt", user.PasswordSalt),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *userDAL) UpdateEnabledState(user *domain.User) error {

	_, err := d.sqlDB.Exec("UsersUpdateEnableState",
		sql.Named("id", user.ID),
		sql.Named("enabled", user.Enabled),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *userDAL) Delete(userID uint) error {

	_, err := d.sqlDB.Exec("UsersDelete",
		sql.Named("id", userID),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *userDAL) Select(id uint) (*domain.User, error) {

	var user domain.User

	row := d.sqlDB.QueryRow("UsersSelect", sql.Named("id", id))
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.PasswordSalt, &user.Enabled)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &user, nil
}

func (d *userDAL) SelectByUsername(username string) (*domain.User, error) {

	var user domain.User

	row := d.sqlDB.QueryRow("UsersSelectByUsername", sql.Named("username", username))
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.PasswordSalt, &user.Enabled)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &user, nil
}

func (d *userDAL) SelectAll() ([]domain.User, error) {

	rows, err := d.sqlDB.Query("UsersSelectAll")

	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var user domain.User
	var users = make([]domain.User, 0)

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(users) == 0 {
		return nil, dal.ErrorDataAccessEmptyResult
	}

	return users, nil
}

func (d *userDAL) Exists(id uint) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("UsersExists",
		sql.Named("id", id),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}

func (d *userDAL) UsernameExists(username string) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("UsersUsernameExists",
		sql.Named("username", username),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}

/*
func (d *userDAL) SelectUserPermissionsByUserID(ID uint) (domain.UserPermissions, error) {

	userPermissions := make(domain.UserPermissions)

	spRes, err := d.sqlDB.Query("UsersSelectAllResourcesPermissionsByUserID", sql.Named("userID", ID))
	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer spRes.Close()

	var rowResource, rowPermission string
	for spRes.Next() {
		err = spRes.Scan(&rowResource, &rowPermission)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		userPermissions[rowResource] = append(userPermissions[rowResource], rowPermission)
	}
	err = spRes.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(userPermissions) == 0 {
		return nil, services.ErrorDataAccessEmptyResult
	}

	return userPermissions, nil
}
*/
