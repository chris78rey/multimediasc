package oracle

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("oracle: invalid credentials")
	ErrAccountLocked      = errors.New("oracle: account locked")
	ErrPasswordExpired    = errors.New("oracle: password expired")
	ErrNoQueryPermission  = errors.New("oracle: no query permission")
	ErrDBUnavailable      = errors.New("oracle: database unavailable")
)

type OracleError struct {
	Kind error
	Err  error
}

func (e *OracleError) Error() string {
	return e.Err.Error()
}

func (e *OracleError) Unwrap() error {
	return e.Err
}

func ClassifyOracleError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToUpper(err.Error())
	switch {
	case strings.Contains(msg, "ORA-01017"):
		return &OracleError{Kind: ErrInvalidCredentials, Err: err}
	case strings.Contains(msg, "ORA-28000"):
		return &OracleError{Kind: ErrAccountLocked, Err: err}
	case strings.Contains(msg, "ORA-28001"):
		return &OracleError{Kind: ErrPasswordExpired, Err: err}
	case strings.Contains(msg, "ORA-01031"):
		return &OracleError{Kind: ErrNoQueryPermission, Err: err}
	case strings.Contains(msg, "ORA-12170"),
		strings.Contains(msg, "ORA-12514"),
		strings.Contains(msg, "ORA-12541"),
		strings.Contains(msg, "ORA-12545"),
		strings.Contains(msg, "ORA-03113"),
		strings.Contains(msg, "ORA-03114"):
		return &OracleError{Kind: ErrDBUnavailable, Err: err}
	default:
		return err
	}
}

func IsKind(err error, kind error) bool {
	var oe *OracleError
	if !errors.As(err, &oe) {
		return false
	}
	return errors.Is(oe.Kind, kind)
}

func (e *OracleError) String() string {
	return fmt.Sprintf("%v: %v", e.Kind, e.Err)
}
