package main

import (
	"fmt"
	"os"

	"github.com/crossplane/crossplane-runtime/pkg/password"
)

func main() {
	pwsettings := password.Settings{
		CharacterSet: `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()<>?{}[]-_=+\/|`,
		Length:       64,
	}

	password, ok := pwsettings.Generate()

	if ok != nil {
		fmt.Println("error during password creation")
		os.Exit(1)
	}

	fmt.Println(password)
}
