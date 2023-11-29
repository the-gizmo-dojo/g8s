package qanda

import (
	"crypto"
	"fmt"

	gopgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	gopgphelper "github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/charmbracelet/keygen"
	"github.com/crossplane/crossplane-runtime/pkg/password"
)

var answererMap = map[string]Answerer{
	"DaimonPassword":   DaimonPassword{},
	"DaimonPGPKeyPair": DaimonPGPKeyPair{},
	"DaimonSSHKeyPair": DaimonSSHKeyPair{},
}

type Answer struct {
	Content any
}

type Answerer interface {
	Respond() (Answer, error)
}

type PasswordAnswer string

func (p PasswordAnswer) String() string {
	return string(p)
}

type DaimonPassword struct {
	PasswordAnswer
	password.Settings
}

const (
	pwCharacterSet = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()<>?{}[]-_=+\/|`
	pwLength       = 64
)

func (m DaimonPassword) Respond() (a Answer, err error) {
	m.Settings = password.Settings{
		CharacterSet: pwCharacterSet,
		Length:       pwLength,
	}

	pw, err := m.Generate()

	if err != nil {
		err = fmt.Errorf("Error creating password.")
	}

	a.Content = PasswordAnswer(pw)
	return a, err
}

type SSHKeyPairAnswer struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

type DaimonSSHKeyPair struct {
	SSHKeyPairAnswer
	Opts []keygen.Option
}

func (m DaimonSSHKeyPair) Respond() (a Answer, err error) {
	m.Opts = []keygen.Option{(keygen.WithKeyType(keygen.Ed25519))}
	newkey, err := keygen.New("", m.Opts...)

	if err != nil {
		err = fmt.Errorf("Error during SSH key creation.")
	}

	keypair := SSHKeyPairAnswer{
		PrivateKey: newkey.PrivateKey,
		PublicKey:  newkey.CryptoPublicKey,
	}
	a.Content = keypair
	return a, err
}

type PGPKeyPairAnswer struct {
	PGPKey   gopgpcrypto.Key
	Password string
}

type DaimonPGPKeyPair struct {
	PGPKeyPairAnswer
}

const (
	pgpKeyType = "x25519"
)

func (m DaimonPGPKeyPair) Respond() (a Answer, err error) {
	entityname, email := "riley", "riley@gmail.com" // TODO replace

	q, _ := NewQuestion(Password, Daimon)
	pw, _ := q.Ask()
	pwstr := pw.Content.(PasswordAnswer).String()
	pgppw := []byte(pwstr)
	pgpstr, err := gopgphelper.GenerateKey(entityname, email, pgppw, pgpKeyType, 0)

	if err != nil {
		err = fmt.Errorf("Error during PGP key creation.")
	}

	pgpkey, err := gopgpcrypto.NewKeyFromArmored(pgpstr)

	if err != nil {
		err = fmt.Errorf("Error during PGP key creation.")
	}

	keypair := PGPKeyPairAnswer{
		PGPKey:   *pgpkey,
		Password: pwstr,
	}
	a.Content = keypair
	return a, err
}
