package users

import (
	"sync"

	"github.com/UnnoTed/authenticaTed/errors"

	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/postgresql"
)

const Table = `users`
const TableBan = `user_bans`
const TableEvents = `user_events`
const TableActivation = `user_activation`

var (
	session sqlbuilder.Database

	uc db.Collection
	bc db.Collection
	ac db.Collection
	ec db.Collection

	isTest   = false
	settings = postgresql.ConnectionURL{
		Database: `authed_test`,
		Host:     `localhost`,
		User:     `authed_test`,
		Password: `123`,
	}
)

// Connect to the database
func Connect() *errors.Error {
	var err error
	logger.Info("[User.Database]: Connecting...")

	// connect to the database
	session, err = postgresql.Open(settings)
	if err != nil {
		return errors.FromErr(err)
	}
	s := Schema

	if isTest {
		s = append(Schema, SchemaTest...)
	}

	// exec the schema
	Exec(s)

	// get tables / collections

	// users
	uc = session.Collection(Table)
	CheckCollection(uc, Table)

	// ban
	bc = session.Collection(TableBan)
	CheckCollection(bc, TableBan)

	// activation
	ac = session.Collection(TableActivation)
	CheckCollection(ac, TableActivation)

	// events
	ec = session.Collection(TableEvents)
	CheckCollection(ec, TableEvents)

	return nil
}

func CheckCollection(c db.Collection, table string) {
	if !c.Exists() {
		logger.WithField("table", table).Fatal(`[User.Database]: there is no table/collection!`)
	} else {
		logger.Info("[User.Database]: Connected!")
	}
}

// Disconnect from the database
func Disconnect() *errors.Error {
	logger.Info("[User.Database]: Disconnecting...")

	err := session.Close()
	if err != nil {
		logger.WithError(err).Error("[User.Database]: Error while disconnecting from the database")
		return errors.FromErr(err)
	}

	logger.Info("[User.Database]: Disconnected!")
	return nil
}

func Exec(list []string) {
	// without the WaitGroup it will only Exec the first one
	var wg sync.WaitGroup
	for _, s := range list {
		wg.Add(1)

		go func(s string, wg *sync.WaitGroup) {
			_, err := session.Exec(s)
			if err != nil {
				panic(err)
			}

			wg.Done()
		}(s, &wg)
	}

	wg.Wait()
}
