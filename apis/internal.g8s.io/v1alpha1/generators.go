package v1alpha1

import (
	gopgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/the-gizmo-dojo/g8s/pkg/qanda"
)

const (
	pwCharacterSet = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()<>?{}[]-_=+\/|`
	pwLength       = 64
)

func NewPasswordGenerator(pgft PasswordGeneratorFuncType, pwopts PasswordOptions) Generator[PasswordGeneratorFuncType] {
	return Generator[PasswordGeneratorFuncType]{
		GeneratorFunc: NewPasswordGeneratorFunc(pwopts),
		Options:       pwopts,
	}
}

type PasswordGeneratorFuncType func(pwopts PasswordOptions) (string, error)

func NewPasswordGeneratorFunc(Options) PasswordGeneratorFuncType {
	pgft := PasswordGeneratorFuncType(
		func(pwopts PasswordOptions) (string, error) {
			settings := password.Settings{
				CharacterSet: pwCharacterSet,
				Length:       int(pwopts.Length),
			}
			password, err := settings.Generate()
			return password, err
		},
	)

	return pgft
}

func NewPGPKeyGenerator(pgpkgft PGPKeyGeneratorFuncType, pgpkopts PGPKeyOptions) Generator[PGPKeyGeneratorFuncType] {
	return Generator[PGPKeyGeneratorFuncType]{
		GeneratorFunc: NewPGPKeyGeneratorFunc(pgpkopts),
		Options:       pgpkopts,
	}
}

type PGPKeyGeneratorFuncType func(pgpopts PGPKeyOptions) (gopgpcrypto.Key, error)

func NewPGPKeyGeneratorFunc(Options) PGPKeyGeneratorFuncType {
	pgpkgft := PGPKeyGeneratorFuncType(
		func(pgpkopts PGPKeyOptions) (gopgpcrypto.Key, error) {
			q, _ := qanda.NewQuestion(qanda.PGPKeyPair, qanda.Daimon)
			pgpkey, _ := q.Ask()
			return pgpkey.Content.(gopgpcrypto.Key), nil
		},
	)

	return pgpkgft
}
