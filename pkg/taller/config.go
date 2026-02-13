package taller

import (
	"fmt"

	"github.com/celsiainternet/elvis/config"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jrpc"
)

func LoadConfig() error {
	StartRpcServer()

	// stage := envar.GetStr("local", "STAGE")
	// return defaultConfig(stage)
	return nil
}

func defaultConfig(stage string) error {
	name := "default"
	result, err := jrpc.CallItem("Module.Services.GetConfig", et.Json{
		"stage": stage,
		"name":  name,
	})
	if err != nil {
		return err
	}

	if !result.Ok {
		return fmt.Errorf(jrpc.MSG_NOT_LOAD_CONFIG, stage, name)
	}

	cfg := result.Json("config")
	return config.Load(cfg)
}
