package main

import (
	"flag"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/jawsboot"
)

var listenaddr = flag.String("listenaddr", "localhost:8081", "address to listen on")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var globals = NewGlobals()

func maybeLogError(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err == nil {
			defer f.Close()
			if err = pprof.StartCPUProfile(f); err == nil {
				defer pprof.StopCPUProfile()
			}
		}
		maybeLogError(err)
	}

	// parse our templates
	templates := template.Must(template.New("").ParseGlob("assets/*.html"))

	mux := http.NewServeMux()                     // create a ServeMux to do routing
	jw := jaws.New()                              // create a default JaWS instance
	jw.AddTemplateLookuper(templates)             // let JaWS know about our templates
	defer jw.Close()                              // ensure we clean up
	jw.Logger = slog.Default()                    // optionally set the logger to use
	jw.Debug = deadlock.Debug                     // optionally set the debug flag
	maybeLogError(jawsboot.Setup(jw, mux.Handle)) // optionally enable the included Bootstrap support
	go jw.Serve()                                 // start the JaWS processing loop
	mux.Handle("/jaws/", jw)                      // ensure the JaWS routes are handled

	// spin up a goroutine to update the clock and cars button text
	go func() {
		now := time.Now()
		time.Sleep(now.Round(time.Second).Sub(now))
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {
			jw.Dirty(uiClock{})
			if (time.Now().Second() % 3) == 0 {
				globals.mu.Lock()
				x := rand.Intn(5) //#nosec G404
				switch x {
				case 0:
					globals.carsLink = "Check out these cars!"
					globals.getUserAgentParam = "runtime: "
				case 1:
					globals.carsLink = "Did you know VIN numbers are encoded?"
					globals.getUserAgentParam = "uptime: "
				case 2:
					globals.carsLink = "DO NOT CLICK HERE!"
					globals.getUserAgentParam = "bored for "
				case 3:
					globals.carsLink = "Cars"
					globals.getUserAgentParam = "waited "
				default:
					globals.carsLink = "This is a boring link to car info."
					globals.getUserAgentParam = "..."
				}
				globals.runtime = time.Since(now).String()
				globals.mu.Unlock()
				jw.Dirty(globals.Runtime())
				jw.Dirty(globals.CarsLink())
				jw.Dirty(globals.GetUserAgent())
				jw.Dirty(globals.Client())
			}
		}
	}()

	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/", jw.Session(jw.Handler("index.html", globals)))
	mux.Handle("/cars", jw.Session(jw.Handler("cars.html", globals)))

	// handle CTRL-C and start listening
	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(breakChan)
		slog.Info("listening", "address", "http://"+*listenaddr)
		slog.Error(http.ListenAndServe(*listenaddr, mux).Error()) //#nosec G114
	}()

	// wait for the HTTP server to stop
	if sig, ok := <-breakChan; ok {
		slog.Info(sig.String())
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err == nil {
			defer f.Close()
			err = pprof.WriteHeapProfile(f)
		}
		maybeLogError(err)
	}
}
