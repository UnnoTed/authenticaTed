package users

// UserPower is the level of power that a user can have access to
type UserPower int

// This is a list of user's power levels
const (
	// Normal powers
	// UserPowerNone is the user that hasn't activaTed his account yet
	UserPowerNone UserPower = iota
	// UserPowerNormal is the user that activaTed his account
	UserPowerNormal
	// UserPowerPremium is the user that paid/donaTed
	UserPowerPremium

	// Limited powers
	// UserPowerMod has the powers to ban and warn users
	UserPowerMod
	// UserPowerBot has the power to read private information (email) but can not modify it
	UserPowerBot

	// All powers
	// UserPowerAdmin has the powers to make mods and edit users' information
	UserPowerAdmin
	// UserPowerOwner can make admins
	UserPowerOwner
	// UserPowerProgrammer can do everything and has access to db info
	UserPowerProgrammer
)
