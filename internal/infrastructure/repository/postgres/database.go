package postgres

import "github.com/istovpets/pgxhelper/sqlsetpgxhelper"

type Database struct {
	*sqlsetpgxhelper.DBHelper
}

func New(db *sqlsetpgxhelper.DBHelper) *Database {
	return &Database{
		DBHelper: db,
	}
}
