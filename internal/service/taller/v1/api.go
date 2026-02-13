package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celsiainternet/elvis/cache"
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/jdb"
	"github.com/celsiainternet/elvis/jrpc"
	"github.com/celsiainternet/elvis/utility"
	"github.com/dimiro1/banner"
	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-colorable"
	pkg "github.com/celsiainternet/taller/pkg/taller"	
)

func New() http.Handler {
	r := chi.NewRouter()

	err := pkg.LoadConfig()
	if err != nil {
		console.Panic(err)
	}

	_, err = cache.Load()
	if err != nil {
		console.Panic(err)
	}

	_, err = event.Load()
	if err != nil {
		console.Panic(err)
	}

	db, err := jdb.Load()
	if err != nil {
		console.Panic(err)
	}

	_pkg := &pkg.Router{
		Repository: &pkg.Controller{
			Db: db,
		},
	}

	r.Mount(pkg.PackagePath, _pkg.Routes())

	return r
}

func Close() {
	jrpc.Close()
	cache.Close()
	event.Close()
}

func Banner() {
	time.Sleep(3 * time.Second)
	templ := utility.BannerTitle(pkg.PackageName, 4)
	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
	fmt.Println()
}
