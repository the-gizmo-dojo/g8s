package config

import (
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ConfigFile string

const (
	DefaultFile ConfigFile = "g8s.yaml"
)

type Config struct {
	Passwords []Password `json:"passwords,omitempty" yaml:"passwords,omitempty"`
	PGPKeys   []PGPKey   `json:"pgpKeys,omitempty"   yaml:"pgpKeys,omitempty"`
	SSHKeys   []SSHKey   `json:"sshKeys,omitempty"   yaml:"sshKeys,omitempty"`
	TLSCerts  []TLSCert  `json:"tlsCerts,omitempty"  yaml:"tlsCerts,omitempty"`
}

func Parse(cf ConfigFile) (Config, error) {
	var config Config
	name := string(cf)

	info, err := os.Stat(name)
	if err == nil && info.IsDir() {
		err = fmt.Errorf("Error: filename is a directory, not a g8s.yaml")
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
