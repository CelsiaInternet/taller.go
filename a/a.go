package a

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/celsiainternet/taller/b"
)

const (
	Empresa = "Celsia Internet"
)

type J map[string]interface{}

type Persona struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

func Hello(name, city string, age int) ([]string, error) {
	if age < 1 {
		return nil, fmt.Errorf("Edad no puede ser menor a 1")
	}

	var result []string
	r := 0
	for i := 0; i < age; i++ {
		r = b.Add(r, i)
		result = append(result, fmt.Sprintf("Hello %s, tengo %d años y vivo en %s y trabajo en %s", name, r, city, Empresa))
		// fmt.Println(result[i])
	}

	return result, nil
}

func HandlerHello(w http.ResponseWriter, r *http.Request) {
	bt, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var p Persona
	err = json.Unmarshal(bt, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := Hello(p.Name, p.City, p.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResult := make(J)
	for i, item := range items {
		jsonResult[fmt.Sprintf("result_%d", i)] = item
	}

	result, err := json.Marshal(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(result)
}
