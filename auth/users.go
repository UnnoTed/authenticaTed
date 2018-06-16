package users

import (
	"github.com/UnnoTed/authenticaTed/errors"
	. "github.com/UnnoTed/authenticaTed/logger"
	db "upper.io/db.v2"
)

// Find gets all users in the database with the given condition
func Find(cond ...interface{}) ([]*User, *errors.Error) {
	Logger.Debug("[Users.Find]: Finding Users...")
	var list []*User
	if err := uc.Find(cond...).All(&list); err != nil {
		return nil, errors.FromErr(err)
	}

	return list, nil
}

func FindExpressive(cond ...interface{}) db.Result {
	return uc.Find(cond...)
}

func SoftDelete() {

}

func HardDelete() {

}

func Contact() {

}
