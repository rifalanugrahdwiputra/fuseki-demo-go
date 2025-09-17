package main

import (
	"html/template"
	"log"
	"net/http"

	"fuseki-demo/controller"
)

type PageData struct {
	Query string
	Vars  []string
	Rows  [][]string
	Error string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	if r.Method == http.MethodPost {
		query := r.FormValue("query")
		data.Query = query

		var err error
		data.Vars, data.Rows, err = controller.QuerySPARQL(query)
		if err != nil {
			data.Error = err.Error()
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
