package g8s

import (
	"github.com/crossplane/crossplane-runtime/pkg/password"
	"github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
)

func GeneratePassword(pw *v1alpha1.Password) string {
	settings := password.Settings{
		Length:       int(pw.Spec.Length),
		CharacterSet: pw.Spec.CharacterSet,
	}

	// error can be ignored because if there's a problem it will be handled in the controller (processNextWorkItem will requeue it)
	pwstr, _ := settings.Generate()

	return pwstr
}
