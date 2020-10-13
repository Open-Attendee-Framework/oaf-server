package packer

import "github.com/gobuffalo/packr"

//PostgresBox holds the psql migrations
var PostgresBox packr.Box

//SQLiteBox holds the sqlite migrations
var SQLiteBox packr.Box

func PackAll() {
	PostgresBox = packr.NewBox("../../scripts/sql/migrations/postgres")
	SQLiteBox = packr.NewBox("../../scripts/sql/migrations/sqlite")
}
