package gin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
)

// rsBodyWrapper is used for response body buffering.
type rsBodyWrapper struct {
	buf *bytes.Buffer
	gin.ResponseWriter
}

func (w *rsBodyWrapper) Write(b []byte) (int, error) {
	w.buf.Write(b)
	length, err := w.ResponseWriter.Write(b)

	if err != nil {
		return length, fmt.Errorf("failed in writer: %w", err)
	}

	return length, nil
}

func (w *rsBodyWrapper) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	length, err := w.ResponseWriter.WriteString(s)

	if err != nil {
		return length, fmt.Errorf("failed in string writer: %w", err)
	}

	return length, nil
}
