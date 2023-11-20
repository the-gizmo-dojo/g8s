package config

type Name string

type Password struct {
	Name   `json:"id,omitempty"`
	Length uint8 `json:"length,omitempty"`
}

type PGPKey struct {
	Name       `json:"id,omitempty"`
	EntityName string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Email      string `json:"email,omitempty"`
}

type SSHKey struct {
	Name `json:"id,omitempty"`
}

type TLSCert struct {
	Name     `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}
