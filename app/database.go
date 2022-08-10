package app

import (
	"database/sql"

	"github.com/maxiancillotti/access-control/app/config"
	"github.com/maxiancillotti/mssqlconn"
)

func GetDBConn(dbconfig *config.DatabaseConfig) *sql.DB {

	return mssqlconn.NewBuilder().
		SetHostname(dbconfig.Hostname).
		SetPort(dbconfig.Port).
		SetInstance(dbconfig.Instance).
		SetDatabaseName(dbconfig.DBName).
		SetCredentials(dbconfig.User, dbconfig.Password).
		//EnableDebug().
		Build().
		OpenConn()
}
