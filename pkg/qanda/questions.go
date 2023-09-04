// a package to ask questions and either get a generated response or prompt for input

/*
how do i want to consume this pkg?

q := qanda.NewQuestion(qanda.Password, qanda.Human)
var answer = $question.Ask
*/
package qanda

import (
	"fmt"
)

type Question struct {
	Subject
	Audience
}

type Subject string

const (
	Password Subject = "Password"
	SSHKey   Subject = "SSHKey"
	TLSKey   Subject = "TLSKey"
	GPGKey   Subject = "GPGKey"
)

func validSubject(s Subject) bool {
	sVals := []Subject{Password, SSHKey, TLSKey, GPGKey}

	for _, v := range sVals {
		if v == s {
			return true
		}
	}

	return false
}

func (s Subject) String() string {
	return string(s)
}

type Audience string

const (
	Human   Audience = "Human"   // bufio.NewReader(os.Stdin), some REST form
	Machine Audience = "Machine" // automatically generated
)

func validAudience(a Audience) bool {
	aVals := []Audience{Human, Machine}

	for _, v := range aVals {
		if v == a {
			return true
		}
	}

	return false
}

func (a Audience) String() string {
	return string(a)
}

// Question constructor
func NewQuestion(s Subject, a Audience) (q Question, err error) {
	sok := validSubject(s)
	aok := validAudience(a)

	if !sok {
		err = fmt.Errorf("Error creating Question: invalid Subject '%s'", s.String())
		return q, err
	} else if !aok {
		err = fmt.Errorf("Error creating Question: invalid Audience '%s'", a.String())
		return q, err
	}

	q = Question{
		Subject:  s,
		Audience: a,
	}

	return q, err
}

// automagically determines the name of the type to call Respond() on
func (q Question) Ask() (a Answer, err error) {
	typeName := q.Audience.String() + q.Subject.String()
	inst, ok := answererMap[typeName]

	if !ok {
		err = fmt.Errorf("Error during Ask(), cannot find correct type.")
		return a, err
	}

	return inst.Respond()
}
