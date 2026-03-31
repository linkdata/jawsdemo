package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/ui"
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
	case "mystical":
		e.Session().Set("mystical", nil)
	case "add":
		AddRandomCar()
		e.Dirty(globals.CarsTable)
		return nil
	}
	return jaws.ErrEventUnhandled
}
