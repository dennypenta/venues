package storages

import (
	"venues/cmd/settings"

	"gopkg.in/mgo.v2"

	"log"
)

var storage *mgo.Database

func initStorage(db string) *mgo.Database {
	dialUrl := settings.MustGetSetting("MONGO_ADDRESS")

	if storage == nil {
		session, err := mgo.Dial(dialUrl)
		if err != nil {
			log.Fatal("Error initializing Storage \n %s", err.Error())
		}

		storage = session.DB(db)
	}

	return storage
}

func GetStorage() *mgo.Database {
	return initStorage(settings.MustGetSetting("MONGO_DB_NAME"))
}

func GetTestStorage() *mgo.Database {
	return initStorage(settings.MustGetSetting("MONGO_DB_NAME_TEST"))
}
