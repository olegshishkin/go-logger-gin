package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// rsBodyWrapper is used for response body buffering.
type rsBodyWrapper struct {
	buf *bytes.Buffer
	gin.ResponseWriter
}

func (w *rsBodyWrapper) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *rsBodyWrapper) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
