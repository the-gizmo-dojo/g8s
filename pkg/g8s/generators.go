package g8s

import (
	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
)

type Gate interface {
	Generate() Answer
	Rotate()
}

type Answer[A any] struct {
	Content A
}

type Password struct {
	v1alpha1.PasswordSpec
}

type PasswordAnswer struct {
	Answer string
}

func (p Password) Generate() Answer {
	content := map[string]int{"asdf": 1}
	return make(Answer[string]{content})
}

func GeneratePassword(pw *v1alpha1.Password) string {
	settings := password.Settings{
		Length:       int(pw.Spec.Length),
		CharacterSet: pw.Spec.CharacterSet,
	}

	// error can be ignored because if there's a problem it will be handled in the controller (processNextWorkItem will requeue it)
	pwstr, _ := settings.Generate()

	return pwstr
}
