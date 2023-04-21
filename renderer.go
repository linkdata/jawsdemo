package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/linkdata/jaws"
)

type renderer struct {
	jw *jaws.Jaws
	g  *Globals
	t  *template.Template
}

func (t *renderer) makeHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// jaws.NewRequest() creates a new jaws.Request that tracks
		// elements registered during the template rendering, and
		// prepares JaWS to accept the incoming WebSocket call for
		// this request.
		rq := t.jw.NewRequest(context.Background(), r)

		t.g.RLock()
		defer t.g.RUnlock()
		if err := t.t.ExecuteTemplate(w, name, NewUiState(rq, t.g)); err != nil {
			log.Fatal(err)
		}
	}
}
