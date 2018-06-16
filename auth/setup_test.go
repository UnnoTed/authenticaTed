package users

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	isTest = true

	err := Setup()
	if err != nil {
		panic(err)
	}

	err = Connect()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = Disconnect()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}
