// a package to ask questions and either get a generated response or prompt for input

/*
how do i want to consume this pkg?

q := qanda.NewQuestion(qanda.Password, qanda.Human)
var answer = $question.Ask
*/

// TODO needs to with with context so that information can be passed through
// to various Human and Daimon handlers

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
	Password   Subject = "Password"
	PGPKeyPair Subject = "PGPKeyPair"
	SSHKeyPair Subject = "SSHKeyPair"
	TLSKeyPair Subject = "TLSKeyPair"
)

func validSubject(s Subject) bool {
	sVals := []Subject{Password, SSHKeyPair, TLSKeyPair, PGPKeyPair}

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

/*
Audience is either Daimon, potentially could have Human option one day (some kind of prompt / input).

	Daimon: means input is not needed, can be generated from information the system already knows.
*/
type Audience string

const (
	Daimon Audience = "Daimon" // automatically generated
)

func validAudience(a Audience) bool {
	aVals := []Audience{Daimon}

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
