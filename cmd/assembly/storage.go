package assembly

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

		storage = session.DB(settings.MustGetSetting("MONGO_DB_NAME"))
	}

	return storage
}

func InitStorage() *mgo.Database {
	return initStorage(settings.MustGetSetting("MONGO_DB_NAME"))
}

func InitTestStorage() *mgo.Database {
	return initStorage(settings.MustGetSetting("MONGO_DB_NAME_TEST"))
}
