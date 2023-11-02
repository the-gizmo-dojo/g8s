package main

import (
	"encoding/json"
	"fmt"

	"github.com/the-gizmo-dojo/g8s/pkg/config"
	// "reflect"
	// "github.com/the-gizmo-dojo/g8/pkg/qanda"
)

func main() {
	content, _ := config.Parse(config.DefaultFile)
	payload, _ := json.Marshal(content)
	fmt.Println(string(payload))

	// q, _ := qanda.NewQuestion(qanda.PGPKeyPair, qanda.Daimon)
	// a, _ := q.Ask()
	// fmt.Println(a.Content, reflect.TypeOf(a.Content))
}
