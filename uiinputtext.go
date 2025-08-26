package main

import (
	"fmt"
	"strings"

	"github.com/linkdata/jaws"
)

func (g *Globals) InputText() any {
	return jaws.Bind(&g.mu, &g.inputText).SetLocked(func(bind jaws.Binder[string], elem *jaws.Element, v string) (err error) {
		if v != "" && g.SelectPet.Set(v, true) {
			elem.Dirty(g.SelectPet)
		}
		if strings.HasPrefix(v, "fail") {
			if v == "fail" {
				err = fmt.Errorf("whaddayamean, fail?")
			} else {
				// try using cut'n'paste or just holding down 'l'
				g.inputText = "well, if you insist..."
			}
		} else {
			g.inputText = v
		}
		return
	})
}
