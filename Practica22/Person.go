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

//GET /person/{id} – возвращает одну модель Person.
func (p *person) FindPersonByID(db *gorm.DB, Id int) (*person, error) {
	var err error
	err = db.Debug().Model(person{}).Where("Id = ?", Id).Take(&p).Error
	if err != nil {
		return &person{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &person{}, errors.New("person Not Found")
	}
	return p, err
}

//DELETE /person/{id} – удаляет модель Person
func (p *person) DeletePerson(db *gorm.DB, Id int) (int64, error) {

	db = db.Debug().Model(&person{}).Where("Id = ?", Id).Take(&person{}).Delete(&person{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//PUT /person/{id} – обновляет модель Person
func (p *person) UpdateAUser(db *gorm.DB, Id int) (*person, error) {

	var err error
	db = db.Debug().Model(&person{}).Where("Id = ?", Id).Take(&person{}).UpdateColumns(
		map[string]interface{}{
			"Id":        p.Id,
			"Email":     p.Email,
			"FirstName": p.FirstName,
			"LastName":  p.LastName,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &person{}, db.Error
	}
	// Отоброжение обновленного пользователя
	err = db.Debug().Model(&person{}).Where("Id = ?", Id).Take(&p).Error
	if err != nil {
		return &person{}, err
	}
	return p, nil
}

//GET /person/ - возвращает список моделей Person о которых у нас есть инфа.
func (p *person) FindAllPersons(db *gorm.DB) (*[]person, error) {
	var err error
	persons := []person{}
	err = db.Debug().Model(&person{}).Limit(100).Find(&persons).Error
	if err != nil {
		return &[]person{}, err
	}
	return &persons, err
}

func handleRequest() {
	http.HandleFunc("/", home_page)
	http.ListenAndServe(":8090", nil)
}

func main() {
	handleRequest()
}
