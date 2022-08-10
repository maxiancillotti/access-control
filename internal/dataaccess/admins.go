package dataaccess

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type adminsDAL struct {
	sqlDB *sql.DB
}

func NewAdminsDAL(sqlDB *sql.DB) dal.AdminsDAL {
	return &adminsDAL{
		sqlDB: sqlDB,
	}
}

func (d *adminsDAL) Insert(admin *domain.Admin) (uint, error) {

	var id uint

	_, err := d.sqlDB.Exec("AdminsInsert",
		sql.Named("username", admin.Username),
		sql.Named("passwordHash", admin.PasswordHash),
		sql.Named("passwordSalt", admin.PasswordSalt),
		sql.Named("enabled", admin.Enabled),
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

func (d *adminsDAL) UpdatePassword(admin *domain.Admin) error {

	_, err := d.sqlDB.Exec("AdminsUpdatePassword",
		sql.Named("id", admin.ID),
		sql.Named("passwordHash", admin.PasswordHash),
		sql.Named("passwordSalt", admin.PasswordSalt),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *adminsDAL) UpdateEnabledState(admin *domain.Admin) error {

	_, err := d.sqlDB.Exec("AdminsUpdateEnableState",
		sql.Named("id", admin.ID),
		sql.Named("enabled", admin.Enabled),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *adminsDAL) Delete(id uint) error {

	_, err := d.sqlDB.Exec("AdminsDelete",
		sql.Named("id", id),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *adminsDAL) Select(id uint) (*domain.Admin, error) {

	var admin domain.Admin

	row := d.sqlDB.QueryRow("AdminsSelect", sql.Named("id", id))
	err := row.Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.PasswordSalt, &admin.Enabled)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &admin, nil
}

func (d *adminsDAL) SelectByUsername(username string) (*domain.Admin, error) {

	var admin domain.Admin

	row := d.sqlDB.QueryRow("AdminsSelectByUsername", sql.Named("username", username))
	err := row.Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.PasswordSalt, &admin.Enabled)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &admin, nil
}

func (d *adminsDAL) SelectAll() ([]domain.Admin, error) {

	rows, err := d.sqlDB.Query("AdminsSelectAll")

	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var admin domain.Admin
	var admins = make([]domain.Admin, 0)

	for rows.Next() {
		err = rows.Scan(&admin.ID, &admin.Username)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		admins = append(admins, admin)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(admins) == 0 {
		return nil, dal.ErrorDataAccessEmptyResult
	}

	return admins, nil
}

func (d *adminsDAL) Exists(id uint) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("AdminsExists",
		sql.Named("id", id),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}

func (d *adminsDAL) UsernameExists(username string) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("AdminsUsernameExists",
		sql.Named("username", username),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}
