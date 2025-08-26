package main

import (
	"context"
	"embed"
	"flag"
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
	"github.com/linkdata/jaws/staticserve"
	"github.com/linkdata/jaws/templatereloader"
)

//go:embed assets
var assetsFS embed.FS

var listenaddr = flag.String("listenaddr", "localhost:8081", "address to listen on")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var globals = NewGlobals()

func maybeLogError(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func setupRoutes(jw *jaws.Jaws, mux *http.ServeMux) (faviconuri string, err error) {
	var tmpl jaws.TemplateLookuper
	if tmpl, err = templatereloader.New(assetsFS, "assets/ui/*.html", ""); err == nil {
		jw.AddTemplateLookuper(tmpl)
		err = jw.Setup(mux.Handle, "/static", jawsboot.Setup,
			staticserve.MustNewFS(assetsFS, "assets/static", "images/favicon.png"))
		if err == nil {
			mux.Handle("/jaws/", jw) // ensure the JaWS routes are handled
			mux.Handle("/", jw.Session(jw.Handler("index.html", globals)))
			mux.Handle("/cars", jw.Session(jw.Handler("cars.html", globals)))
		}
		faviconuri = jw.FaviconURL()
	}
	return
}

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

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

	mux := http.NewServeMux() // create a ServeMux to do routing
	jw, err := jaws.New()     // create a default JaWS instance
	if err != nil {
		panic(err)
	}
	defer jw.Close()           // ensure we clean up
	jw.Logger = slog.Default() // optionally set the logger to use
	jw.Debug = deadlock.Debug  // optionally set the debug flag

	faviconuri, err := setupRoutes(jw, mux)
	maybeLogError(err)
	globals.FaviconURI = faviconuri

	go jw.Serve() // start the JaWS processing loop, Serve returns when Close is called.

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
				case 1:
					globals.carsLink = "Did you know VIN numbers are encoded?"
				case 2:
					globals.carsLink = "DO NOT CLICK HERE!"
				case 3:
					globals.carsLink = "Cars"
				default:
					globals.carsLink = "This is a boring link to car info."
				}
				globals.runtime = time.Since(now).String()
				globals.mu.Unlock()
				jw.Dirty(globals.Runtime())
				jw.Dirty(globals.CarsLink())
				jw.Dirty(globals.Client())
			}
		}
	}()

	go func() {
		slog.Info("listening", "address", "http://"+*listenaddr)
		slog.Error(http.ListenAndServe(*listenaddr, mux).Error()) //#nosec G114
	}()

	// wait for stop
	<-ctx.Done()
	slog.Info("stopped", "reason", ctx.Err())

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err == nil {
			defer f.Close()
			err = pprof.WriteHeapProfile(f)
		}
		maybeLogError(err)
	}
}
