// a package to get the answers to whatever question was asked
package qanda

/*
how do i want to consume this package?

q := qanda.New(qanda.Password, qanda.Human)
var answer = $question.Ask
*/

import (
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/password"
)

var answererMap = map[string]Answerer{
	"MachinePassword": MachinePassword{},
	// "MachineSSHKeyResponse": MachineSSHKeyResponse{},
	"HumanPassword": HumanPassword{},
}

type Answer struct {
	Content any
}

type Answerer interface {
	Respond() (Answer, error)
}

type MachinePassword struct {
	password.Settings
}

const (
	pwCharacterSet = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()<>?{}[]-_=+\/|`
	pwLength       = 64
)

func (mp MachinePassword) Respond() (a Answer, err error) {
	mp.Settings = password.Settings{
		CharacterSet: pwCharacterSet,
		Length:       pwLength,
	}

	pw, err := mp.Generate()

	if err != nil {
		err = fmt.Errorf("Error creating password.")
	}

	a.Content = pw
	return a, err
}

type HumanPassword struct {
}

func (hp HumanPassword) Respond() (a Answer, err error) {
	// TODO filler, need to call prompt
	a.Content = "yay"
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
