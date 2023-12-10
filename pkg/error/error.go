package error

import (
	"errors"
	"fmt"
)

var (
	ErrBadResponse   = errors.New("api responded with bad status code")
	ErrBadEnvLoading = errors.New("can not load env variable")
)

func CertainError(prefix string, err error) error {
	return errors.New(fmt.Sprintf("%s: %v", prefix, err))
}
