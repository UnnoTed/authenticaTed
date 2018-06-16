package users

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/c2h5oh/hide"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"upper.io/db.v2"

	"github.com/UnnoTed/authenticaTed/errors"
	. "github.com/UnnoTed/authenticaTed/logger"
	"github.com/UnnoTed/authenticaTed/util"
)

// cache for a nil time
var nilTime = time.Time{}

// User holds all needed user information
// also includes validation and db management
type User struct {
	ID hide.Int64 `db:"id,omitempty" json:"id,string"`

	Name     string `db:"name"       json:"name"               valid:"optional,length(3|20),alphanum"`
	Username string `db:"username"   json:"username"           valid:"optional,length(3|25),matches(^[a-zA-Z0-9_]+$)"`
	Password string `db:"password"   json:"password,omitempty" valid:"optional,length(3|255)"`
	Email    string `db:"email"      json:"email"              valid:"optional,length(6|255),email"`

	Token string `db:"-"             json:"token"` // jwt
	Power int    `db:"power"         json:"power"`

	Deleted bool      `db:"deleted"  json:"deleted"`
	Created time.Time `db:"created"  json:"created"`
	Seen    time.Time `db:"seen"     json:"seen"`

	//

	LastName string `db:"last_name" valid:"optional,length(3|50),alphanum"`

	// cache, it will only look for activation codes
	// when this is set to false
	Activated bool `db:"activated"   json:"activated"`

	// other structs
	Banned     *Ban        `db:"-"   json:"banned"`
	Activation *Activation `db:"-"   json:"activation"`
}

// Ban information of a user
type Ban struct {
	ID     hide.Int64 `db:"id,omitempty" json:"id,string"`
	UserID hide.Int64 `db:"user_id"      json:"user_id,string"`

	State     bool `db:"state"           json:"state"`
	Temporary bool `db:"temporary"       json:"temporary"`

	Starts time.Time `db:"starts"        json:"starts"`
	Until  time.Time `db:"until"         json:"until"`
}

// Activation code for a user
type Activation struct {
	ID     int64  `db:"id,omitempty" json:"id,string"`
	UserID int64  `db:"user_id"      json:"user_id,string"`
	Code   string `db:"code"         json:"code"`
}

// NewUser creates a new user
func NewUser() *User {
	u := new(User)
	return u
}

// Exists check if the user exists by counting the results found
func (u *User) Exists() (bool, *errors.Error) {
	Logger.Debug("[User.Exists]: Checking if user exists")

	if u.ID != 0 {
		Logger.WithField("ID", u.ID).Debug("[User.Exists]: found ID")
		return u.ExistsWithCond(db.Cond{"id": u.ID})

	} else if u.Username != "" {
		Logger.WithField("Username", u.Username).Debug("[User.Exists]: found Username")
		return u.ExistsWithCond(db.Cond{"username": u.Username})

	} else if u.Email != "" {
		Logger.WithField("Email", u.Email).Debug("[User.Exists]: found Email")
		return u.ExistsWithCond(db.Cond{"email": u.Email})
	}

	Logger.WithField("User", u).Warn("[User.Exists]: Not enough info to check if user exists")
	return false, errors.FromCode(errors.ErrorNotEnoughInfo)
}

// ExistsWithCond check if the user exists by using the given condition and counting the results found
func (u *User) ExistsWithCond(cond db.Cond) (bool, *errors.Error) {
	if cond == nil {
		return false, errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	var count uint64
	var err error

	// find and count using a condition
	Logger.WithField("cond", cond).Debug("[User.ExistsWithCond]: finding...")
	if count, err = uc.Find(cond).Count(); err != nil {

		Logger.WithError(err).Error("[User.ExistsWithCond]: error")
		return false, errors.FromErr(err)
	}

	// count of users found
	exists := count > 0

	Logger.WithField("found", exists).Debug("[User.ExistsWithCond]: got response")
	return exists, nil
}

// Create validates and check if a user exists before inserting into the db
func (u *User) Create() (hide.Int64, *errors.Error) {
	if u.Username == "" || u.Email == "" || u.Password == "" {
		Logger.WithField("User", u).Error("[User.Create]: Not enough info!")
		return 0, errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	Logger.Debug("[User.Create] Creating user...")
	valid, err := u.Validate()
	if err != nil {
		return 0, errors.Mask(err, errors.ErrorUserInvalid)
	}

	// check if valid
	Logger.WithField("valid", valid).Debug("[User.Create]: User validated")
	if !valid {
		return 0, errors.FromCode(errors.ErrorUserInvalid)
	}

	// data to check if exists
	var exists bool
	info := map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
	}

	// check if exists
	for field, value := range info {
		if exists, err = u.ExistsWithCond(db.Cond{field: value}); err != nil {
			Logger.WithError(err).Error("[User.Create]: error while checking if the user exists")
			return 0, err
		}

		if exists {
			switch field {
			case "username":
				Logger.Debug("[User.Create]: username exists")
				return 0, errors.FromCode(errors.ErrorUsernameExists)

			case "email":
				Logger.Debug("[User.Create]: email exists")
				return 0, errors.FromCode(errors.ErrorEmailExists)
			}
		}
	}

	// hash the password
	err = u.Hash()
	if err != nil {
		Logger.WithError(err).Error("[User.Create]: error while hashing password")
		return 0, err
	}

	// default values
	Logger.WithField("username", u.Username).Debug("[User.Create]: Setting default values for user")
	u.Created = time.Now()

	// insert into the database
	Logger.WithField("username", u.Username).Debug("[User.Create]: Inserting user into the database")
	id, gErr := uc.Insert(u)
	if gErr != nil {
		Logger.WithError(err).Error("[User.Create]: error while inserting the user into the database")
		return 0, err
	}

	// insert returned id into User
	u.ID = hide.Int64(id.(int64))

	// gives error when id is 0 (nil)
	if u.ID == 0 {
		Logger.Debug("[User.Create]: id returned 0, trying User.Find")
		found, err := u.Find()
		if err != nil {
			return 0, err
		}

		// gives error when not found
		if !found {
			return 0, errors.FromCode(errors.ErrorUserDoesntExists)
		}
	}

	Logger.WithFields(log.Fields{
		"id":       u.ID,
		"username": u.Username,
	}).Debug("[User.Create]: User inserted into the database")

	return u.ID, nil
}

//

// Hash the user's password
func (u *User) Hash() *errors.Error {
	// check for empty password
	if len(u.Password) == 0 {
		Logger.Error("[User.Encrypt]: invalid password")
		return errors.FromCode(errors.ErrorUserInvalidPassword)
	}

	// encrypt the password using bcrypt
	p, err := bcrypt.GenerateFromPassword([]byte(u.Password), Config.EncryptionLevel)
	if err != nil {
		Logger.WithError(err).Error("[User.Encrypt]: error during encryption")
		return errors.FromErr(err)
	}

	// writes the encrypTed password into the user struct
	u.Password = string(p)
	Logger.Debug("[User.Encrypt]: password encrypted")
	return nil
}

// Save the user's data into the database
// aka update
func (u *User) Save() *errors.Error {
	Logger.Debug("[User.Save]: Finding user info...")
	cond := db.Cond{}

	if u.ID != 0 {
		cond["id"] = u.ID

	} else if u.Username != "" {
		cond["username"] = u.Username

	} else if u.Email != "" {
		cond["email"] = u.Email

	} else {
		Logger.Warn("[User.Save]: No user info found!")
		return errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	Logger.Debug("[User.Save]: Found user info...")
	return u.SaveWithCond(cond)
}

// SaveWithCond updates the user's data on the db with conditions
func (u *User) SaveWithCond(cond db.Cond) *errors.Error {
	Logger.WithField("cond", cond).Debug("[User.SaveWithCond]: Saving user...")
	err := uc.Find(cond).Update(u)

	if err != nil {
		Logger.WithError(err).Error("[User.SaveWithCond]: Error while saving the user")
		return errors.FromErr(err)
	}

	Logger.Debug("[User.SaveWithCond]: User saved!")
	return nil
}

// IsBanned checks if a ban expired
// then removes the ban state and save to the database
func (u *User) IsBanned() (bool, *errors.Error) {
	l := Logger.WithField("User", u.ID)
	l.Debug("[User.IsBanned]: Checking user banned state...")

	var err error

	// when a ban is not cached at "u.banned"
	// it checks the database for it
	if u.Banned == nil {
		l.Debug("[User.IsBanned]: ban is not cached, checking database...")

		// find and count ban records
		r := bc.Find(db.Cond{"user_id": u.ID, "state": true})
		c, err := r.Count()
		if err != nil {
			l.WithError(err).Error("[User.IsBanned]: Error while counting ban records")
			return false, errors.FromErr(err)
		}

		// no ban records were found
		if c == 0 {
			l.Debug("[User.IsBanned]: No ban records found, User is not banned")
			return false, nil
		}

		err = r.One(&u.Banned)
		if err != nil {
			l.WithError(err).Error("[User.IsBanned]: Error while finding bans")
			return false, errors.FromErr(err)
		}
	}

	// user is not banned
	if !u.Banned.State {
		l.Debug("[User.IsBanned]: User is not banned")
		return false, nil
	}

	// check if ban expired to disable it
	if u.Banned.Until != nilTime && u.Banned.Until.After(time.Now()) {
		l.Debug("[User.IsBanned]: User is banned")
		return true, nil
	}

	l.Debug("[User.IsBanned]: Ban expired, updating database...")
	u.Banned.State = false

	// update the ban state
	err = bc.Find(db.Cond{"id": u.Banned.ID}).Update(u.Banned)
	if err != nil {
		l.WithError(err).Error("[User.IsBanned]: Error while finding bans")
		return false, errors.FromErr(err)
	}

	l.Debug("[User.IsBanned]: Ban state updated")
	return false, nil
}

// Ban a user and save the state to the database
// it can be temporary or permanent
func (u *User) Ban(temporary bool, until time.Time) *errors.Error {
	if u.ID == 0 {
		return errors.FromCode(errors.ErrorUserInvalid)
	}

	b := &Ban{
		UserID: u.ID,

		// apply the ban
		State:     true,
		Temporary: temporary,
		Starts:    time.Now(),
	}

	// only insert the until time when it is temporary
	if temporary {
		b.Until = until
	}

	// insert the ban into the database
	_, err := bc.Insert(b)
	if err != nil {
		return errors.FromErr(err)
	}

	return nil
}

//

// HardDelete removes the user from the database
// use SoftDelete to disable a account
func (u *User) HardDelete() *errors.Error {
	if u.ID == 0 {
		return errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	var err *errors.Error
	cond := db.Cond{"id": u.ID}

	// deletes rows from a table/collection based on a condition
	del := func(c db.Collection, cond db.Cond) *errors.Error {
		if err != nil {
			return err
		}

		// find the rows
		r := c.Find(cond)
		if r.Err() != nil {
			return errors.FromErr(r.Err())
		}

		// delete from the database
		err := r.Delete()
		return errors.FromErr(err)
	}

	err = del(uc, cond) // user account

	cond = db.Cond{"user_id": u.ID}
	err = del(bc, cond) // user bans
	err = del(ac, cond) // user activation
	err = del(ec, cond) // user events

	return err
}

// SoftDelete disables the account
// without removing it from the database
func (u *User) SoftDelete() *errors.Error {
	u.Deleted = true
	u.Activated = false
	return u.Save()
}

// Find the user using the data on the struct
// id -> username -> email -> fail
func (u *User) Find() (bool, *errors.Error) {
	Logger.Debug("[User.Find]: Finding user...")

	// validate before finding
	_, err := u.Validate()
	if err != nil {
		return false, err
	}

	Logger.Debug("[User.Find]: Checking for user info")

	// check if theres some data
	// id, username or email
	if u.ID > 0 {
		Logger.Debug("[User.Find]: found ID, trying FindWithCond(id)")
		return u.FindWithCond(db.Cond{"id": u.ID})

	} else if u.Username != "" {
		Logger.Debug("[User.Find]: found Username, trying FindWithCond(username)")
		return u.FindWithCond(db.Cond{"username": u.Username})

	} else if u.Email != "" {
		Logger.Debug("[User.Find]: found Email, trying FindWithCond(email)")
		return u.FindWithCond(db.Cond{"email": u.Email})
	}

	Logger.Warn("[User.Find]: No info found")

	// theres no info so it can find
	return false, errors.FromCode(errors.ErrorNotEnoughInfo)
}

// FindWithCond tries to find the user using the give conditions
func (u *User) FindWithCond(cond db.Cond) (bool, *errors.Error) {
	l := Logger.WithField("cond", cond)
	l.Debug("[User.FindWithCond]: Finding User...")
	user := NewUser()

	// find user on db and insert into the struct
	err := uc.Find(cond).One(&user)
	if err != nil {
		l.WithError(err).Error()
		return false, errors.FromErr(err)
	}

	// checks if ban expired when one is found
	_, cErr := user.IsBanned()
	if cErr != nil {
		l.WithError(cErr).Error("")
		return false, cErr
	}

	found := user != nil

	// replace the current user with
	// the one from the database
	u.Replace(user)
	Logger.WithFields(log.Fields{
		"id":    u.ID,
		"found": found,
	}).Debug("[User.Find]: found user")
	return found, nil
}

// Replace a user mem address
func (u *User) Replace(user *User) {
	Logger.WithField("from", user).Debug("[User.Replace]: Replacing User")
	*u = *user
}

// Validate the user's struct
// only fields that aren't empty will be validaTed!
func (u *User) Validate() (bool, *errors.Error) {
	Logger.WithField("User", u).Debug("[User.Validate]: Validating user")

	result, err := govalidator.ValidateStruct(u)
	if err != nil {
		Logger.WithError(err).Error("[User.Validate]: error")
		return false, errors.FromErr(err)
	}

	Logger.WithField("result", result).Debug("[User.Validate]: User validated")
	return result, nil
}

// ComparePassword checks if the given password is the same as the one
// in the database after hashing it, a u.Find() is required before using it
// when error is nil, means the passwords are equal
func (u *User) ComparePassword(password string) *errors.Error {
	l := Logger.WithFields(log.Fields{
		"ID":       u.ID,
		"Username": u.Username,
		"Email":    u.Email,
	})

	l.Debug("[User.ComparePassword]: Comparing passwords...")

	// check for empty passwords
	if u.Password == "" || password == "" {
		return errors.FromCode(errors.ErrorNoPasswordToCompare)
	}

	// compare the password
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		l.Debug("[User.ComparePassword]: Passwords are equal")
	} else {
		l.Debug("[User.ComparePassword]: Wrong password")
	}

	return errors.FromErr(err)
}

// Auth authenticates a user and return a jwt token
func (u *User) Auth(password string) (string, *errors.Error) {
	l := Logger.WithFields(log.Fields{
		"ID":       u.ID,
		"Username": u.Username,
		"Email":    u.Email,
	})
	l.Debug("[User.Auth]: Authenticating user...")

	// check for empty fields
	if password == "" {
		l.Error("[User.Auth]: The given password is empty")
		return "", errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	if u.Username == "" && u.Email == "" {
		l.Error("[User.Auth]: No Username and Email were specified")
		return "", errors.FromCode(errors.ErrorNotEnoughInfo)
	}

	// validate the username or email
	valid, err := u.Validate()
	if err != nil {
		l.Debug("[User.Auth]: User isn't valid")
		return "", err
	}

	if !valid {
		l.Debug("[User.Auth]: User isn't valid")
		return "", errors.FromCode(errors.ErrorUserInvalid)
	}

	// tries to find the user with the username or email
	found, err := u.Find()
	if err != nil {
		l.Debug("[User.Auth]: Error while finding user")
		return "", err
	}

	if !found {
		l.Debug("[User.Auth]: User not found")
		return "", errors.FromCode(errors.ErrorUserDoesntExists)
	}

	// check if the password given is equal
	err = u.ComparePassword(password)
	if err != nil {
		return "", err
	}

	// obfuscate the user id
	obfuscatedID, gErr := u.ID.MarshalJSON()
	if gErr != nil {
		l.WithError(gErr).Error("[User.Auth]: Can't obfuscate user id")
		return "", errors.FromErr(gErr)
	}

	// encrypt the obfuscaTed user id
	id, err := util.Encrypt(string(obfuscatedID), Config.EncryptionKey)
	if err != nil {
		return "", err
	}

	// encrypt the user's power
	power, err := util.Encrypt(strconv.Itoa(u.Power), Config.EncryptionKey)
	if err != nil {
		return "", err
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":    id,
		"power": power,
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(Config.TokenExpirationTime).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret
	tokenString, gErr := token.SignedString(Config.TokenSecret)
	if gErr != nil {
		l.WithError(gErr).Error("[User.Auth]: Can't sign jwt token")
		return "", errors.FromErr(gErr)
	}

	l.Debug("[User.Auth]: Token created for user")
	return tokenString, nil
}

// SetIDFromString parses a user id from a string and insert it
// into the user
func (u *User) SetIDFromString(id string) *errors.Error {
	// convert the id to int64
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.FromErr(err)
	}

	// insert the id into the user
	u.ID = hide.Int64(i)
	return nil
}

// JUser prevents a loop in User.MarshalJSON()
type JUser User

// MarshalJSON hides the user password before transforming it into a json
func (u *User) MarshalJSON() ([]byte, error) {
	pwd := u.Password
	u.Password = ""

	defer func() {
		u.Password = pwd
	}()

	return json.Marshal(JUser(*u))
}
