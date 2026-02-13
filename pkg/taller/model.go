package taller

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/jdb"
)

func initModels(db *jdb.DB) error {
	if err := DefineTaller(db); err != nil {
		return console.Panic(err)
	}

	return nil
}
