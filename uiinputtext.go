package main

import (
	"fmt"
	"strings"

	"github.com/linkdata/jaws"
)

func (g *Globals) InputText() any {
	return jaws.Bind(&g.mu, &g.inputText).
		SetLocked(func(bind jaws.Binder[string], elem *jaws.Element, v string) error {
			// if user entered a pet, change the dropdown to that pet
			if v != "" && g.SelectPet.Count(v) > 0 {
				g.SelectPet.Set(v, true)
				elem.Dirty(g.SelectPet)
			}

			// make it harder to enter text that starts with "fail"
			// try using cut'n'paste or just holding down 'l'
			if strings.HasPrefix(v, "fail") {
				if v == "fail" {
					return fmt.Errorf("whaddayamean, fail?")
				}
				v = "well, if you insist..."
			}

			return bind.JawsSetLocked(elem, v)
		})
}
