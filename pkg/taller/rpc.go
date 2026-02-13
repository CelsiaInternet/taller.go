package taller

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jrpc"
)

type Services struct{}

func StartRpcServer() {
	pkg, err := jrpc.Load(PackageName)
	if err != nil {
		console.Panic(err)
	}

	services := new(Services)
	err = jrpc.Mount(services)
	if err != nil {
		console.Fatal(err)
	}

	go pkg.Start()
}

func (c *Services) Version(require et.Json, response *et.Item) error {
	company := envar.EnvarStr("", "COMPANY")
	web := envar.EnvarStr("", "WEB")
	version := envar.EnvarStr("", "VERSION")
	help := envar.EnvarStr("", "RPC_HELP")
	response.Ok = true
	response.Result = et.Json{
		"methos":  "RPC",
		"version": version,
		"service": PackageName,
		"host":    HostName,
		"company": company,
		"web":     web,
		"help":    help,
	}

	return console.Rpc(response)
}

func (c *Services) GetTallerById(require et.Json, response *et.Item) error {
	id := require.Str("id")

	result, err := GetTallerById(id)
	if err != nil {
		return err
	}

	*response = result

	return console.Rpc(response)
}	
