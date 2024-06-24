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

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/jawsboot"
)

var listenaddr = flag.String("listenaddr", "localhost:8081", "address to listen on")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var globals = NewGlobals()

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// parse our templates
	templates := template.Must(template.New("").ParseGlob("assets/*.html"))

	jw := jaws.New()                  // create a default JaWS instance
	jw.AddTemplateLookuper(templates) // let JaWS know about our templates
	defer jw.Close()                  // ensure we clean up
	jw.Logger = slog.Default()        // optionally set the logger to use
	jawsboot.Setup(jw)                // optionally enable the included Bootstrap support
	go jw.Serve()                     // start the JaWS processing loop

	mux := http.NewServeMux() // create a ServeMux to do routing
	mux.Handle("/jaws/", jw)  // ensure the JaWS routes are handled

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
				switch rand.Intn(5) {
				case 0:
					globals.carsLink = "Check out these cars!"
				case 1:
					globals.carsLink = "Did you know VIN numbers are encoded?"
				case 2:
					globals.carsLink = "DO NOT CLICK HERE!"
				case 3:
					globals.carsLink = "Cars"
				default:
					globals.carsLink = "This is a boring link to car info."
				}
				globals.mu.Unlock()
				jw.Dirty(globals.CarsLink())
			}
		}
	}()

	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/", jw.Handler("index.html", globals))
	mux.Handle("/cars", jw.Handler("cars.html", globals))

	// handle CTRL-C and start listening
	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(breakChan)
		slog.Info("listening", "address", "http://"+*listenaddr)
		http.ListenAndServe(*listenaddr, mux)
	}()

	// wait for the HTTP server to stop
	if sig, ok := <-breakChan; ok {
		slog.Info(sig.String())
	}

	if *memprofile != "" {
		if f, err := os.Create(*memprofile); err == nil {
			defer f.Close()
			pprof.WriteHeapProfile(f)
		} else {
			slog.Error(err.Error())
		}
	}
}
