package taller

import (
	"context"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
)

type Controller struct {
	Db *jdb.DB
}

func (c *Controller) Version(ctx context.Context) (et.Json, error) {
	company := envar.GetStr("", "COMPANY")
	web := envar.GetStr("", "WEB")
	version := envar.EnvarStr("0.0.1", "VERSION")
  service := et.Json{
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    "",
	}

	return service, nil
}

func (c *Controller) Init(ctx context.Context) {
	initModels(c.Db)
	initEvents()
}

type Repository interface {
	Version(ctx context.Context) (et.Json, error)
	Init(ctx context.Context)
}
