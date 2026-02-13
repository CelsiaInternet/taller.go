package main

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/taller/a"
	"github.com/go-chi/chi/v5"
)

func main() {
	// var name string
	// var edad = 48
	// ciudad := "San Gil"
	// name = "Hello Cesar mi edad es:" + fmt.Sprintf("%d", edad) + " y vivo en " + ciudad
	// fmt.Println(name)
	// fmt.Println(a.Empresa)
	// _, err := a.Hello("Cesar", "San Gil", 48)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(message)
	server()
}

func server() {
	r := chi.NewRouter()
	r.Post("/hello", a.HandlerHello)
	fmt.Println("Server running on port 3001")
	http.ListenAndServe(":3001", r)
}
