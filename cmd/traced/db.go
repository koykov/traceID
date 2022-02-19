package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
)

var dbi *sql.DB

func dbConnect(addr string) (err error) {
	var di int
	if di = strings.Index(addr, "://"); di == -1 {
		err = fmt.Errorf("couldn't get driver name from DSN '%s'", addr)
		return
	}
	drv := addr[:di]
	if len(drv) == 0 {
		return errors.New("empty DB driver")
	}
	addr = addr[di+3:]
	if dbi, err = sql.Open(drv, addr); err != nil {
		return
	}
	if err = dbi.Ping(); err != nil {
		return
	}
	return
}

func dbClose() error {
	if dbi == nil {
		return nil
	}
	return dbi.Close()
}
