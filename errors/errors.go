package errors

// we'll not use "<< iota" because
// it may break things if i add a new error
const (
	ErrorUserInvalidUsername ID = 1 // invalid username - when a username is not in [a-zA-Z0-9_]
	ErrorUserInvalidEmail    ID = 2 // invalid email
	ErrorUserInvalidPassword ID = 3 // invalid password - when a password is bigger than 255

	ErrorUserUsernameExists ID = 4
	ErrorUserEmailExists    ID = 5

	ErrorUserNotFound     ID = 6 // user doesn't exists
	ErrorUserNotSpecified ID = 7 // not specified - when a field doesn't have a value when looking for a user

	ErrorUserBanned             ID = 8  // banned
	ErrorUserNoPower            ID = 9  // no power - when a user doesn't have enough power
	ErrorUserTokenInvalidMethod ID = 10 // invalid token method - when a JWT header is invalid

	ErrorUserUnauthorized ID = 97 // Unauthorized
	ErrorUserInvalidInfo  ID = 98 // Invalid info - when a user information is invalid when provided from a url like (GET: api/user/ufo@usa.gov/exists)
	ErrorUserUnknown      ID = 99 // I Don't know what the error is, and at this point, i'm too afraid to ask
)
