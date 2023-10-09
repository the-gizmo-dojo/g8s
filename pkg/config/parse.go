package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	gatesv1 "github.com/the-gizmo-dojo/g8/pkg/apis/gates/v1"
	yaml "gopkg.in/yaml.v3"
)

type ConfigFile string

const (
	DefaultFile    ConfigFile = "gates.yaml"
	passwordsField string     = "passwords"
	pgpKeysField   string     = "pgpKeys"
	sshKeysField   string     = "sshKeys"
	tlsCertsField  string     = "tlsCerts"
)

type Config struct {
	gatesv1.Passwords
	gatesv1.PGPKeys
	gatesv1.SSHKeys
	gatesv1.TLSCerts
}

func Parse(cf ConfigFile) (Config, error) {
	yml := make(map[string]interface{})
	contents, err := os.ReadFile(string(cf))
	if err != nil {
		err = fmt.Errorf("Error reading %s file", cf)
	}

	err = yaml.Unmarshal(contents, yml)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling %s file", cf)
	}

	var config Config
	var passwords gatesv1.Passwords
	var pgpkeys gatesv1.PGPKeys
	var sshkeys gatesv1.SSHKeys
	var tlscerts gatesv1.TLSCerts
	for k, v := range yml {
		switch k {
		case passwordsField:
			for _, m := range v.([]interface{}) {
				var password gatesv1.Password
				err = mapstructure.Decode(m, &password)
				passwords.Passwords = append(passwords.Passwords, password)
			}
		case pgpKeysField:
			for _, m := range v.([]interface{}) {
				var pgpkey gatesv1.PGPKey
				err = mapstructure.Decode(m, &pgpkey)
				pgpkeys.PGPKeys = append(pgpkeys.PGPKeys, pgpkey)
			}
		case sshKeysField:
			for _, m := range v.([]interface{}) {
				var sshkey gatesv1.SSHKey
				err = mapstructure.Decode(m, &sshkey)
				sshkeys.SSHKeys = append(sshkeys.SSHKeys, sshkey)
			}
		case tlsCertsField:
			for _, m := range v.([]interface{}) {
				var tlscert gatesv1.TLSCert
				err = mapstructure.Decode(m, &tlscert)
				tlscerts.TLSCerts = append(tlscerts.TLSCerts, tlscert)
			}
		default:
			err = fmt.Errorf("Unknown field in %s file!", cf)
			fmt.Println(err)
			os.Exit(1)
		}
	}

	config = Config{
		Passwords: passwords,
		PGPKeys:   pgpkeys,
		SSHKeys:   sshkeys,
		TLSCerts:  tlscerts,
	}

	return config, err
}
