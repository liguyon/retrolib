package pktsrv

import "github.com/liguyon/retrolib/login"

type AccountSelectServerEncode struct {
	Success bool
	Addr    string
	Port    string
	ErrID   login.SelectServerErrorID
}

func (a *AccountSelectServerEncode) TypeID() string {
	return "AX"
}

func (a *AccountSelectServerEncode) Marshal() ([]byte, error) {

}

func (a *AccountSelectServerEncode) Unmarshal(bytes []byte) error {

}
