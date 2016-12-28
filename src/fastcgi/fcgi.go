package fcgi

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"io"

	"gopkg.in/labstack/echo.v3"
)

type notFound struct {
}

func (notFound) ServeHTTP(http.ResponseWriter, *http.Request) (int, error) {
	return http.StatusNotFound, errors.New("Not forund")
}

type dummyResponseWriter struct {
	buffer  *bytes.Buffer
	status  int
	headers http.Header
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (d *dummyResponseWriter) Header() http.Header {
	return d.headers
}

// Write writes the data to the connection as part of an HTTP reply.
//
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType.
//
// Depending on the HTTP protocol version and the client, calling
// Write or WriteHeader may prevent future reads on the
// Request.Body. For HTTP/1.x requests, handlers should read any
// needed request body data before writing the response. Once the
// headers have been flushed (due to either an explicit Flusher.Flush
// call or writing enough data to trigger a flush), the request body
// may be unavailable. For HTTP/2 requests, the Go HTTP server permits
// handlers to continue to read the request body while concurrently
// writing the response. However, such behavior may not be supported
// by all HTTP/2 clients. Handlers should read before writing if
// possible to maximize compatibility.
func (d *dummyResponseWriter) Write(p []byte) (int, error) {
	return d.buffer.Write(p)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (d *dummyResponseWriter) WriteHeader(s int) {
	d.status = s
}

// NewPHPFastCGIHandler return a simple handler
func NewPHPFastCGIHandler(root, path string, php string, rTimeout, sTimeout, dialTimeout time.Duration) echo.HandlerFunc {
	r := Rule{
		Path: path,

		Address: "tcp://" + php,

		Ext:         ".php",
		SplitPath:   ".php",
		IndexFiles:  []string{"index.php"},
		ReadTimeout: rTimeout,
		SendTimeout: sTimeout,
	}
	d := make([]dialer, 10)
	for i := range d {
		d[i] = basicDialer{
			address: php,
			network: "tcp",
			timeout: dialTimeout,
		}
	}
	r.dialer = &loadBalancingDialer{dialers: d}

	h := Handler{
		Next:    notFound{},
		Rules:   []Rule{r},
		Root:    root,
		AbsRoot: root,
		FileSys: http.Dir(root),
	}

	return func(e echo.Context) error {
		dummy := &dummyResponseWriter{
			status:  0,
			buffer:  &bytes.Buffer{},
			headers: e.Response().Header(),
		}
		s, err := h.ServeHTTP(dummy, e.Request())
		e.Response().Status = s
		_, err = io.Copy(e.Response(), dummy.buffer)
		if err != nil {
			return err
		}

		return err
	}

}
