package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type requestBodyLogger struct {
	logWriter io.Writer
	handler   http.Handler
}

func requestBodyLoggerMiddleware(out io.Writer, h http.Handler) http.Handler {
	return &requestBodyLogger{out, h}
}

func (l *requestBodyLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	closer := ioutil.NopCloser(bytes.NewBuffer(buf))
	req.Body = closer
	buf = append(buf, '\n')

	l.handler.ServeHTTP(w, req)

	inputPrefix := "> "
	var buffer bytes.Buffer

	buffer.Write([]byte(inputPrefix))
	err = json.Indent(&buffer, buf, inputPrefix, "  ")
	if err != nil {
		l.logWriter.Write(append([]byte(inputPrefix), buf...))
	} else {
		buffer.WriteTo(l.logWriter)
	}
}

type responseBodyLogger struct {
	responseWriter http.ResponseWriter
	body           []byte
	handler        http.Handler
	logWriter      io.Writer
}

func responseBodyLoggerMiddleware(out io.Writer, handler http.Handler) http.Handler {
	return &responseBodyLogger{nil, make([]byte, 0), handler, out}
}

func (l *responseBodyLogger) Write(bytes []byte) (int, error) {
	l.body = append(l.body, bytes...)
	l.body = append(l.body, '\n')
	return l.responseWriter.Write(bytes)
}

func (l *responseBodyLogger) WriteHeader(header int) {
	l.responseWriter.WriteHeader(header)
}

func (l *responseBodyLogger) Header() http.Header {
	return l.responseWriter.Header()
}

func (l *responseBodyLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.responseWriter = w
	l.handler.ServeHTTP(l, req)

	outputPrefix := "< "
	var buffer bytes.Buffer

	buffer.Write([]byte(outputPrefix))
	err := json.Indent(&buffer, l.body, outputPrefix, "  ")
	if err != nil {
		l.logWriter.Write(append([]byte(outputPrefix), l.body...))
	} else {
		buffer.WriteTo(l.logWriter)
	}
	l.body = make([]byte, 0)
}

type requestHeaderLogger struct {
	logWriter io.Writer
	handler   http.Handler
}

func requestHeaderLoggerMiddleware(out io.Writer, handler http.Handler) http.Handler {
	return &requestHeaderLogger{out, handler}
}

func (l *requestHeaderLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.handler.ServeHTTP(w, req)
	headersPrefix := "  "
	for k, v := range req.Header {
		l.logWriter.Write([]byte(fmt.Sprintf("%s[%s: %s]\n", headersPrefix, k, strings.Join(v, " "))))
	}
}

type requestURILogger struct {
	logWriter      io.Writer
	handler        http.Handler
	statusCode     int
	responseWriter http.ResponseWriter
}

func (l *requestURILogger) Header() http.Header {
	return l.responseWriter.Header()
}

func (l *requestURILogger) Write(buf []byte) (int, error) {
	return l.responseWriter.Write(buf)
}

func (l *requestURILogger) WriteHeader(status int) {
	l.statusCode = status
	l.responseWriter.WriteHeader(status)
}

func requestURILoggerMiddleware(out io.Writer, handler http.Handler) http.Handler {
	return &requestURILogger{out, handler, 200, nil}
}

func (l *requestURILogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	l.responseWriter = w
	l.handler.ServeHTTP(l, req)

	str := fmt.Sprintf("%s: [%s, \"%s\", %s] = (%d=%s)\n",
		time.Now().String(),
		req.RemoteAddr, req.RequestURI, req.Method,
		l.statusCode, http.StatusText(l.statusCode))

	l.logWriter.Write([]byte(str))
}
