package dataaccess

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type resourcesDAL struct {
	sqlDB *sql.DB
}

func NewResourcesDAL(sqlDB *sql.DB) dal.ResourcesDAL {
	return &resourcesDAL{
		sqlDB: sqlDB,
	}
}

func (d *resourcesDAL) Insert(resource *domain.Resource) (uint, error) {

	var id uint

	_, err := d.sqlDB.Exec("ResourcesInsert",
		sql.Named("path", resource.Path),
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

func (d *resourcesDAL) Delete(id uint) error {

	_, err := d.sqlDB.Exec("ResourcesDelete",
		sql.Named("id", id),
	)

	if err != nil {
		return errors.Wrap(err, errMsgSPExecFailed)
	}
	return nil
}

func (d *resourcesDAL) Select(id uint) (*domain.Resource, error) {

	var resource domain.Resource

	row := d.sqlDB.QueryRow("ResourcesSelect", sql.Named("id", id))
	err := row.Scan(&resource.ID, &resource.Path)

	if err != nil {
		if err == sql.ErrNoRows {
			// the caller needs to use errors.Is() to check which err is returned.
			//return nil, errors.Wrap(err, services.ErrorDataAccessEmptyResult.Error())
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &resource, nil
}

func (d *resourcesDAL) SelectByPath(path string) (*domain.Resource, error) {

	var resource domain.Resource

	row := d.sqlDB.QueryRow("ResourcesSelectByPath", sql.Named("path", path))
	err := row.Scan(&resource.ID, &resource.Path)

	if err != nil {
		if err == sql.ErrNoRows {
			// the caller needs to use errors.Is() to check which err is returned.
			//return nil, errors.Wrap(err, services.ErrorDataAccessEmptyResult.Error())
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &resource, nil
}

func (d *resourcesDAL) SelectAll() ([]dto.Resource, error) {

	rows, err := d.sqlDB.Query("ResourcesSelectAll")

	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var resources = make([]dto.Resource, 0)

	for rows.Next() {
		var resource domain.Resource

		err = rows.Scan(&resource.ID, &resource.Path)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		resources = append(resources, dto.Resource{ID: &resource.ID, Path: &resource.Path})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(resources) == 0 {
		return nil, dal.ErrorDataAccessEmptyResult
	}

	return resources, nil
}

func (d *resourcesDAL) Exists(id uint) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("ResourcesExists",
		sql.Named("id", id),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}

func (d *resourcesDAL) PathExists(path string) (bool, error) {

	var exists bool

	_, err := d.sqlDB.Exec("ResourcesPathExists",
		sql.Named("path", path),
		sql.Named("exists", sql.Out{Dest: &exists}),
	)

	if err != nil {
		return false, errors.Wrap(err, errMsgSPExecFailed)
	}
	return exists, nil
}
