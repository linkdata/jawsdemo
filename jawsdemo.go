package main

import (
	"context"
	"embed"
	"flag"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/jawsboot"
	"github.com/linkdata/jaws/lib/templatereloader"
	"github.com/linkdata/jaws/lib/ui"
	"github.com/linkdata/staticserve"
	"github.com/linkdata/webserv"
)

//go:embed assets
var assetsFS embed.FS

//go:generate go run github.com/linkdata/gitsemver@latest -gopackage -package main -out version.gen.go

var (
	flagAddress    = flag.String("address", os.Getenv("WEBSERV_ADDRESS"), "serve HTTP requests on given [address][:port]")
	flagCertDir    = flag.String("certdir", os.Getenv("WEBSERV_CERTDIR"), "where to find fullchain.pem and privkey.pem")
	flagUser       = flag.String("user", envOrDefault("WEBSERV_USER", ""), "switch to this user after startup (*nix only)")
	flagDataDir    = flag.String("datadir", envOrDefault("WEBSERV_DATADIR", "$HOME"), "where to store data files after startup")
	flagListenURL  = flag.String("listenurl", os.Getenv("WEBSERV_LISTENURL"), "specify the external URL clients can reach us at")
	flagCpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flagMemprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func envOrDefault(envvar, defval string) (s string) {
	if s = os.Getenv(envvar); s == "" {
		s = defval
	}
	return
}

var globals = NewGlobals()

func maybeLogError(err error) {
	if err != nil {
		slog.Error("error", "text", err)
	}
}

func setupRoutes(jw *jaws.Jaws, mux *http.ServeMux) (err error) {
	var tmpl jaws.TemplateLookuper
	if tmpl, err = templatereloader.New(assetsFS, "assets/ui/*.html", ""); err == nil {
		if err = jw.AddTemplateLookuper(tmpl); err == nil {
			err = jw.Setup(mux.Handle, "/static", jawsboot.Setup,
				staticserve.MustNewFS(assetsFS, "assets/static", "images/favicon.png", "mousetrack.js"))
			if err == nil {
				mux.Handle("/jaws/", jw) // ensure the JaWS routes are handled
				mux.Handle("GET /{$}", jw.Session(ui.Handler(jw, "index.html", globals)))
				mux.Handle("GET /cars", jw.Session(ui.Handler(jw, "cars.html", globals)))
			}
		}
	}
	return
}

func backgroundUpdates(jw *jaws.Jaws) {
	started := time.Now()
	time.Sleep(time.Until(started.Truncate(time.Second).Add(time.Second)))
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
			globals.runtime = time.Since(started).String()
			globals.mu.Unlock()
			jw.Dirty(globals.Runtime())
			jw.Dirty(globals.CarsLink())
			jw.Dirty(globals.Client())
		}
	}
}

func main() {
	flag.Parse()

	if *flagCpuprofile != "" {
		f, err := os.Create(*flagCpuprofile)
		if err == nil {
			defer f.Close()
			if err = pprof.StartCPUProfile(f); err == nil {
				defer pprof.StopCPUProfile()
			}
		}
		maybeLogError(err)
	}

	cfg := webserv.Config{
		Address:   *flagAddress,
		CertDir:   *flagCertDir,
		User:      *flagUser,
		DataDir:   *flagDataDir,
		ListenURL: *flagListenURL,
		Logger:    slog.Default(),
	}

	l, err := cfg.Listen()
	if err == nil {
		defer l.Close()

		var jw *jaws.Jaws
		// create a default JaWS instance
		if jw, err = jaws.New(); err == nil {
			defer jw.Close()          // ensure we clean up
			jw.Logger = cfg.Logger    // optionally set the logger to use
			jw.Debug = deadlock.Debug // optionally set the debug flag

			mux := http.NewServeMux() // create a ServeMux to do routing
			if err = setupRoutes(jw, mux); err == nil {
				globals.FaviconURI = jw.FaviconURL()

				// start the JaWS processing loop, Serve returns when Close is called.
				go jw.Serve()

				// spin up a goroutine to update the clock and cars button text
				go backgroundUpdates(jw)

				// start serving requests using the default secure headers and with a JaWS session
				err = cfg.Serve(context.Background(), l, jw.Session(jw.SecureHeadersMiddleware(mux)))
			}
		}
	}

	maybeLogError(err)

	if *flagMemprofile != "" {
		f, err := os.Create(*flagMemprofile)
		if err == nil {
			defer f.Close()
			err = pprof.WriteHeapProfile(f)
		}
		maybeLogError(err)
	}
}
