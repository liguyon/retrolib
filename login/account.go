package login

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// AccountBase represents a Retro game account containing a minimal set of fields. It is sufficient for a Login Server
// implementation, but should be embedded if more features are required.
type AccountBase struct {
	ID                 int
	Username           string
	PasswordHash       []byte
	Nickname           string
	CommunityID        int
	Question           string
	Answer             string
	SubscriptionExpiry time.Time
	BanExpiry          time.Time
}

var passwordRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{6,32}$")

// ValidatePassword validates the given password again the password policy. It returns nil if the password is valid.
// Password policy:
// - between 6 and 32 characters
// - allowed characters: [a-zA-Z0-9_-]
// - does not contain the account's username
func (a *AccountBase) ValidatePassword(pass string) error {
	if matched := passwordRegex.MatchString(pass); !matched {
		return errors.New("invalid")
	}
	if strings.Contains(pass, a.Username) {
		return errors.New("must not contain username")
	}
	return nil
}

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9]{3,16}$")

// ValidateUsername validates the given username against the username policy. It returns nil if the username is valid.
// Username policy:
// - between 3 and 16 characters
// - allowed characters: [a-zA-Z0-9]
func (a *AccountBase) ValidateUsername(user string) error {
	if matched := usernameRegex.MatchString(user); !matched {
		return errors.New("invalid")
	}
	return nil
}

var nicknameRegex = regexp.MustCompile("^[a-zA-Z]{3,16}$")

// ValidateNickname validates the given nickname against the nickname policy. It returns nil if the nickname is valid.
// Nickname policy:
// - between 3 and 16 characters
// - allowed characters: [a-zA-Z]
func (a *AccountBase) ValidateNickname(nick string) error {
	if matched := nicknameRegex.MatchString(nick); !matched {
		return errors.New("invalid")
	}
	return nil
}

// ValidateQuestion validates the given question against the question policy. It returns nil if question is valid.
// Question policy:
// - between 1 and 100 characters
// - allowed characters: printable ASCII chars
func (a *AccountBase) ValidateQuestion(question string) error {
	if len(question) == 0 {
		return errors.New("can't be blank")
	}
	if len(question) > 100 {
		return errors.New("too long")
	}
	if !stringIsPrintableASCII(question) {
		return errors.New("contains invalid characters")
	}
	return nil
}

// ValidateAnswer validates the given answer against the answer policy. It returns nil if answer is valid.
// Answer policy:
// - between 1 and 100 characters
// - allowed characters: printable ASCII chars
func (a *AccountBase) ValidateAnswer(answer string) error {
	if len(answer) == 0 {
		return errors.New("can't be blank")
	}
	if len(answer) > 100 {
		return errors.New("too long")
	}
	if !stringIsPrintableASCII(answer) {
		return errors.New("contains invalid characters")
	}
	return nil
}

func stringIsPrintableASCII(s string) bool {
	for _, v := range s {
		if v < 32 || v > 126 {
			return false
		}
	}
	return true
}
