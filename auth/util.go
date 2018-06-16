package users

import (
	"errors"
	"strconv"
	"time"

	"github.com/UnnoTed/authenticaTed/util"

	"github.com/c2h5oh/hide"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	issuer = "auth.service"
)

// UserToken holds a user id inside a jwt
type UserToken struct {
	UID   string `json:"id,string"`
	Power string `json:"power"`
	jwt.StandardClaims
}

// CreateToken creates a jwt token with a ID
func CreateToken(hID interface{}, power string, encrypt bool) (string, error) {
	var id string

	switch hID.(type) {
	case hide.Int64:
		id = strconv.FormatInt(int64(hID.(hide.Int64)), 10)
	case string:
		id = hID.(string)
	case int64:
		id = strconv.FormatInt(hID.(int64), 10)
	}

	var (
		encryptedID = id
		err         error
	)

	if encrypt {
		// encrypt the user's ID
		encryptedID, err = util.Encrypt(id, Config.EncryptionKey)
		if err != nil {
			return "", err
		}
	}

	// create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&UserToken{
			encryptedID,
			power,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(Config.TokenExpirationTime).Unix(),
				Id:        encryptedID,
				IssuedAt:  time.Now().Unix(),
				Issuer:    issuer,
				NotBefore: time.Now().Unix(),
			},
		})

	// sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(Config.TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

// WillTokenExpire checks if a token will expire
// the current range is 5-30 minutes
func WillTokenExpire(expAt int64) bool {
	exp := time.Unix(expAt, 0)
	almostNow := time.Now().Add(5 * time.Minute)
	later := time.Now().Add(30 * time.Minute)

	// after 5min and before 30min
	return exp.After(almostNow) && exp.Before(later)
}

// GetPower gets the user's power from the jwt and decrypts it
func GetPower(c echo.Context) (UserPower, error) {
	usr := c.Get(middleware.DefaultJWTConfig.ContextKey)
	if usr == nil {
		return 0, errors.New("There is no token")
	}

	switch usr.(type) {
	case *jwt.Token:
		// decrypt the ID from the token
		p, err := util.Decrypt(usr.(*jwt.Token).Claims.(*UserToken).Power, Config.EncryptionKey)
		if err != nil {
			return 0, err
		}

		power, derr := strconv.Atoi(p)
		if derr != nil {
			return 0, derr
		}

		return UserPower(power), nil
	}

	return 0, errors.New("Not able to find token")
}

// GetID gets the user's ID from the jwt and decrypts it
func GetID(c echo.Context) (hide.Int64, error) {
	usr := c.Get(middleware.DefaultJWTConfig.ContextKey)
	if usr == nil {
		return 0, errors.New("There is no token")
	}

	switch usr.(type) {
	case *jwt.Token:
		// decrypt the ID from the token
		id, err := util.Decrypt(usr.(*jwt.Token).Claims.(*UserToken).UID, Config.EncryptionKey)
		if err != nil {
			return 0, err
		}

		hID, derr := strconv.ParseInt(id, 10, 64)
		if derr != nil {
			return 0, derr
		}

		return hide.Int64(hID), nil
	}

	return 0, errors.New("Not able to find token")
}

// GetUserID decrypts the user id from the user token stored in echo's context
func GetUserID(c echo.Context) (int64, error) {
	usr := c.Get(middleware.DefaultJWTConfig.ContextKey).(*jwt.Token).Claims

	m := structs.Map(usr)

	id, err := util.Decrypt(m["UID"].(string), Config.EncryptionKey)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(id, 10, 64)
}

// ByCreated sorts users by the time that it was inserted into the database
type ByCreated []*User

func (t ByCreated) Len() int           { return len(t) }
func (t ByCreated) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByCreated) Less(i, j int) bool { return t[i].Created.Before(t[j].Created) }
