package main

import (
	"context"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/jawsboot"
	"github.com/linkdata/jaws/jawsecho"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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
	jawsboot.Setup(jw)
	jw.Logger = log.Default()
	go jw.Serve()
	defer jw.Close()

	g := NewGlobals()
	go g.ClockFn(jw)

	e := echo.New()
	e.Renderer = &Template{templates: template.Must(template.New("").Funcs(jaws.FuncMap).ParseGlob("assets/*.html"))}
	jawsecho.Setup(e, jw)
	e.GET("/", func(c echo.Context) (err error) {
		rq := jw.NewRequest(context.Background(), c.Request())
		if cookie := rq.EnsureSession(1, 60); cookie != nil {
			c.SetCookie(cookie)
			log.Println("new session", cookie.Value)
		}
		g.RLock()
		defer g.RUnlock()
		if err = c.Render(http.StatusOK, "index.html", NewUiState(rq, g)); err != nil {
			log.Println(c.Request().RequestURI, err)
		}
		return
	})
	e.GET("/cars", func(c echo.Context) (err error) {
		rq := jw.NewRequest(context.Background(), c.Request())
		cookie := rq.SessionCookie()
		log.Println("cars session", cookie.Value)
		g.RLock()
		defer g.RUnlock()
		if err = c.Render(http.StatusOK, "cars.html", NewUiState(rq, g)); err != nil {
			log.Println(c.Request().RequestURI, err)
		}
		return
	})

	breakChan := make(chan os.Signal, 1)
	signal.Notify(breakChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(breakChan)
		log.Print(e.Start(":8081"))
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
