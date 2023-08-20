package user

type User struct {
	Name       string
	Email      string
	Passphrase []byte // TODO make interface so can be retrieved from prompt or passed in programatically
	KeyType    string
}

func GetNewUser(name, email string, pass []byte) User {
	return User{
		Name:       name,
		Email:      email,
		Passphrase: pass,
		KeyType:    "x25519",
	}
}
