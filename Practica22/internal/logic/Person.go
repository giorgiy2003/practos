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

