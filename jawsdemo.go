package main

import (
	"context"
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

type Template struct {
	jw        *jaws.Jaws
	g         *Globals
	templates *template.Template
}

func (t *Template) renderer(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rq := t.jw.NewRequest(context.Background(), r)
		t.g.RLock()
		defer t.g.RUnlock()
		if err := t.templates.ExecuteTemplate(w, name, NewUiState(rq, t.g)); err != nil {
			log.Fatal(err)
		}
	}
}

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

	jw := jaws.New()
	jw.Logger = log.Default()
	jawsboot.Setup(jw)
	defer jw.Close()
	go jw.Serve()

	g := NewGlobals()
	go g.ClockFn(jw)

	t := &Template{
		jw:        jw,
		g:         g,
		templates: template.Must(template.New("").Funcs(jaws.FuncMap).ParseGlob("assets/*.html")),
	}
	http.DefaultServeMux.Handle("/jaws/", jw)
	http.DefaultServeMux.HandleFunc("/", t.renderer("index.html"))
	http.DefaultServeMux.HandleFunc("/cars", t.renderer("cars.html"))

	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(breakChan)
		log.Print("listening on \"http://localhost:8081/\"")
		log.Print(http.ListenAndServe("localhost:8081", nil))
	}()

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
