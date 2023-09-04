// this package defines various agents who can be queried for responses or to take action
package agent

type Daimon struct {
	Name         string
	Representing Human
}

func NewDaimon(n string, h Human) Daimon {
	// must prompt?

	return Daimon{
		Name:         n,
		Representing: h,
	}
}
