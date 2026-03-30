package main

import (
	"github.com/linkdata/jaws"
)

func (g *Globals) InputButton() jaws.Binder[string] {
	rv := jaws.Bind(&g.mu, &g.inputButton).
		GetLocked(func(bind jaws.Binder[string], elem *jaws.Element) (value string) {
			value = g.inputButton
			if elem.Session().Get("mystical") != nil {
				elem.SetAttr("disabled", "")
			} else {
				elem.RemoveAttr("disabled")
			}
			return
		}).
		Clicked(func(elem *jaws.Element, name string) (err error) {
			err = jaws.ErrEventUnhandled
			if name == "clicky" {
				err = nil
				g.mu.Lock()
				defer g.mu.Unlock()
				if g.inputButton == "Meh" {
					g.inputButton = "Mystical?"
					elem.Session().Set("mystical", true)
				} else {
					g.inputButton = "Meh"
					elem.Session().Set("mystical", nil)
				}
				elem.Dirty(g.InputButton())
			}
			return
		})
	return rv
}
