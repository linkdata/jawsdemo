module github.com/linkdata/jawsdemo

go 1.19

require (
	github.com/labstack/echo/v4 v4.10.0
	github.com/linkdata/deadlock v0.3.0
	github.com/linkdata/jaws v0.11.0
)

// For debugging local copy of JaWS. Also set jaws version above to v0.0.0
// replace github.com/linkdata/jaws => ../jaws

require (
	github.com/klauspost/compress v1.15.14 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
