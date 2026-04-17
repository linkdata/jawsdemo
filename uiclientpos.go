package main

import (
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
	"github.com/linkdata/jaws/lib/tag"
)

type uiClientPos struct{ *Globals }

func (uic uiClientPos) JawsGetHTML(e *jaws.Element) (v template.HTML) {
	var activeclients []*Client
	sessions := e.Jaws.Sessions()
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID() < sessions[j].ID() })
	uic.mu.RLock()
	defer uic.mu.RUnlock()
	for _, sess := range sessions {
		if c, _ := sess.Get("client").(*Client); c != nil {
			if c.X != -1 || c.Y != -1 {
				activeclients = append(activeclients, c)
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "(%d/%d)", len(activeclients), len(sessions))
	for _, c := range activeclients {
		fmt.Fprintf(&sb, " %.0fx%.0f-%b", c.X, c.Y, int(c.B))
	}
	v = template.HTML(html.EscapeString(sb.String())) //#nosec G203
	return
}

func (uic uiClientPos) JawsGetTag(tag.Context) any {
	return uiClientPos{}
}

func (g *Globals) ClientPos() bind.HTMLGetter {
	return uiClientPos{g}
}
