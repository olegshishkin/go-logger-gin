package gin

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
	params      string
	headers     http.Header
	body        string
}

func (t *rqTemplate) shortString() string {
	return fmt.Sprintf("%v %v %s %s", t.correlation, request, t.method, t.path)
}

func (t *rqTemplate) String() string {
	params := ""
	if t.params != "" {
		params = fmt.Sprintf(", params: %s", t.params)
	}

	args := []any{
		t.shortString(),
		t.size,
		t.remoteAddr,
		t.clientIP,
		params,
		t.headers,
	}

	return fmt.Sprintf("%s, %d bytes from %s ip %s%s, headers: %v", args...)
}

func (t *rqTemplate) fullString() string {
	if t.size > 0 {
		return fmt.Sprintf("%s, body: %s", t.String(), t.body)
	}

	return fmt.Sprintf("%s, %s", t.String(), noBody)
}
