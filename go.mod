module github.com/linkdata/jawsdemo

go 1.17

require (
	github.com/labstack/echo/v4 v4.10.2
	github.com/linkdata/deadlock v0.4.0
	github.com/linkdata/jaws v0.19.1
)

// For debugging local copy of JaWS. Also set jaws version above to v0.0.0
// replace github.com/linkdata/jaws => ../jaws

require (
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
