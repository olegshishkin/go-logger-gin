package gin

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/olegshishkin/go-logger"
	"io"
	"time"
)

// WebServerLogger provides a Gin handler for logging through the Logger interface.
func WebServerLogger(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch l.GetLevel() {
		case logger.Trace:
			full(c, l)
		case logger.Debug, logger.Info:
			normal(c, l)
		case logger.Warn, logger.Error, logger.Fatal:
			short(c)
		}
	}
}

// full writes verbose logs.
func full(c *gin.Context, l logger.Logger) {
	start := time.Now()

	rq := rqTemp(c, l)
	rq.body = rqBody(c, l)
	l.Trace("%s", rq.fullString())

	wrapRsBody(c)
	c.Next()

	rs := &rsTemplate{
		correlation: rq.correlation,
		status:      c.Writer.Status(),
		size:        c.Writer.Size(),
		body:        c.Writer.(*rsBodyWrapper).buf.String(),
		latency:     time.Since(start),
	}
	l.Trace("%s", rs.fullString())
}

// normal writes lightweight logs.
func normal(c *gin.Context, l logger.Logger) {
	start := time.Now()

	rq := rqTemp(c, l)

	switch l.GetLevel() {
	case logger.Debug:
		l.Debug("%s", rq.String())
	case logger.Info:
		l.Info("%s", rq.shortString())
	}

	c.Next()

	rs := &rsTemplate{
		correlation: rq.correlation,
		status:      c.Writer.Status(),
		size:        c.Writer.Size(),
		latency:     time.Since(start),
	}

	switch l.GetLevel() {
	case logger.Debug:
		l.Debug("%s", rs.String())
	case logger.Info:
		l.Info("%s", rs.String())
	}
}

// short doesn't log anything.
func short(c *gin.Context) {
	c.Next()
}

// rqTemp creates a template for the request logging.
func rqTemp(c *gin.Context, l logger.Logger) *rqTemplate {
	rq := &rqTemplate{
		method:     c.Request.Method,
		size:       c.Request.ContentLength,
		remoteAddr: c.Request.RemoteAddr,
		clientIP:   c.ClientIP(),
		path:       c.Request.URL.Path,
		params:     c.Request.URL.RawQuery,
		headers:    c.Request.Header,
	}

	if u, err := uuid.NewV4(); err != nil {
		l.Error(err, "UUID hasn't been generated for request %s", rq.shortString())
	} else {
		rq.correlation = &u
	}

	return rq
}

// rqBody reads the request body with the ability to read it again.
func rqBody(c *gin.Context, l logger.Logger) string {
	if c.Request.ContentLength <= 0 {
		return noBody
	}

	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		l.Error(err, "Request body hasn't been read for request %s", c.Request.URL.Path)
		return noBody
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(b))

	return string(b)
}

// wrapRsBody adds the ability to read response body multiple times.
func wrapRsBody(c *gin.Context) {
	c.Writer = &rsBodyWrapper{
		buf:            &bytes.Buffer{},
		ResponseWriter: c.Writer,
	}
}
