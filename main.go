package main

import (
	"fmt"
	"os"

	gopgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	gopgphelper "github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/charmbracelet/keygen"
	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/the-gizmo-dojo/core-secrets/pkg/user"
)

const (
	SSHDIR string = "/home/jrodonnell/.ssh/core-ssh"
)

func main() {
	// configure then create password
	// TODO turn into pkg
	// method(num) []Password
	pwsettings := password.Settings{
		CharacterSet: `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()<>?{}[]-_=+\/|`,
		Length:       64,
	}

	password, err := pwsettings.Generate()

	if err != nil {
		fmt.Println("error during password creation")
		os.Exit(1)
	}

	// configure then create ssh key
	opts := []keygen.Option{(keygen.WithKeyType(keygen.Ed25519))}
	keypair, err := keygen.New(SSHDIR, opts...)

	if err != nil {
		fmt.Println("error during ssh key creation")
		os.Exit(1)
	} else if keypair.KeyPairExists() == false {
		keypair.WriteKeys()
	}

	user := user.GetNewUser("riley", "riley@gmail.com", []byte("hello world"))

	gpgstr, err := gopgphelper.GenerateKey(user.Name, user.Email, user.Passphrase, user.KeyType, 0)

	if err != nil {
		fmt.Println("error during gpg key creation")
		os.Exit(1)
	}

	gpgkey, err := gopgpcrypto.NewKeyFromArmored(gpgstr)

	if err != nil {
		fmt.Println("error during gpg key creation")
		os.Exit(1)
	}

	fmt.Println(password, keypair.PrivateKey(), keypair.PublicKey(), gpgkey.GetEntity())
}
