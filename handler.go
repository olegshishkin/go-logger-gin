package gin

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/olegshishkin/go-logger"
	"io"
	"time"
)

// WebServerLogger provides a Gin handler for logging through the Logger interface.
func WebServerLogger(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		level := l.GetLevel()
		if level == logger.Warn || level == logger.Error || level == logger.Fatal {
			c.Next()
			return
		}
		start := time.Now()

		rq := rq(c, l)
		logRq(c, l, rq)

		if l.GetLevel() == logger.Trace {
			wrapRqBody(c, l, rq)
			wrapRsBody(c)
		}

		c.Next()

		rs := rs(c, rq, start)
		logRs(c, l, rs)
	}
}

func wrapRqBody(c *gin.Context, l logger.Logger, rq *rqTemplate) {
	if rq.size == 0 {
		return
	}
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		l.Error(err, "Request body hasn't been read for request %s", rq.shortString())
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(b))
	rq.body = string(b)
}

func wrapRsBody(c *gin.Context) {
	c.Writer = &rsBodyWrapper{
		buf:            &bytes.Buffer{},
		ResponseWriter: c.Writer,
	}
}

func rq(c *gin.Context, l logger.Logger) *rqTemplate {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	if query != "" {
		path = path + "?" + query
	}

	rq := &rqTemplate{
		method:     c.Request.Method,
		size:       c.Request.ContentLength,
		remoteAddr: c.Request.RemoteAddr,
		clientIP:   c.ClientIP(),
		path:       path,
		headers:    c.Request.Header,
	}

	u, err := uuid.NewV4()
	if err != nil {
		l.Error(err, "UUID hasn't been generated for request %s", rq.shortString())
	} else {
		rq.correlation = &u
	}
	return rq
}

func logRq(c *gin.Context, l logger.Logger, rq *rqTemplate) {
	switch l.GetLevel() {
	case logger.Trace:
		if rq.size > 0 {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				rq.body = noBody
				break
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(b))
			rq.body = string(b)
		}
		l.Trace("%s", rq.fullString())
	case logger.Debug:
		l.Debug("%s", rq.String())
	case logger.Info, logger.Warn:
		l.Info("%s", rq.shortString())
	}
}

func rs(c *gin.Context, rq *rqTemplate, start time.Time) *rsTemplate {
	return &rsTemplate{
		correlation: rq.correlation,
		status:      c.Writer.Status(),
		size:        c.Writer.Size(),
		latency:     time.Since(start),
	}
}

func logRs(c *gin.Context, l logger.Logger, rs *rsTemplate) {
	switch l.GetLevel() {
	case logger.Trace:
		rs.body = c.Writer.(*rsBodyWrapper).buf.String()
		l.Trace("%s", rs.fullString())
	case logger.Debug:
		l.Debug("%s", rs.String())
	case logger.Info, logger.Warn:
		l.Info("%s", rs.String())
	case logger.Error:
		l.Error(errors.New(c.Errors.String()), "%s", rs.String())
	}
}
