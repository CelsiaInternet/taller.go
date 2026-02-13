package taller

import (
	"context"
	"net/http"
	"os"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	er "github.com/celsiainternet/elvis/router"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
)

var PackageName = "taller"
var PackageTitle = "taller"
var PackagePath = envar.GetStr("/api/taller", "PATH_URL")
var PackageVersion = envar.EnvarStr("0.0.1", "VERSION")
var HostName, _ = os.Hostname()

type Router struct {
	Repository Repository
}

func (rt *Router) Routes() http.Handler {
	defaultHost := strs.Format("http://%s", HostName)
	var host = strs.Format("%s:%d", envar.GetStr(defaultHost, "HOST"), envar.GetInt(3300, "PORT"))

	r := chi.NewRouter()

	er.PublicRoute(r, er.Get, "/version", rt.version, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/routes", rt.routes, PackageName, PackagePath, host)
	// Taller
	er.ProtectRoute(r, er.Get, "/{id}", rt.getTallerById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/", rt.upSertTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/state/{id}", rt.stateTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/{id}", rt.deleteTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/", rt.allTaller, PackageName, PackagePath, host)
	// Users
	er.ProtectRoute(r, er.Get, "/usuarios/{id}", rt.getUsuariosById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/usuarios", rt.upSertUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/usuarios/state/{id}", rt.stateUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/usuarios/{id}", rt.deleteUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/usuarios/all", rt.allUsuarios, PackageName, PackagePath, host)

	ctx := context.Background()
	rt.Repository.Init(ctx)
	middleware.SetServiceName(PackageName)

	console.LogKF(PackageName, "Router version:%s", PackageVersion)
	return r
}

func (rt *Router) version(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := rt.Repository.Version(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

func (rt *Router) routes(w http.ResponseWriter, r *http.Request) {
	_routes := er.GetRoutes()
	routes := []et.Json{}
	for _, route := range _routes {
		routes = append(routes, et.Json{
			"method": route.Str("method"),
			"path":   route.Str("path"),
		})
	}

	result := et.Items{
		Ok:     true,
		Count:  len(routes),
		Result: routes,
	}

	response.ITEMS(w, r, http.StatusOK, result)
}
