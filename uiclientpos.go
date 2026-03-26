package main

import (
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"

	"github.com/linkdata/jaws"
)

type uiClientPos struct{ *Globals }

func (uic uiClientPos) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	var sessions []*jaws.Session
	uic.mu.RLock()
	for sess := range uic.client {
		sessions = append(sessions, sess)
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID() < sessions[j].ID() })
	var sb strings.Builder
	for _, sess := range sessions {
		c := uic.client[sess]
		fmt.Fprintf(&sb, "%.0fx%.0f-%b ", c.X, c.Y, int(c.B))
	}
	uic.mu.RUnlock()
	v = template.HTML(html.EscapeString(sb.String())) //#nosec G203
	return
}

func (uic uiClientPos) JawsGetTag(*jaws.Request) any {
	return uiClientPos{}
}

func (g *Globals) ClientPos() jaws.HTMLGetter {
	return uiClientPos{g}
}
