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

var errInvalidWriterType = errors.New("invalid writer type")

// WebServerLogger provides a Gin handler for logging through the Logger interface.
func WebServerLogger(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		switch log.GetLevel() {
		case logger.Trace:
			full(ctx, log)
		case logger.Debug, logger.Info:
			normal(ctx, log)
		case logger.Warn, logger.Error, logger.Fatal:
			short(ctx)
		}
	}
}

// full writes verbose logs.
func full(ctx *gin.Context, log logger.Logger) {
	start := time.Now()

	rqTempl := rqTemp(ctx, log)
	rqTempl.body = rqBody(ctx, log)
	log.Trace("%s", rqTempl.fullString())

	wrapRsBody(ctx)
	ctx.Next()

	wrapper, ok := ctx.Writer.(*rsBodyWrapper)
	if !ok {
		log.Error(errInvalidWriterType, "failed type assertion to custom body writer")

		return
	}

	rsTempl := &rsTemplate{
		correlation: rqTempl.correlation,
		status:      ctx.Writer.Status(),
		size:        ctx.Writer.Size(),
		body:        wrapper.buf.String(),
		latency:     time.Since(start),
	}
	log.Trace("%s", rsTempl.fullString())
}

// normal writes lightweight logs.
func normal(ctx *gin.Context, log logger.Logger) {
	start := time.Now()

	rqTempl := rqTemp(ctx, log)

	switch log.GetLevel() {
	case logger.Debug:
		log.Debug("%s", rqTempl.String())
	case logger.Info:
		log.Info("%s", rqTempl.shortString())
	}

	ctx.Next()

	rsTempl := &rsTemplate{
		correlation: rqTempl.correlation,
		status:      ctx.Writer.Status(),
		size:        ctx.Writer.Size(),
		latency:     time.Since(start),
	}

	switch log.GetLevel() {
	case logger.Debug:
		log.Debug("%s", rsTempl.String())
	case logger.Info:
		log.Info("%s", rsTempl.String())
	}
}

// short doesn't log anything.
func short(c *gin.Context) {
	c.Next()
}

// rqTemp creates a template for the request logging.
func rqTemp(ctx *gin.Context, log logger.Logger) *rqTemplate {
	rqTempl := &rqTemplate{
		method:     ctx.Request.Method,
		size:       ctx.Request.ContentLength,
		remoteAddr: ctx.Request.RemoteAddr,
		clientIP:   ctx.ClientIP(),
		path:       ctx.Request.URL.Path,
		params:     ctx.Request.URL.RawQuery,
		headers:    ctx.Request.Header,
	}

	if u, err := uuid.NewV4(); err != nil {
		log.Error(err, "UUID hasn't been generated for request %s", rqTempl.shortString())
	} else {
		rqTempl.correlation = &u
	}

	return rqTempl
}

// rqBody reads the request body with the ability to read it again.
func rqBody(ctx *gin.Context, log logger.Logger) string {
	if ctx.Request.ContentLength <= 0 {
		return noBody
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err, "Request body hasn't been read for request %s", ctx.Request.URL.Path)

		return noBody
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

	return string(body)
}

// wrapRsBody adds the ability to read response body multiple times.
func wrapRsBody(c *gin.Context) {
	c.Writer = &rsBodyWrapper{
		buf:            &bytes.Buffer{},
		ResponseWriter: c.Writer,
	}
}
