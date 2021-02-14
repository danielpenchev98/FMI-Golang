package dbconn

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//PostgresDialectorCreator - creates postgre database specific dialector
func PostgresDialectorCreator(dbDns string) gorm.Dialector {
	return postgres.Open(dbDns)
}
