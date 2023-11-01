package config

import (
	"fmt"
	"os"

	gatesv1 "github.com/the-gizmo-dojo/g8/pkg/apis/gates/v1"
	yaml "gopkg.in/yaml.v3"
)

type ConfigFile string

const (
	DefaultFile ConfigFile = "gates.yaml"
)

type Config struct {
	Passwords []gatesv1.Password `json:"passwords,omitempty" yaml:"passwords,omitempty"`
	PGPKeys   []gatesv1.PGPKey   `json:"pgpKeys,omitempty"   yaml:"pgpKeys,omitempty"`
	SSHKeys   []gatesv1.SSHKey   `json:"sshKeys,omitempty"   yaml:"sshKeys,omitempty"`
	TLSCerts  []gatesv1.TLSCert  `json:"tlsCerts,omitempty"  yaml:"tlsCerts,omitempty"`
}

func Parse(cf ConfigFile) (Config, error) {
	var config Config
	contents, err := os.ReadFile(string(cf))
	if err != nil {
		err = fmt.Errorf("Error reading %s file", cf)
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling %s file", cf)
	}

	return config, err
}
