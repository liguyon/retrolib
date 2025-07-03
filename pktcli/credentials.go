package pktcli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Credentials represents a packet containing the credentials to log into an account.
// Type ID: none
// Wire format: [Username]\n#[EncID][PasswordCt]
type Credentials struct {
	// Username is the unique identifier used to log into the account
	Username string
	// EncID is the identifier for the encryption method used for password encryption
	EncID int
	// PasswordCT is the password ciphertext
	PasswordCT string
}

func (c *Credentials) TypeID() string {
	return ""
}

func (c *Credentials) Marshal() ([]byte, error) {
	if len(c.Username) == 0 || len(c.PasswordCT) == 0 {
		return nil, errors.New("insufficient data")
	}
	return []byte(fmt.Sprintf("%s\n#%d%s", c.Username, c.EncID, c.PasswordCT)), nil
}

func (c *Credentials) Unmarshal(bytes []byte) error {
	tks := strings.Split(string(bytes), "\n")
	if len(tks) < 2 || len(tks[1]) < 3 || len(tks[0]) == 0 {
		return errors.New("insufficient data")
	}
	c.Username = tks[0]
	var err error
	c.EncID, err = strconv.Atoi(string(tks[1][1]))
	if err != nil {
		return err
	}
	c.PasswordCT = tks[1][2:]
	return nil
}
