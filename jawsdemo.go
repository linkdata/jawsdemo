package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/jawsboot"
)

var listenaddr = flag.String("listenaddr", "localhost:8081", "address to listen on")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	jw := jaws.New()           // create a default JaWS instance
	defer jw.Close()           // ensure we clean up
	jw.CookieName = "jawsdemo" // optionally set a session cookie name
	jw.Logger = log.Default()  // optionally set the logger to use
	jawsboot.Setup(jw)         // optionally enable the included Bootstrap support
	go jw.Serve()              // start the JaWS processing loop

	mux := http.NewServeMux() // create a ServeMux to do routing
	mux.Handle("/jaws/", jw)  // ensure the JaWS routes are handled

	g := NewGlobals() // "Globals" contains the data we want to render
	go g.ClockFn(jw)  // spin up a goroutine to update the clock and cars button text

	// the renderer simplifies making http.HanderFunc functions for us
	t := &renderer{
		jw: jw,
		g:  g,
		t:  template.Must(template.New("").Funcs(jaws.FuncMap).ParseGlob("assets/*.html")),
	}
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("/cars", t.makeHandler("cars.html"))
	mux.HandleFunc("/", t.makeHandler("index.html"))

	// handle CTRL-C and start listening
	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(breakChan)
		log.Printf("listening on %q", *listenaddr)
		log.Print(http.ListenAndServe(*listenaddr, mux))
	}()

	// wait for the HTTP server to stop
	if sig, ok := <-breakChan; ok {
		log.Print(sig)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}
}
