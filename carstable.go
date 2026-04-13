package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/ui"
)

type CarsTable struct{}

func (ct *CarsTable) JawsContains(e *jaws.Element) (tl []jaws.UI) {
	globals.mu.RLock()
	defer globals.mu.RUnlock()
	for _, c := range globals.Cars {
		tl = append(tl, ui.NewTemplate("car_row.html", c))
	}
	tl = append(tl, ui.NewTemplate("car_row.html", nil))
	return tl
}

func (ct *CarsTable) JawsClick(e *jaws.Element, name string) (err error) {
	switch name {
	case "add":
		AddRandomCar()
		e.Dirty(globals.CarsTable)
		return nil
	}
	return jaws.ErrEventUnhandled
}

func (ct *CarsTable) Mystical() jaws.ClickHandler {
	return ui.Clickable("Mystical", func(elem *jaws.Element, name string) (err error) {
		elem.Session().Set("mystical", nil)
		return
	})
}
