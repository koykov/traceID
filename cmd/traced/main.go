package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

var (
	fcnf string
	cnf  *Config
)

func init() {
	var err error

	rf := func(v *string, names []string, value, usage string) {
		for i := range names {
			flag.StringVar(v, names[i], value, usage)
		}
	}
	rf(&fcnf, []string{"config", "c"}, "", "Path to config file.")
	flag.Parse()

	if len(fcnf) == 0 {
		log.Fatalln("param -config is required")
	}
	if _, err = os.Stat(fcnf); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("config file '%s' doesn't exists", fcnf)
	}
	if cnf, err = ParseConfig(fcnf); err != nil {
		log.Fatalf("error '%s' caught on parse config '%s'", err.Error(), fcnf)
	}
	if len(cnf.Listeners) == 0 {
		log.Fatalln("no listeners available")
	}
}

func main() {
	_ = cnf
	// ...
}
