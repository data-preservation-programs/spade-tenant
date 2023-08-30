package main

import "github.com/data-preservation-programs/spade-tenant/db"

func main() {
	dbDsn := "postgres://postgres:password@localhost:5432/spade-tenant"
	debug := true

	db.OpenDatabase(dbDsn, debug)
}
