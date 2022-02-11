package main

import "database/sql"

var dbi *sql.DB

func dbConnect(addr string) error {
	dbi = new(sql.DB)
	_ = dbi
	return nil
}
