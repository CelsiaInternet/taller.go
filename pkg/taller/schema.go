package taller

import (
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"	
)

var schemaTaller *linq.Schema

func defineSchema(db *jdb.DB) error {
	if schemaTaller == nil {
		schemaTaller = linq.NewSchema(db, "taller")
	}

	return nil
}
