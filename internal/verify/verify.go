// Package verify is for verifying the credentials of a user are valid.
package verify

import "errors"

var (
	ErrNotFound        = errors.New("could not find entry for given user")
	ErrHashDoesntMatch = errors.New("hashed password is not hash of the given password")
)

type Verifier interface {
	Verify(user string, pw string) (bool, error)
}

type SimpleVerifier struct {
}

func (s *SimpleVerifier) Verify(_ string, _ string) (bool, error) {
	return true, nil
}
