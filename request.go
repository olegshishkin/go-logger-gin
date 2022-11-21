package main

import (
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
)

const (
	noBody = "<no body>"
)

type rqTemplate struct {
	correlation *uuid.UUID
	method      string
	size        int64
	remoteAddr  string
	clientIP    string
	path        string
	headers     http.Header
	body        string
}

func (t *rqTemplate) shortString() string {
	return fmt.Sprintf("%v %v %s %s", t.correlation, request, t.method, t.path)
}

func (t *rqTemplate) String() string {
	return fmt.Sprintf("%s %d bytes from %s ip %s, headers: %v", t.shortString(), t.size, t.remoteAddr, t.clientIP, t.headers)
}

func (t *rqTemplate) fullString() string {
	if t.size > 0 {
		return fmt.Sprintf("%s, body: %s", t.String(), t.body)
	}
	return fmt.Sprintf("%s, %s", t.String(), noBody)
}
