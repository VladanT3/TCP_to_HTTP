# TCP to HTTP

Making HTTP/1.1 from scratch in Go (kind of, go's 'net' library is still used).

## About

- cmd/tcplistener - used before httpserver was made it was the initial thing used to debug request parsing.

- cmd/udpsender - silly little thing, just made to see the differences between tcp and udp.

- cmd/httpserver - the http server, might rewrite it a bit because it feels a bit cramped and there are some stuff you might not need there necesarilly.

- internal/ - request parsing, response writing and server logic.

Was planning to maybe make this into a mini http library but I will not be doing that, at least not now. Really just wanted to learn the HTTP protocol.
