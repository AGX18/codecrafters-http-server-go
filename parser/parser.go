package parser

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

var readerPool = sync.Pool{
	New: func() any {
		return bufio.NewReader(nil)
	},
}

func getReader(conn net.Conn) *bufio.Reader {
	r := readerPool.Get().(*bufio.Reader)
	r.Reset(conn)
	return r
}

func putReader(r *bufio.Reader) {
	r.Reset(nil)
	readerPool.Put(r)
}

func parseRequestLine(line string) (string, string, string, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", line)
	}
	return parts[0], parts[1], parts[2], nil
}

func ParseRequest(conn net.Conn) (*Request, error) {
	r := getReader(conn)
	defer putReader(r)

	requestLine, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request line: %w", err)
	}

	requestLine = strings.TrimRight(requestLine, "\r\n")

	method, path, version, err := parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}
	req := &Request{
		Method:  method,
		Path:    path,
		Version: version,
	}

	return req, nil
}
