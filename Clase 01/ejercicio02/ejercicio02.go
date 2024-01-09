package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Persona struct {
	Nombre   string `json:"firstName"`
	Apellido string `json:"lastName"`
}

func main() {
	router := chi.NewRouter()
	router.Post("/greetings", func(w http.ResponseWriter, r *http.Request) {
		var p Persona
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			panic(err)
		}
		respuesta := fmt.Sprintf("Hola %s %s", p.Nombre, p.Apellido)
		w.WriteHeader(200)
		w.Write([]byte(respuesta))
	})
	http.ListenAndServe(":8000", router)
}
