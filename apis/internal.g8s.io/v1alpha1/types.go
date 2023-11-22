package v1alpha1

type Name string

type Generator[G any] struct {
	GeneratorFunc G
	Options
}

type Options interface {
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
