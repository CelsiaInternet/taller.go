package taller

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/linq"
	"github.com/celsiainternet/elvis/msg"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/utility"
	"github.com/go-chi/chi/v5"
)

var Usuarios *linq.Model

func DefineUsuarios(db *jdb.DB) error {
	if err := defineSchema(db); err != nil {
		return console.Panic(err)
	}

	if Usuarios != nil {
		return nil
	}

	Usuarios = linq.NewModel(schemaTaller, "USUARIOS", "Tabla", 1)
	Usuarios.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Usuarios.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Usuarios.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Usuarios.DefineColum(jdb.KEY, "", "VARCHAR(80)", "-1")
	Usuarios.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Usuarios.DefineColum("name", "", "VARCHAR(250)", "")
	Usuarios.DefineColum("description", "", "TEXT", "")
	Usuarios.DefineColum("_data", "", "JSONB", "{}")
	Usuarios.DefineColum("index", "", "INTEGER", 0)
	Usuarios.DefinePrimaryKey([]string{jdb.KEY})
	Usuarios.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"name",
		"index",
	})
	Usuarios.DefineRequired([]string{
		"name:Atributo requerido (name)",
	})
	Usuarios.IntegrityAtrib(true)
	Usuarios.IndexSource(true)
	Usuarios.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.Trigger(linq.BeforeUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.Trigger(linq.AfterUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.Trigger(linq.BeforeDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.Trigger(linq.AfterDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Usuarios.OnListener = func(data et.Json) {
		console.Debug(data.ToString())
	}

	if err := Usuarios.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
*	GetUsuariosById
* @param id string
* @return et.Item, error
**/
func GetUsuariosById(id string) (et.Item, error) {
	result, err := Usuarios.Data().
		Where(Usuarios.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* InsertUsuarios
* @params project_id, id, name, description string
* @params data et.Json
* @return et.Item, error
**/
func InsertUsuarios(project_id, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(project_id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	if !utility.ValidId(id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, jdb.KEY)
	}

	current, err := Usuarios.Data().
		Where(Usuarios.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return et.Item{Ok: false, Result: current.Result}, nil
	}

	id = utility.GenKey(id)
	now := utility.Now()
	data["date_make"] = now
	data["date_update"] = now
	data["project_id"] = project_id
	data[jdb.KEY] = id
	data["name"] = name
	data["description"] = description
	item, err := Usuarios.Insert(data).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* UpSertUsuarios
* @param project_id string
* @param id string
* @param data et.Json
* @return et.Item, error
**/
func UpSertUsuarios(project_id, id, name, description string, data et.Json) (et.Item, error) {
	current, err := InsertUsuarios(project_id, id, name, description, data)
	if err != nil {
		return et.Item{}, err
	}

	if current.Ok {
		return current, nil
	}

	current_state := current.Key("_state")
	if current_state != utility.ACTIVE {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_UPDATE)
	}

	id = current.Str(jdb.KEY)
	now := utility.Now()
	data["date_update"] = now
	data["project_id"] = project_id
	data[jdb.KEY] = id
	data["name"] = name
	data["description"] = description
	result, err := Usuarios.Update(data).
		Where(Usuarios.Column(jdb.KEY).Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* StateUsuarios
* @param id, state string
* @return et.Item, error
**/
func StateUsuarios(id, state string) (et.Item, error) {
	if !utility.ValidId(state) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := Usuarios.Data("_state").
		Where(Usuarios.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	if !current.Ok {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_FOUND)
	}

	old_state := current.Key("_state")
	if old_state == state {
		return et.Item{}, console.AlertF(msg.RECORD_NOT_CHANGE)
	}

	return Usuarios.Update(et.Json{
		"_state": state,
	}).
		Where(Usuarios.Column(jdb.KEY).Eq(id)).
		CommandOne()
}

/**
* DeleteUsuarios
* @param id string
* @return et.Item, error
**/
func DeleteUsuarios(id string) (et.Item, error) {
	return StateUsuarios(id, utility.FOR_DELETE)
}

/**
* AllUsuarios
* @param project_id, state, search string
* @param page, rows int
* @param _select string
* @return et.List, error
**/
func AllUsuarios(project_id, state, search string, page, rows int, _select string) (et.List, error) {
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Usuarios.Data(_select).
			Where(Usuarios.Column("project_id").In("-1", project_id)).
			And(Usuarios.Concat("NAME:", Usuarios.Column("name"), "DESCRIPTION:", Usuarios.Column("description"), "DATA:", Usuarios.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Usuarios.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Usuarios.Data(_select).
			Where(Usuarios.Column("_state").Neg(state)).
			And(Usuarios.Column("project_id").In("-1", project_id)).
			OrderBy(Usuarios.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Usuarios.Data(_select).
			Where(Usuarios.Column("_state").In("-1", state)).
			And(Usuarios.Column("project_id").In("-1", project_id)).
			OrderBy(Usuarios.Column("name"), true).
			List(page, rows)
	} else {
		return Usuarios.Data(_select).
			Where(Usuarios.Column("_state").Eq(state)).
			And(Usuarios.Column("project_id").In("-1", project_id)).
			OrderBy(Usuarios.Column("name"), true).
			List(page, rows)
	}
}

/**
* upSertUsuarios
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) upSertUsuarios(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	project_id := body.Str("project_id")
	id := body.Str("id")
	name := body.Str("name")
	description := body.Str("description")

	result, err := UpSertUsuarios(project_id, id, name, description, body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* getUsuariosById
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) getUsuariosById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := GetUsuariosById(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* stateUsuarios
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) stateUsuarios(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	body, _ := response.GetBody(r)
	state := body.Str("state")

	result, err := StateUsuarios(id, state)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* deleteUsuarios
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) deleteUsuarios(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := DeleteUsuarios(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* allUsuarios
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) allUsuarios(w http.ResponseWriter, r *http.Request) {
	query := response.GetQuery(r)
	project_id := query.Str("project_id")
	state := query.Str("state")
	search := query.Str("search")
	page := query.ValInt(1, "page")
	rows := query.ValInt(30, "rows")
	_select := query.Str("select")

	result, err := AllUsuarios(project_id, state, search, page, rows, _select)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/** Copy this code to router.go
	// Usuarios
	er.ProtectRoute(r, er.Get, "/usuarios/{id}", rt.getUsuariosById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/usuarios", rt.upSertUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/usuarios/state/{id}", rt.stateUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/usuarios/{id}", rt.deleteUsuarios, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/usuarios/all", rt.allUsuarios, PackageName, PackagePath, host)
**/

/** Copy this code to func initModel in model.go
	if err := DefineUsuarios(db); err != nil {
		return console.Panic(err)
	}
**/
