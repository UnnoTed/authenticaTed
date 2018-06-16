package main

type AuthenticaTedUser struct {
	LastName string `db:"last_name" valid:"optional,length(3|50),alphanum"`
}
