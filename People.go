package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var ()

type Database struct {
	dbClient *sql.DB
	log      *logrus.Logger
}
type Person struct {
	UID        string
	canAdd     bool
	drinkFancy bool
}

func CreateDatabase(logger *logrus.Logger) (database Database, err error) {
	database.log = logger
	database.dbClient, err = sql.Open("postgres", "host=pgsql port=5432 user=postgres password=Thunder@01 dbname=postgres sslmode=disable")
	if err != nil {
		return
	}
	err = database.dbClient.Ping()
	if err != nil {
		return
	}
	database.log.Info("Connected to Database")
	return
}
func (db *Database) hasUID(UID string) (Person Person, exists bool) {
	err, DataMap := db.GetPersonData(UID)
	if err != nil {
		return
	}
	exists = true
	//TODO COPY PERSON INFO here
	newUID, ok := DataMap["UID"].(string)
	if ok {
		Person.UID = newUID
	}
	return
}
func (db *Database) AddFriend(UID, NewUID string) {
}
func (db *Database) AddPerson(UID, firstName, lastName string, partyBarner bool) (err error) {
	err = db.RunInsert("INSERT INTO public.users(\"UID\", \"First Name\", \"Last Name\", \"Barner\", \"Added By\")	VALUES ($1, $2, $3, $4, $5);",
		UID, firstName, lastName, partyBarner, "")
	if err != nil {
		db.log.Error(err)
	}
	return
}
func (db *Database) GetPersonData(UID string) (err error, DataMap map[string]interface{}) {
	err, _, DataMapArray := db.RunQueryRows(fmt.Sprintf("SELECT * from \"users\" WHERE \"users\".\"UID\" = '%v'", UID))
	if err != nil {
		return
	}
	db.log.Tracef("DataMapArray: %+v", DataMapArray)
	//TODO if this is greater than 0 database has wrong information len(NodeMapArray)
	if DataMapArray == nil {
		err = errors.New("DataMapArray is null")
		return
	}
	DataMap = DataMapArray[0]
	if DataMap == nil {
		err = errors.New("User is not in database")
		return
	}
	/*
		SerialNumber, ok := NodeMap["SerialNumber"].(string)
		if ok == false {
			err = errors.New("SerialNumber not of type string")
			return
		}
	*/
	return
}

func (db *Database) RunInsert(queryString string, params ...interface{}) (err error) {
	_, err = db.dbClient.Exec(queryString, params...)
	if err != nil {
		db.log.Error(err)
	}
	return
}
func (db *Database) RunQueryRows(queryString string) (error, int, []map[string]interface{}) {
	db.log.Trace("Query: ", queryString)
	var dataMapArray []map[string]interface{}
	rows, err := db.dbClient.Query(queryString)
	if err != nil {
		return err, 0, dataMapArray
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err, 0, dataMapArray
	}

	for rows.Next() {
		rowMap := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return err, 0, dataMapArray
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			rowMap[colName] = *val
		}
		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		dataMapArray = append(dataMapArray, rowMap)
	}
	return nil, len(dataMapArray), dataMapArray
}
