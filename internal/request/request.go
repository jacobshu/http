package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

var (
	ErrInvalidRequestLineNumber  = errors.New("invalid parameter number in request line\n")
	ErrInvalidRequestLineMethod  = errors.New("invalid method in request line\n")
	ErrInvalidRequestLineVersion = errors.New("invalid HTTP version in request line\n")
	ErrInvalidRequestLineOrder   = errors.New("invalid order of request line parameters\n")
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("error reading request:", err)
	}

	lines := strings.Split(string(req), "\r\n")
	parts := strings.Split(lines[0], " ")

	if len(parts) != 3 {
		return &Request{}, ErrInvalidRequestLineNumber
	}

	rl := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.Split(parts[2], "/")[1],
	}

	fmt.Println(rl)
	r := Request{RequestLine: rl}

	if !isValidMethod(rl.Method) {
		return &Request{}, ErrInvalidRequestLineMethod
	}

	if !isValidVersion(rl.HttpVersion) {
		return &Request{}, ErrInvalidRequestLineVersion
	}

	if isValidMethod(rl.HttpVersion) || isValidMethod(rl.RequestTarget) || isValidVersion(rl.RequestTarget) || isValidVersion(rl.Method) {
		return &Request{}, ErrInvalidRequestLineOrder
	}

	return &r, nil
}

func isValidVersion(v string) bool {
	return v == "1.1"
}

func isValidMethod(m string) bool {
	methods := []string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
		"PATCH",
	}

	if slices.Contains(methods, m) {
		return true
	}
	return false
}
