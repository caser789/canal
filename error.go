package mysql

import (
	"errors"
	"fmt"
)

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")
	ErrTxDone        = errors.New("sql: Transaction has already been committed or rolled back")
)

type MySQLError struct {
	Code    uint16
	Message string
	State   string
}

func (e *MySQLError) String() string {
	return fmt.Sprintf("ERROR %d (%s): %s", e.Code, e.State, e.Message)
}

func NewDefaultError(errCode uint16, args ...interface{}) *MySQLError {
	e := new(MySQLError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DEFAULT_MYSQL_STATE
	}

	if format, ok := MySQLErrName[errCode]; ok {
		e.Message = fmt.Sprintf(format, args...)
	} else {
		e.Message = fmt.Sprint(args...)
	}

	return e
}

func NewError(errCode uint16, message string) *MySQLError {
	e := new(MySQLError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DEFAULT_MYSQL_STATE
	}

	e.Message = message
	return e
}
