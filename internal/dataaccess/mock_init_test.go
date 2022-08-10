package dataaccess

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	sqlDBConn, sqlMock = initSqlMock()

	adminsDALMock       = NewAdminsDAL(sqlDBConn)
	userDALMock         = NewUserDAL(sqlDBConn)
	userRESTPermDALMock = NewUserRESTPermissionsDAL(sqlDBConn)
	httpMethodsDALMock  = NewHttpMethodsDAL(sqlDBConn)
	resourcesDALMock    = NewResourcesDAL(sqlDBConn)
)

func initSqlMock() (*sql.DB, sqlmock.Sqlmock) {

	var t testing.T

	sqlDBConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("cannot initialize mock:", err)
	}
	return sqlDBConn, mock
}
