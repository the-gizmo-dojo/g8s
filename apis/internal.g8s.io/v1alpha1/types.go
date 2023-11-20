package v1alpha1

type Source string

const (
	Generated Source = "generated"
	Input     Source = "input"
)

type Generator struct {
	Spec GeneratorSpec
}

type GeneratorSpec struct {
	GenerateFunc func(any) any
	Options      []Option
}

type Option struct {
}

type PasswordOptions struct {
}

type PasswordGenerator struct {
}
