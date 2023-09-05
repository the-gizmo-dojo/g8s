// this package defines various agents who can be queried for responses or to take action
package agent

type Human struct {
	Agent
	Email string
}

func NewHuman(n, a, e string) Human {
	// must prompt?
	return Human{
		Agent: Agent{
			Name:    n,
			Address: a,
		},
		Email: e,
	}
}

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

// func GetName() Name

// func GetEmail() Email

// func GetPassphrase() Passphrase
