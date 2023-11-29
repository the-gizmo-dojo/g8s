package v1alpha1

import (
	"github.com/the-gizmo-dojo/g8s/pkg/qanda"
)

type Name string

type Generator[A qanda.Answerer] struct {
	qanda.Answerer
	Validator
}

func Generate(g Generator[qanda.Answerer]) error {
	_, ok := g.Respond()

	return ok
}

type Validator interface {
	Validate() error
}

type NamespaceOptions struct {
	Name `json:"name,omitempty"`
}

func (nso NamespaceOptions) Validate() error {
	// TODO for now just getting the basics working
	return nil
}

type PasswordOptions struct {
	Name   `json:"name,omitempty"`
	Length uint8 `json:"length,omitempty"`
}

func (po PasswordOptions) Validate() error {
	// TODO for now just getting the basics working
	return nil
}

type PGPKeyOptions struct {
	Name       `json:"name,omitempty"`
	EntityName string `json:"entityName,omitempty"`
	Address    string `json:"address,omitempty"`
	Email      string `json:"email,omitempty"`
}

func (pgpo PGPKeyOptions) Validate() error {
	// TODO for now just getting the basics working
	return nil
}

type SSHKeyOptions struct {
	Name `json:"name,omitempty"`
}

func (ssho SSHKeyOptions) Validate() error {
	// TODO for now just getting the basics working
	return nil
}
