package v1

type ID string

type Password struct {
	ID     `json:"id,omitempty"`
	Length uint8 `json:"length,omitempty"`
}

type Passwords struct {
	Passwords []Password `json:"passwords"`
}

type PGPKey struct {
	ID      `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Email   string `json:"email,omitempty"`
}

type PGPKeys struct {
	PGPKeys []PGPKey `json:"pgpKeys"`
}

type SSHKey struct {
	ID `json:"id,omitempty"`
}

type SSHKeys struct {
	SSHKeys []SSHKey `json:"sshKeys"`
}

type TLSCert struct {
	ID       `json:"id,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

type TLSCerts struct {
	TLSCerts []TLSCert `json:"tlsCerts"`
}
