package main

import (
	users "crud_3/controllers"
	"html/template"
	"log"
	"net/http"
	"path"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// log process
	log.Printf(r.URL.Path)

	tmpl, err := template.ParseFiles(
		path.Join("views", "home.html"),
		path.Join("views", "layouts/navbar.html"),
		path.Join("views", "layouts/main.html"),
	)

	// if error occur
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

	// execute template
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

}

func main() {

	// mux is a http router
	mux := http.NewServeMux()

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/users", users.Index)
	mux.HandleFunc("/users/show", users.Show)
	mux.HandleFunc("/users/create", users.Create)
	mux.HandleFunc("/users/store", users.Store)
	mux.HandleFunc("/users/edit", users.Edit)
	mux.HandleFunc("/users/update", users.Update)
	mux.HandleFunc("/users/", users.Delete)
	// mux.HandleFunc("/posts", posts.PostHandlers)

	// create log process
	log.Println("Run on port : 9000")

	// config the server
	err := http.ListenAndServe(":9000", mux)

	// create handler for handle an error
	log.Fatal(err)

}
