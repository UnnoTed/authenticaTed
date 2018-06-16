package users

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/UnnoTed/authenticaTed/errors"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go-dry"
	"upper.io/db.v2"
)

var user = NewUser()

func TestMarshallJSON(t *testing.T) {
	u := NewUser()
	u.Password = "yes"
	data, err := u.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// json -> interface{}
	var fu map[string]interface{}
	err = json.Unmarshal(data, &fu)
	assert.NoError(t, err)

	// checks if password is present
	p, ok := fu["password"]
	assert.False(t, ok)
	assert.Empty(t, p)

	// checks if User still have the password
	assert.Equal(t, "yes", u.Password)
}

func TestExistsFalse(t *testing.T) {
	u := NewUser()

	// expect error
	exists, err := u.Exists()
	assert.NotNil(t, err)
	assert.False(t, exists)

	u.Username = "Flava_Hustle"

	// expect not found
	exists, err = u.Exists()
	assert.Nil(t, err)
	assert.False(t, exists)

	u.Username = ""
	u.Email = "Flava_Hustle@rapster.com"

	// expect not found
	exists, err = u.Exists()
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestValidate(t *testing.T) {
	u := NewUser()

	// expect error
	u.Username = "Not$valid&Username"
	valid, err := u.Validate()

	assert.NotNil(t, err)
	assert.False(t, valid)

	// expect error
	u.Username = ""
	u.Email = "ayy@.lmao"
	valid, err = u.Validate()

	assert.NotNil(t, err)
	assert.False(t, valid)

	// valid
	u.Email = "ufo@usa.com"
	valid, err = u.Validate()

	assert.Nil(t, err)
	assert.True(t, valid)

	// valid
	u.Username = "JetFuelcantMelt5733LBeans"
	valid, err = u.Validate()

	assert.Nil(t, err)
	assert.True(t, valid)
}

func TestCreate(t *testing.T) {
	user = NewUser()
	user.Password = "password"

	// expect error
	id, err := user.Create()
	assert.NotNil(t, err)
	assert.Zero(t, id)

	// expect error: invalid
	user.Password = ""
	id, err = user.Create()

	assert.NotNil(t, err)
	assert.Equal(t, errors.ErrorNotEnoughInfo, err.Code)
	assert.Zero(t, id)

	// ok
	user.Username = "Flava_Hustle"
	user.Email = "Flava_Hustle@mail.com"
	user.Password = "password"
	id, err = user.Create()

	assert.Nil(t, err)
	assert.NotZero(t, id)
}

func TestExistsTrue(t *testing.T) {
	// ok
	exists, err := user.Exists()
	assert.Nil(t, err)
	assert.True(t, exists)

	// ok
	exists, err = user.ExistsWithCond(db.Cond{"id": user.ID})
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestComparePassword(t *testing.T) {
	err := user.ComparePassword("password")
	assert.Nil(t, err)

	err = user.ComparePassword("err")
	assert.NotNil(t, err)
}

func TestAuth(t *testing.T) {
	// expect error
	token, err := user.Auth("err")
	assert.NotNil(t, err)
	assert.Empty(t, token)

	// expect error
	fu := *user
	fu.Username = ""
	fu.Email = ""

	token, err = fu.Auth("err")
	assert.NotNil(t, err)
	assert.Empty(t, token)

	// ok
	fu.Username = user.Username
	fu.Email = ""

	token, err = fu.Auth("password")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// ok
	fu.Username = ""
	fu.Email = user.Email

	token, err = fu.Auth("password")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// ok
	token, err = user.Auth("password")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestBan(t *testing.T) {
	err := user.Ban(true, time.Now().Add(30*24*time.Hour))
	assert.Nil(t, err)

	// check if banned
	found, err := user.Find()
	assert.Nil(t, err)
	assert.True(t, found)

	assert.NotNil(t, user.Banned)
	assert.True(t, user.Banned.State)

	// check if ban is true on db
	found, err = user.Find()
	assert.Nil(t, err)
	assert.True(t, found)
	assert.True(t, user.Banned.State)
}

func TestDelete(t *testing.T) {
	// ok
	err := user.SoftDelete()
	assert.Nil(t, err)

	// check if still exists and have [Deleted = true]
	exists, err := user.Exists()
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.True(t, user.Deleted)

	// ok
	err = user.HardDelete()
	assert.Nil(t, err)

	// check if still exists
	exists, err = user.Exists()
	assert.Nil(t, err)
	assert.False(t, exists)

	cond := db.Cond{"user_id": user.ID}

	// activations
	count, fErr := ac.Find(cond).Count()
	assert.NoError(t, fErr)
	assert.Zero(t, count)

	// bans
	count, fErr = bc.Find(cond).Count()
	assert.NoError(t, fErr)
	assert.Zero(t, count)

	// events
	count, fErr = ec.Find(cond).Count()
	assert.NoError(t, fErr)
	assert.Zero(t, count)
}

var (
	n string
	c bool
)

func BenchmarkInit(b *testing.B) {
	err := Connect()
	if err != nil {
		b.Error(err)
	}
	b.SkipNow()
}

func BenchmarkCreate(b *testing.B) {
	logger.Level = log.InfoLevel
	Config.EncryptionLevel = 1
	b.ResetTimer()

	for a := 0; a < b.N; a++ {
		n = dry.RandomHexString(10)
		u := NewUser()
		u.Username = n
		u.Email = n + "@gmail.com"
		u.Password = n
		_, err := u.Create()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	us := NewUser()
	us.Username = n
	us.Email = n + "@gmail.com"
	us.Password = n
	b.ResetTimer()

	for a := 0; a < b.N; a++ {
		_, err := us.Find()
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEnd(b *testing.B) {
	err := Disconnect()
	if err != nil {
		b.Error(err)
	}
	b.SkipNow()
}
