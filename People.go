package main

import (
	"fmt"

	"gorm.io/gorm"
)

var ()

type Person struct {
	gorm.Model
	UID        string `gorm:"unique;not null"`
	FirstName  string
	LastName   string
	AddedBy    string
	CanAdd     bool
	DrinkFancy bool
}

func (database *DatabaseConnection) GetPersonData(UID string) (record Person, err error) {
	result := database.sqlClient.Where("UID = ?", UID).Find(&record)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected != 1 {
		err = fmt.Errorf("found no one!")
	}

	return
}
func (database *DatabaseConnection) AddFriend(UID, NewUID string) {
}
func (database *DatabaseConnection) AddPerson(record Person) (err error) {
	result := database.sqlClient.Create(record)
	err = result.Error
	return
}
