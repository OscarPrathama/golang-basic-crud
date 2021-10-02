package users

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT*FROM employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	tmpl, err := template.ParseFiles(
		path.Join("views", "users/index.html"), // harus diletakkan paling atas
		path.Join("views", "layouts/main.html"),
		path.Join("views", "layouts/navbar.html"),
	)

	// if error occur
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, res)
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	tmpl, err := template.ParseFiles(
		path.Join("views", "users/show.html"), // harus diletakkan paling atas
		path.Join("views", "layouts/main.html"),
		path.Join("views", "layouts/navbar.html"),
	)

	// if error occur
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

	// tmpl.ExecuteTemplate(w, "Show", emp)
	tmpl.Execute(w, emp)
	defer db.Close()
}

func Create(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		path.Join("views", "users/create.html"), // harus diletakkan paling atas
		path.Join("views", "layouts/main.html"),
		path.Join("views", "layouts/navbar.html"),
	)

	// if error occur
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Store(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")

		insForm, err := db.Prepare("INSERT INTO employee(name, city) VALUES (?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name : " + name + ", City : " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/users", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT*FROM employee WHERE id = ?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	tmpl, err := template.ParseFiles(
		path.Join("views", "users/edit.html"), // harus diletakkan paling atas
		path.Join("views", "layouts/main.html"),
		path.Join("views", "layouts/navbar.html"),
	)

	// if error occur
	if err != nil {
		log.Println(err)
		http.Error(w, "500, Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, emp)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + ", City: " + city)
		defer db.Close()
		http.Redirect(w, r, "/users/edit?id="+id, 301)
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE user")
	defer db.Close()
	http.Redirect(w, r, "/users", 301)
}
