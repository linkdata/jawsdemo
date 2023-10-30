module github.com/linkdata/jawsdemo

go 1.21

require (
	github.com/linkdata/deadlock v0.4.0
	github.com/linkdata/jaws v0.40.0
)

require (
	github.com/klauspost/compress v1.15.5 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

// For debugging local copy of JaWS. Also set jaws version above to v0.0.0
// replace github.com/linkdata/jaws => ../jaws
