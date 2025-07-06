package pktsrv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/liguyon/retrolib/login"
)

// AccountLogin represents a packet containing login status information.
// Type ID: Al
// Wire format: Al[Success:K][IsGM] or Al[Success:E][ErrID][Extra]
type AccountLogin struct {
	// Success indicates whether the login process succeeded
	Success bool
	// IsGM indicates whether the account has Game Master access rights
	IsGM bool
	// ErrID is the identifier for the error that indicates why login failed
	ErrID login.LoginErrorID
	// Extra is the extra information sent when the client gets kicked
	Extra string
}

func (p *AccountLogin) TypeID() string {
	return "Al"
}

func (p *AccountLogin) Marshal() ([]byte, error) {
	if p.Success {
		if p.IsGM {
			return []byte(p.TypeID() + "K1"), nil
		} else {
			return []byte(p.TypeID() + "K0"), nil
		}
	} else {
		return []byte(fmt.Sprintf("%sE%c%s", p.TypeID(), p.ErrID, p.Extra)), nil
	}
}

func (p *AccountLogin) Unmarshal(bytes []byte) error {
	pl, _ := strings.CutPrefix(string(bytes), p.TypeID())
	if len(pl) < 2 {
		return errors.New("insufficient data")
	}
	switch pl[0] {
	case 'K':
		p.Success = true
		if pl[1] == '1' {
			p.IsGM = true
		}
	case 'E':
		p.ErrID = login.LoginErrorID(pl[1])
		if len(pl) > 2 {
			p.Extra = pl[2:]
		}
	default:
		return errors.New("invalid data")
	}
	return nil
}
