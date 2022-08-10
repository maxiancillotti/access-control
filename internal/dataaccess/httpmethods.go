package dataaccess

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type httpmethodsDAL struct {
	sqlDB *sql.DB
}

func NewHttpMethodsDAL(sqlDB *sql.DB) dal.HttpMethodsDAL {
	return &httpmethodsDAL{
		sqlDB: sqlDB,
	}
}

func (d *httpmethodsDAL) SelectByName(name string) (*domain.HttpMethod, error) {

	var method domain.HttpMethod

	row := d.sqlDB.QueryRow("HttpMethodsSelectByName", sql.Named("name", name))
	err := row.Scan(&method.ID, &method.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			// the caller needs to use errors.Is() to check which err is returned.
			//return nil, errors.Wrap(err, services.ErrorDataAccessEmptyResult.Error())
			return nil, dal.ErrorDataAccessEmptyResult
		}
		return nil, errors.Wrap(err, errMsgSPResultScanFailed)
	}
	return &method, nil
}

func (d *httpmethodsDAL) SelectAll() ([]dto.HttpMethod, error) {

	rows, err := d.sqlDB.Query("HttpMethodsSelectAll")

	if err != nil {
		return nil, errors.Wrap(err, errMsgSPExecFailed)
	}
	defer rows.Close()

	var methods = make([]dto.HttpMethod, 0)

	for rows.Next() {
		var method domain.HttpMethod

		err = rows.Scan(&method.ID, &method.Name)
		if err != nil {
			return nil, errors.Wrap(err, errMsgSPResultScanFailed)
		}
		methods = append(methods, dto.HttpMethod{ID: &method.ID, Name: &method.Name})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, errMsgRowsScanReturnedAnError)
	}

	if len(methods) == 0 {
		return nil, dal.ErrorDataAccessEmptyResult
	}

	return methods, nil
}
