package config

import (
	"fmt"
	"io"
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
	name := string(cf)

	info, err := os.Stat(name)
	if err == nil && info.IsDir() {
		err = fmt.Errorf("Error: filename is a directory, not a gates.yaml")
		fmt.Println(err)
		os.Exit(1)
	} else if err != nil {
		err = fmt.Errorf("Error: cannot find file")
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(name)
	if err != nil {
		return config, err
	}

	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("Error reading %s file", cf)
		fmt.Println(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling %s file", cf)
		fmt.Println(err)
		os.Exit(1)
	}

	return config, err
}
