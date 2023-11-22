// this package defines various agents who can be queried for responses or to take action
package agent

type Daimon struct {
	Agent
}

func NewDaimon(n, a string) Daimon {
	// must prompt?

	return Daimon{
		Agent: Agent{
			Name:    n,
			Address: a,
		},
	}
}
