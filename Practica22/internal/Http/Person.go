package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p person) formatStyle() string {
	return fmt.Sprintf("%d. Person name is %s. His lastname is %s. His phone is %s. Also email: %s", p.Id, p.FirstName, p.LastName, p.Phone, p.Email)
}

func home_page(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "./persons.db") //Для открытия соединения с базой данных используем функцию sql.Open()
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS persons (Id INTEGER PRIMARY KEY,email TEXT,phone TEXT,firstName TEXT,lastName TEXT)")
	statement.Exec()
	//statement, _ = database.Prepare("INSERT INTO persons (email,phone,firstName,lastName) VALUES (?,?,?,?)")
	//statement.Exec("Masha@email", "+79545454767", "masha", "veselova")
	rows, _ := database.Query("SELECT Id,email,phone,firstName,lastName FROM persons")
	persons := []person{}
	for rows.Next() {
		p := person{}
		rows.Scan(&p.Id, &p.Email, &p.Phone, &p.FirstName, &p.LastName)
		persons = append(persons, p)
	}

	for _, p := range persons {
		fmt.Fprintln(w, p.formatStyle())
	}
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.ListenAndServe(":8090", nil)
}

func main() {
	handleRequest()
}
