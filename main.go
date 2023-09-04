package main

import (
	"fmt"
	"reflect"

	"github.com/the-gizmo-dojo/core-secrets/pkg/qanda"
)

func main() {
	q, _ := qanda.NewQuestion(qanda.PGPKeyPair, qanda.Daimon)
	a, _ := q.Ask()
	fmt.Println(a.Content, reflect.TypeOf(a.Content))
}
