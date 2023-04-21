module github.com/linkdata/jawsdemo

go 1.16

require (
	github.com/linkdata/deadlock v0.4.0
	github.com/linkdata/jaws v0.20.0
)

// For debugging local copy of JaWS. Also set jaws version above to v0.0.0
// replace github.com/linkdata/jaws => ../jaws

require github.com/klauspost/compress v1.16.5 // indirect
