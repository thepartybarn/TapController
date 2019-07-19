package main

import (
	"encoding/json"
	"io/ioutil"
)

var ()

type Database struct {
	People map[string]Person
}
type Person struct {
	UID     string
	canAdd  bool
	isFancy bool
	isCheap bool
}

func CreateDatabase() (database Database, err error) {
	database.People = make(map[string]Person)
	return
}
func (db *Database) hasUID(UID string) (Person Person, ok bool) {
	Person, ok = db.People[UID]
	return
}
func (db *Database) AddAdmin(UID string) (err error) {
	db.AddPerson(UID, true, false, false)
	return
}
func (db *Database) AddBarner(UID string) (err error) {
	db.AddPerson(UID, false, true, false)
	return
}
func (db *Database) AddFriend(UID string) (err error) {
	db.AddPerson(UID, false, false, true)
	return
}
func (db *Database) Commit() (err error) {
	dataToWrite, err := json.Marshal(db.People)
	if err != nil {
		log.Error("Marshal Error: %v")
		return
	}
	err = ioutil.WriteFile("data/people.json", dataToWrite, 0666)
	if err != nil {
		log.Error("Error Commiting: %v")
	}
	return
}
func (db *Database) Load() (err error) {
	dataToLoad, err := ioutil.ReadFile("data/people.json")
	if err != nil {
		log.Error("Read file: %v")
	}
	err = json.Unmarshal(dataToLoad, &db)
	if err != nil {
		log.Error("Error Loading: %v")
	}
	return
}
func (db *Database) AddPerson(UID string, isAdmin, isBarner, isFriend bool) (err error) {
	_, exists := db.People[UID]
	if UID != "" && !exists {
		toAdd := Person{UID, false, false, true}
		db.People[UID] = toAdd
	}
	return
}
