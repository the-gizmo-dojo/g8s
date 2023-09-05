package qanda

import (
	"crypto"
	"fmt"

	gopgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	gopgphelper "github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/charmbracelet/keygen"
	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/the-gizmo-dojo/core-secrets/pkg/agent"
)

var answererMap = map[string]Answerer{
	"DaimonPassword":   DaimonPassword{},
	"DaimonPGPKeyPair": DaimonPGPKeyPair{},
	"DaimonSSHKeyPair": DaimonSSHKeyPair{},
	"HumanPassword":    HumanPassword{},
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
	agent.Daimon
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
	agent.Daimon
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
	agent.Daimon
	PGPKeyPairAnswer
}

const (
	pgpKeyType = "x25519"
)

func (m DaimonPGPKeyPair) Respond() (a Answer, err error) {
	user := agent.NewHuman("riley", "123 fake st", "riley@gmail.com") // TODO replace

	q, _ := NewQuestion(Password, Daimon)
	pw, _ := q.Ask()
	pwstr := pw.Content.(PasswordAnswer).String()
	pgppw := []byte(pwstr)
	pgpstr, err := gopgphelper.GenerateKey(user.Name, user.Email, pgppw, pgpKeyType, 0)

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

type HumanPassword struct {
}

func (hp HumanPassword) Respond() (a Answer, err error) {
	// TODO filler, need to call prompt
	a.Content = PasswordAnswer("yay")
	return a, err
}

// TODO to human pkg?
type Input string

const (
	CLI Input = "CLI"
	GUI Input = "GUI"
)

type Prompter interface {
	Prompt() string
}

func (in Input) Prompt() string {
	// TODO which kind of prompt, cli vs rest form?
	return "prompt"
}
