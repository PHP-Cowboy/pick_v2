package ecode

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

var (
	_codes = make(map[int]string)
)

//var _codes map[int]string //register ecode

type Codes interface {
	// sometimes Error return Code in string form
	Error() string
	// Code get error ecode.
	Code() int
	// Message get ecode message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
}

// New new a ecode.Codes by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(e int, msg string) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = msg
	return Code(e)
}

// A Code is an int error ecode spec.
type Code int

func (c Code) Error() string {
	return strconv.Itoa(int(c))
}

// Code return error ecode
func (c Code) Code() int {
	return int(c)
}

// Message return error message
func (c Code) Message() string {
	msg, ok := _codes[c.Code()]
	if ok {
		return msg
	}
	return c.Error()
}

// Details return details.
func (c Code) Details() []interface{} {
	return nil
}

// Int parse ecode int to error.
func Int(i int) Code {
	return Code(i)
}

// String parse ecode string to error.
// 如果是数字字符串通过转换找到错误对象
func String(e string) Code {
	if e == "" {
		return Success
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Code(i)
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return Success
	}
	if eCode, ok := errors.Cause(e).(Codes); ok {
		return eCode
	}
	return String(e.Error())
}

// 比较两个Codes 是否相等
func Equal(c1, c2 Codes) bool {
	if c1 == nil {
		c1 = Success
	}
	if c2 == nil {
		c2 = Success
	}
	return c1.Code() == c2.Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return code.Code() == Cause(err).Code()
}
