package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
)

func (g *Globals) InputButton() bind.Binder[string] {
	rv := bind.New(&g.mu, &g.inputButton).
		GetLocked(func(bind bind.Binder[string], elem *jaws.Element) (value string) {
			value = g.inputButton
			if elem.Session().Get("mystical") != nil {
				elem.SetAttr("disabled", "")
			} else {
				elem.RemoveAttr("disabled")
			}
			return
		}).
		Clicked(func(bind bind.Binder[string], elem *jaws.Element, name string) (err error) {
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
				elem.Dirty(bind)
			}
			return
		})
	return rv
}
