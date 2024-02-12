package main

import "github.com/linkdata/jaws"

type CarsTable struct{}

func (ct *CarsTable) JawsContains(e *jaws.Element) (tl []jaws.UI) {
	for _, c := range globals.Cars {
		tl = append(tl, e.MakeTemplate("car_row.html", c))
	}
	tl = append(tl, e.MakeTemplate("car_row.html", nil))
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
