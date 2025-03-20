package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/jacobshu/http/internal/logger"
)

var (
	ErrInvalidRequestLineNumber  = errors.New("invalid parameter number in request line\n")
	ErrInvalidRequestLineMethod  = errors.New("invalid method in request line\n")
	ErrInvalidRequestLineVersion = errors.New("invalid HTTP version in request line\n")
	ErrInvalidRequestLineOrder   = errors.New("invalid order of request line parameters\n")
)

const crlf = "\r\n"

type ParserState int

const (
	initialized ParserState = iota
	done
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

	bytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("error reading request:", err)
	}

	requestLine, err := parseRequestLine(bytes)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	log := logger.SetupLogger("development")
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request line")
	}

	requestLineStr := string(data[:idx])
	log.Info(requestLineStr)

	requestLine, err := requestLineFromString(requestLineStr)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")

	if len(parts) != 3 {
		return nil, ErrInvalidRequestLineNumber
	}

	rl := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.Split(parts[2], "/")[1],
	}

	fmt.Println(rl)

	if !isValidMethod(rl.Method) {
		return nil, ErrInvalidRequestLineMethod
	}

	if !isValidVersion(rl.HttpVersion) {
		return nil, ErrInvalidRequestLineVersion
	}

	if isValidMethod(rl.HttpVersion) || isValidMethod(rl.RequestTarget) || isValidVersion(rl.RequestTarget) || isValidVersion(rl.Method) {
		return nil, ErrInvalidRequestLineOrder
	}

	return &rl, nil

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
