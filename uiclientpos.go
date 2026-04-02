package main

import (
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
	"github.com/linkdata/jaws/lib/jtag"
)

type uiClientPos struct{ *Globals }

func (uic uiClientPos) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	var sessions []*jaws.Session
	uic.mu.RLock()
	defer uic.mu.RUnlock()
	for sess := range uic.client {
		sessions = append(sessions, sess)
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID() < sessions[j].ID() })
	var sb strings.Builder
	for _, sess := range sessions {
		c := uic.client[sess]
		if c.X != -1 || c.Y != -1 {
			fmt.Fprintf(&sb, "%.0fx%.0f-%b ", c.X, c.Y, int(c.B))
		}
	}
	v = template.HTML(html.EscapeString(sb.String())) //#nosec G203
	return
}

func (uic uiClientPos) JawsGetTag(jtag.Context) any {
	return uiClientPos{}
}

func (g *Globals) ClientPos() bind.HTMLGetter {
	return uiClientPos{g}
}
