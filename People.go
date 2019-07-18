package main

import ()

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

func CreateDatabase() (database *Database, err error) {
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
func (db *Database) AddPerson(UID string, isAdmin, isBarner, isFriend bool) (err error) {
	_, exists := db.People[UID]
	if UID != "" && !exists {
		toAdd := Person{UID, false, false, true}
		db.People[UID] = toAdd
	}
	return
}
