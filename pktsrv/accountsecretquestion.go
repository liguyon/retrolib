package pktsrv

import (
	"errors"
	"strings"
)

// AccountSecretQuestion represents a packet containing the question asked when deleting a character on the account.
type AccountSecretQuestion struct {
	// Question is the question that must be answered before deleting a character
	Question string
}

func (p *AccountSecretQuestion) TypeID() string {
	return "AQ"
}

func (p *AccountSecretQuestion) Marshal() ([]byte, error) {
	if len(p.Question) == 0 {
		return nil, errors.New("empty question")
	}
	return []byte(p.TypeID() + p.Question), nil
}

func (p *AccountSecretQuestion) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), p.TypeID())
	if len(pl) == 0 {
		return errors.New("insufficient data")
	}
	p.Question = pl
	return nil
}
