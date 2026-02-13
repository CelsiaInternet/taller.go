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

var Taller *linq.Model

func DefineTaller(db *jdb.DB) error {
	if err := defineSchema(db); err != nil {
		return console.Panic(err)
	}

	if Taller != nil {
		return nil
	}

	Taller = linq.NewModel(schemaTaller, "TALLER", "Tabla", 1)
	Taller.DefineColum("date_make", "", "TIMESTAMP", "NOW()")
	Taller.DefineColum("date_update", "", "TIMESTAMP", "NOW()")
	Taller.DefineColum("_state", "", "VARCHAR(80)", utility.ACTIVE)
	Taller.DefineColum(jdb.KEY, "", "VARCHAR(80)", "-1")
	Taller.DefineColum("project_id", "", "VARCHAR(80)", "-1")
	Taller.DefineColum("name", "", "VARCHAR(250)", "")
	Taller.DefineColum("description", "", "TEXT", "")
	Taller.DefineColum("_data", "", "JSONB", "{}")
	Taller.DefineColum("index", "", "INTEGER", 0)
	Taller.DefinePrimaryKey([]string{jdb.KEY})
	Taller.DefineIndex([]string{
		"date_make",
		"date_update",
		"_state",
		"project_id",
		"name",
		"index",
	})
	Taller.DefineRequired([]string{
		"name:Atributo requerido (name)",
	})
	Taller.IntegrityAtrib(true)
	Taller.IndexSource(true)
	Taller.Trigger(linq.BeforeInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.Trigger(linq.AfterInsert, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.Trigger(linq.BeforeUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.Trigger(linq.AfterUpdate, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.Trigger(linq.BeforeDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.Trigger(linq.AfterDelete, func(model *linq.Model, old, new *et.Json, data et.Json) error {
		return nil
	})
	Taller.OnListener = func(data et.Json) {
		console.Debug(data.ToString())
	}
	
	if err := Taller.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
*	GetTallerById
* @param id string
* @return et.Item, error
**/
func GetTallerById(id string) (et.Item, error) {
	result, err := Taller.Data().
		Where(Taller.Column(jdb.KEY).Eq(id)).
		First()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil	
}

/**
* InsertTaller
* @params project_id, id, name, description string
* @params data et.Json
* @return et.Item, error
**/
func InsertTaller(project_id, id, name, description string, data et.Json) (et.Item, error) {
	if !utility.ValidId(project_id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "project_id")
	}

	if !utility.ValidStr(name, 0, []string{""}) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "name")
	}

	if !utility.ValidId(id) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, jdb.KEY)
	}

	current, err := Taller.Data().
		Where(Taller.Column(jdb.KEY).Eq(id)).
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
	item, err := Taller.Insert(data).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* UpSertTaller
* @param project_id string
* @param id string
* @param data et.Json
* @return et.Item, error
**/
func UpSertTaller(project_id, id, name, description string, data et.Json) (et.Item, error) {
	current, err := InsertTaller(project_id, id, name, description, data)
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
	result, err := Taller.Update(data).
		Where(Taller.Column(jdb.KEY).Eq(id)).
		CommandOne()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* StateTaller
* @param id, state string
* @return et.Item, error
**/
func StateTaller(id, state string) (et.Item, error) {
	if !utility.ValidId(state) {
		return et.Item{}, fmt.Errorf(msg.MSG_ATRIB_REQUIRED, "state")
	}

	current, err := Taller.Data("_state").
		Where(Taller.Column(jdb.KEY).Eq(id)).
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

	return Taller.Update(et.Json{
		"_state":   state,
	}).
		Where(Taller.Column(jdb.KEY).Eq(id)).
		CommandOne()	
}

/**
* DeleteTaller
* @param id string
* @return et.Item, error
**/
func DeleteTaller(id string) (et.Item, error) {
	return StateTaller(id, utility.FOR_DELETE)
}

/**
* AllTaller
* @param project_id, state, search string
* @param page, rows int
* @param _select string
* @return et.List, error
**/
func AllTaller(project_id, state, search string, page, rows int, _select string) (et.List, error) {	
	if state == "" {
		state = utility.ACTIVE
	}

	auxState := state

	if search != "" {
		return Taller.Data(_select).
			Where(Taller.Column("project_id").In("-1", project_id)).
			And(Taller.Concat("NAME:", Taller.Column("name"), "DESCRIPTION:", Taller.Column("description"), "DATA:", Taller.Column("_data"), ":").Like("%"+search+"%")).
			OrderBy(Taller.Column("name"), true).
			List(page, rows)
	} else if auxState == "*" {
		state = utility.FOR_DELETE

		return Taller.Data(_select).
			Where(Taller.Column("_state").Neg(state)).
			And(Taller.Column("project_id").In("-1", project_id)).
			OrderBy(Taller.Column("name"), true).
			List(page, rows)
	} else if auxState == "0" {
		return Taller.Data(_select).
			Where(Taller.Column("_state").In("-1", state)).
			And(Taller.Column("project_id").In("-1", project_id)).
			OrderBy(Taller.Column("name"), true).
			List(page, rows)
	} else {
		return Taller.Data(_select).
			Where(Taller.Column("_state").Eq(state)).
			And(Taller.Column("project_id").In("-1", project_id)).
			OrderBy(Taller.Column("name"), true).
			List(page, rows)
	}
}

/**
* upSertTaller
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) upSertTaller(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	project_id := body.Str("project_id")
	id := body.Str("id")
	name := body.Str("name")
	description := body.Str("description")

	result, err := UpSertTaller(project_id, id, name, description, body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* getTallerById
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) getTallerById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := GetTallerById(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* stateTaller
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) stateTaller(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	body, _ := response.GetBody(r)
	state := body.Str("state")

	result, err := StateTaller(id, state)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* deleteTaller
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) deleteTaller(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := DeleteTaller(id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* allTaller
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (rt *Router) allTaller(w http.ResponseWriter, r *http.Request) {
	query := response.GetQuery(r)
	project_id := query.Str("project_id")
	state := query.Str("state")
	search := query.Str("search")
	page := query.ValInt(1, "page")
	rows := query.ValInt(30, "rows")
	_select := query.Str("select")

	result, err := AllTaller(project_id, state, search, page, rows, _select)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/** Copy this code to router.go
	// Taller
	er.ProtectRoute(r, er.Get, "/taller/{id}", rt.getTallerById, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Post, "/taller", rt.upSertTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Put, "/taller/state/{id}", rt.stateTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Delete, "/taller/{id}", rt.deleteTaller, PackageName, PackagePath, host)
	er.ProtectRoute(r, er.Get, "/taller/all", rt.allTaller, PackageName, PackagePath, host)
**/

/** Copy this code to func initModel in model.go
	if err := DefineTaller(db); err != nil {
		return console.Panic(err)
	}
**/
