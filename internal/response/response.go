package response

import (
	"fmt"
	"strconv"
	"time"

	"github.com/VladanT3/TCP_to_HTTP/internal/request/headers"
)

type ResponseWriter struct {
	Data []byte
}

func (w *ResponseWriter) WriteResponse(status_code int, headers headers.Headers, body string) {
	w.WriteStatusLine(status_code)

	default_headers := GetDefaultHeaders()
	for key, val := range default_headers {
		_, ok := headers[key]
		if !ok {
			headers[key] = val
		}
	}

	headers["Content-Length"] = strconv.Itoa(len(body))

	w.WriteHeaders(headers)

	w.WriteBody([]byte(body))
}

func (w *ResponseWriter) WriteStatusLine(status_code int) {
	switch status_code {
	case 100:
		w.Write([]byte("HTTP/1.1 100 Continue\r\n"))
	case 101:
		w.Write([]byte("HTTP/1.1 101 Switching Protocols\r\n"))
	case 103:
		w.Write([]byte("HTTP/1.1 103 Early Hints\r\n"))
	case 200:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case 201:
		w.Write([]byte("HTTP/1.1 201 Created\r\n"))
	case 202:
		w.Write([]byte("HTTP/1.1 202 Accepted\r\n"))
	case 203:
		w.Write([]byte("HTTP/1.1 203 Non-Authoritative Information\r\n"))
	case 204:
		w.Write([]byte("HTTP/1.1 204 No Content\r\n"))
	case 205:
		w.Write([]byte("HTTP/1.1 205 Reset Content\r\n"))
	case 206:
		w.Write([]byte("HTTP/1.1 206 Partial Content\r\n"))
	case 207:
		w.Write([]byte("HTTP/1.1 207 Multi Status\r\n"))
	case 208:
		w.Write([]byte("HTTP/1.1 208 Already Reported\r\n"))
	case 226:
		w.Write([]byte("HTTP/1.1 226 IM Used\r\n"))
	case 300:
		w.Write([]byte("HTTP/1.1 300 Multiple Choices\r\n"))
	case 301:
		w.Write([]byte("HTTP/1.1 301 Moved Permanently\r\n"))
	case 302:
		w.Write([]byte("HTTP/1.1 302 Found\r\n"))
	case 303:
		w.Write([]byte("HTTP/1.1 303 See Other\r\n"))
	case 304:
		w.Write([]byte("HTTP/1.1 304 Not Modified\r\n"))
	case 307:
		w.Write([]byte("HTTP/1.1 307 Temporary Redirect\r\n"))
	case 308:
		w.Write([]byte("HTTP/1.1 308 Permanent Redirect\r\n"))
	case 400:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case 401:
		w.Write([]byte("HTTP/1.1 401 Unauthorized\r\n"))
	case 402:
		w.Write([]byte("HTTP/1.1 402 Payment Required\r\n"))
	case 403:
		w.Write([]byte("HTTP/1.1 403 Forbidden\r\n"))
	case 404:
		w.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
	case 405:
		w.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n"))
	case 406:
		w.Write([]byte("HTTP/1.1 406 Not Acceptable\r\n"))
	case 407:
		w.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\n"))
	case 408:
		w.Write([]byte("HTTP/1.1 408 Request Timeout\r\n"))
	case 409:
		w.Write([]byte("HTTP/1.1 409 Conflict\r\n"))
	case 410:
		w.Write([]byte("HTTP/1.1 410 Gone\r\n"))
	case 411:
		w.Write([]byte("HTTP/1.1 411 Length Required\r\n"))
	case 412:
		w.Write([]byte("HTTP/1.1 412 Precondition Failed\r\n"))
	case 413:
		w.Write([]byte("HTTP/1.1 413 Content Too Large\r\n"))
	case 414:
		w.Write([]byte("HTTP/1.1 414 URI Too Long\r\n"))
	case 415:
		w.Write([]byte("HTTP/1.1 415 Unsupported Media Type\r\n"))
	case 416:
		w.Write([]byte("HTTP/1.1 416 Range Not Satisfiable\r\n"))
	case 417:
		w.Write([]byte("HTTP/1.1 417 Expectation Failed\r\n"))
	case 418:
		w.Write([]byte("HTTP/1.1 418 I'm a teapot\r\n"))
	case 421:
		w.Write([]byte("HTTP/1.1 421 Misdirected Request\r\n"))
	case 422:
		w.Write([]byte("HTTP/1.1 422 Unprocessable Content\r\n"))
	case 423:
		w.Write([]byte("HTTP/1.1 423 Locked\r\n"))
	case 424:
		w.Write([]byte("HTTP/1.1 424 Failed Dependency\r\n"))
	case 426:
		w.Write([]byte("HTTP/1.1 426 Upgrade Required\r\n"))
	case 428:
		w.Write([]byte("HTTP/1.1 428 Precondition Required\r\n"))
	case 429:
		w.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n"))
	case 431:
		w.Write([]byte("HTTP/1.1 431 Request Header Fields Too Large\r\n"))
	case 451:
		w.Write([]byte("HTTP/1.1 451 Unavailable For Legal Reasons\r\n"))
	case 500:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	case 501:
		w.Write([]byte("HTTP/1.1 501 Not Implemented\r\n"))
	case 502:
		w.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n"))
	case 503:
		w.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n"))
	case 504:
		w.Write([]byte("HTTP/1.1 504 Gateway Timeout\r\n"))
	case 505:
		w.Write([]byte("HTTP/1.1 505 HTTP Version Not Supported\r\n"))
	case 506:
		w.Write([]byte("HTTP/1.1 506 Variant Also Negotiates\r\n"))
	case 507:
		w.Write([]byte("HTTP/1.1 507 Insufficient Storage\r\n"))
	case 508:
		w.Write([]byte("HTTP/1.1 508 Loop Detected\r\n"))
	case 510:
		w.Write([]byte("HTTP/1.1 510 Not Extended\r\n"))
	case 511:
		w.Write([]byte("HTTP/1.1 511 Network Authentication Required\r\n"))
	default:
		w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(status_code)) + " \r\n"))
	}
}

func GetDefaultHeaders() headers.Headers {
	header := make(headers.Headers)
	header["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	header["Connection"] = "close"
	header["Content-Type"] = "text/html"

	return header
}

func (w *ResponseWriter) WriteHeaders(headers headers.Headers) {
	data := ""
	for key, val := range headers {
		data += fmt.Sprintf("%s:%s\r\n", key, val)
	}
	data += "\r\n"

	w.Write([]byte(data))
}

func (w *ResponseWriter) WriteBody(data []byte) {
	w.Write(data)
}

func (w *ResponseWriter) WriteChunkedBody(p []byte) int {
	data := fmt.Sprintf("%X\r\n%s\r\n", len(p), p)
	w.Write([]byte(data))
	return len(data)
}

func (w *ResponseWriter) WriteChunkedBodyDone() int {
	w.Write([]byte("0\r\n\r\n"))
	return 5
}

func (w *ResponseWriter) WriteTrailers(headers headers.Headers) {
	data := ""
	for key, val := range headers {
		data += fmt.Sprintf("%s:%s\r\n", key, val)
	}
	data += "\r\n"

	w.Write([]byte(data))
}

func (w *ResponseWriter) Write(data []byte) {
	w.Data = append(w.Data, data...)
}
