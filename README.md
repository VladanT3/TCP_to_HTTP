# TCP to HTTP

Making HTTP/1.1 from scratch in Go (kind of, go's 'net' library is still used).

## About

cmd/tcplistener - used before httpserver was made it was the initial thing used to debug request parsing.

cmd/udpsender - silly little thing, just made to see the differences between tcp and udp.

cmd/httpserver - the http server, might rewrite it a bit because it feels a bit cramped and there are some stuff you might not need there necesarilly.

internal/ - request parsing, response writing and server logic.

### To Do (maybe at some point)

- move handler from httpserver to its own package, maybe get rid of it completely but its good for testing still.
- rewrite the response logic a bit, not a fan of it atm.
- maybe broaden it some more so it could actually be used at least for some minimal web servers.
- http/2???
