package main

import (
	"fmt"
	"html"
	"html/template"
	"slices"
	"sort"
	"strings"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
	"github.com/linkdata/jaws/lib/jtag"
)

type uiClientPos struct{ *Globals }

func (uic uiClientPos) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	var activeclientsessions []*jaws.Session
	sessions := e.Jaws.Sessions()
	uic.mu.RLock()
	defer uic.mu.RUnlock()
	for sess, c := range uic.client {
		if slices.Contains(sessions, sess) {
			if c.X != -1 || c.Y != -1 {
				activeclientsessions = append(activeclientsessions, sess)
			}
		} else {
			delete(uic.client, sess)
		}
	}
	sort.Slice(activeclientsessions, func(i, j int) bool { return activeclientsessions[i].ID() < activeclientsessions[j].ID() })
	var sb strings.Builder
	fmt.Fprintf(&sb, "(%d/%d)", len(activeclientsessions), len(uic.client))
	for _, sess := range activeclientsessions {
		c := uic.client[sess]
		if c.X != -1 || c.Y != -1 {
			fmt.Fprintf(&sb, " %.0fx%.0f-%b", c.X, c.Y, int(c.B))
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
