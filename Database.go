package main

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	sqlClient  *gorm.DB
	log        *logrus.Logger
	serialData sync.Mutex
}

func SetupDatabaseConnections(logger *logrus.Logger) (database *DatabaseConnection, err error) {
	database = new(DatabaseConnection)
	database.log = logger
	database.log.Trace("Connecting to databases")

	database.log.Trace("Connecting to postgres")
	database.sqlClient, err = gorm.Open(postgres.Open("host=eieio_postgres_1 port=5432 user=postgres password=Thunder@01 dbname=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		return
	}
	database.log.Trace("Connected to postgres")
	database.log.Trace("Auto migrating postgres")
	//Setup postgres tables here
	err = database.sqlClient.AutoMigrate(&Person{})
	if err != nil {
		return
	}
	database.log.Trace("Done Auto migrating postgres")

	return
}
func (database *DatabaseConnection) Close() (err error) {

	return
}
