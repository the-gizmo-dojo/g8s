package main

import (
	"fmt"

	"github.com/the-gizmo-dojo/core-secrets/pkg/qanda"
)

const (
	SSHDIR string = "/home/jrodonnell/.ssh/core-ssh"
)

func main() {

	pq, _ := qanda.NewQuestion(qanda.Password, qanda.Machine)
	pw, _ := pq.Ask()
	fmt.Println(pw.Content)
	/*
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

	   user := human.NewHuman("riley", "riley@gmail.com", []byte("hello world"))

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
	*/
}
