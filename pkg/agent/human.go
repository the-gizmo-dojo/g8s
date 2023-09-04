// this package defines various agents who can be queried for responses or to take action
package agent

type Human struct {
	Name       string
	Email      string
	Passphrase []byte
}

func NewHuman(n, e string, p []byte) Human {
	// must prompt?
	return Human{
		Name:       n,
		Email:      e,
		Passphrase: p,
	}
}

// func GetName() Name

// func GetEmail() Email

// func GetPassphrase() Passphrase
