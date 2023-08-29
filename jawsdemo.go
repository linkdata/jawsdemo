package main

import (
	"flag"
	"html/template"
	"log"
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

	// spin up a goroutine to update the clock and cars button text
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		lastMin := -1
		for range t.C {
			if minute := time.Now().Minute(); minute != lastMin {
				lastMin = minute
				g.ClockString.Store(ClockString())
				jw.Update([]interface{}{g.ClockString})
			}
			if (time.Now().Second() % 3) == 0 {
				switch rand.Intn(5) {
				case 0:
					g.CarsLink.Store("Check out these cars!")
				case 1:
					g.CarsLink.Store("Did you know VIN numbers are encoded?")
				case 2:
					g.CarsLink.Store("DO NOT CLICK HERE!")
				case 3:
					g.CarsLink.Store("Cars")
				default:
					g.CarsLink.Store("This is a boring link to car info.")
				}
				jw.Update([]interface{}{g.CarsLink})
			}
		}
	}()

	// the renderer simplifies making http.HanderFunc functions for us
	t := &renderer{
		jw: jw,
		g:  g,
		t:  template.Must(template.New("").ParseGlob("assets/*.html")),
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
