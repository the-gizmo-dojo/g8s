package user

const (
	KeyType = "x25519"
)

type Name string

type Email string

type Passphrase string

type User struct {
	Name       Name
	Email      Email
	Passphrase Passphrase
	KeyType    string
}

func NewUser(n Name, e Email, p Passphrase) User {
	return User{
		Name:       n,
		Email:      e,
		Passphrase: p,
		KeyType:    KeyType,
	}
}

// func GetName() Name

// func GetEmail() Email

// func GetPassphrase() Passphrase
