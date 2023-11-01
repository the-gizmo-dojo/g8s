package v1

type ID string

type Password struct {
	ID     `json:"id,omitempty"`
	Length uint8 `json:"length,omitempty"`
}

type PGPKey struct {
	ID      `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Email   string `json:"email,omitempty"`
}

type SSHKey struct {
	ID `json:"id,omitempty"`
}

type TLSCert struct {
	ID       `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}
