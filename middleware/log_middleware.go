package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type resonseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *resonseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *resonseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}

	o.wroteHeader = true
	o.status = code
}

// just for reference
func CustomLogger(out io.Writer, h http.Handler) http.Handler {
	logger := log.New(out, "", 0)

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		o := &resonseObserver{ResponseWriter: rw}
		h.ServeHTTP(o, r)

		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		logger.Printf("%s - - [%s] %q %d %d %q %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			o.status,
			o.written,
			r.Referer(),
			r.UserAgent())
	})
}
