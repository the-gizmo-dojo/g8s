package human

const (
	KeyType = "x25519"
)

type Human struct {
	Name       string
	Email      string
	Passphrase []byte
	KeyType    string
}

func NewHuman(n, e string, p []byte) Human {
	return Human{
		Name:       n,
		Email:      e,
		Passphrase: p,
		KeyType:    KeyType,
	}
}

// func GetName() Name

// func GetEmail() Email

// func GetPassphrase() Passphrase
