package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var templates = template.Must(template.ParseGlob("templates/*"))

//Here I create the struct as is in the database to be able to manage the information
type Employee struct {
	Id    int
	Name  string
	Email string
}

func Home(w http.ResponseWriter, r *http.Request) {

	stablishedConection := conectionDB()
	values, err := stablishedConection.Query("SELECT * FROM employees")
	if err != nil {
		panic(err.Error())
	}
	employee := Employee{}
	arrayEmployee := []Employee{}

	for values.Next() { //Here I get the whole information and asign in this variables
		var id int
		var name, email string
		err = values.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email

		arrayEmployee = append(arrayEmployee, employee)
	}
	//fmt.Println(arrayEmployee)

	// fmt.Fprintf(w, "Hello dev")
	templates.ExecuteTemplate(w, "home", arrayEmployee)
}

//Add a new data
func Add(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "add", nil)
}

//Insert the data in the table from de database
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		stablishedConection := conectionDB()
		insertQuery, err := stablishedConection.Prepare("INSERT INTO employees(name, email) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}

		insertQuery.Exec(name, email)

		http.Redirect(w, r, "/", 301)
	}
}

//Delete the register selected
func Delete(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	// fmt.Println(idEmployee)

	stablishedConection := conectionDB()
	deleteData, err := stablishedConection.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	deleteData.Exec(idEmployee)

	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	employee := Employee{}
	// fmt.Println(idEmployee)

	stablishedConection := conectionDB()
	value, err := stablishedConection.Query("SELECT * FROM employees WHERE id=?", idEmployee)

	for value.Next() { //Here I get the whole information and asign in this variables
		var id int
		var name, email string
		err = value.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email
	}

	// fmt.Println(employee)
	templates.ExecuteTemplate(w, "edit", employee)

}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		stablishedConection := conectionDB()
		editQuery, err := stablishedConection.Prepare("UPDATE employees set Name = ?, email = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}

		editQuery.Exec(name, email, id)
		http.Redirect(w, r, "/", 301)
	}
}

//Create the conection to the data base
func conectionDB() (conection *sql.DB) {
	driver := "mysql"
	user := "root"
	password := ""
	name := "system"

	conection, err := sql.Open(driver, user+":"+password+"@tcp(localhost:3308)/"+name)
	if err != nil {
		panic(err.Error())
	}
	return conection
}

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	log.Println("Running server...")
	http.ListenAndServe(":8080", nil)

}
