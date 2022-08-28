package conn

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Database

func InitializeDB() error {
	// Mongodb
	host := os.Getenv("MONGO_HOST")
	// dbName := os.Getenv("MONGO_DB_NAME")
	dbName := "login"

	session, err := mgo.Dial(host)
	if err != nil {
		return err
	}
	db = session.DB(dbName)
	return nil
}

// GetMongoDB function to return DB connection
func GetMongoDB() *mgo.Database {
	return db
}
