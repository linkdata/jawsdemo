package main

import (
	"log"
	"net/http"

	"github.com/linkdata/jaws"
)

type renderer struct {
	jw *jaws.Jaws
	g  *Globals
}

func (rndr *renderer) makeHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// jaws.NewRequest() creates a new jaws.Request that tracks
		// elements registered during the template rendering, and
		// prepares JaWS to accept the incoming WebSocket call for
		// this request.
		rq := rndr.jw.NewRequest(r)

		if err := rndr.jw.Template.ExecuteTemplate(w, name, NewUiState(rq, rndr.g)); err != nil {
			log.Fatal(err)
		}
	}
}
