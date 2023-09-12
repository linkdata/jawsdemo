package main

import (
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/linkdata/jaws"
)

type CarsTable struct{}

func (ct *CarsTable) JawsContains(rq *jaws.Request) (tl []jaws.Template) {
	for _, c := range globals.Cars {
		tl = append(tl, rq.NewTemplate("car_row.html", c))
	}
	tl = append(tl, rq.NewTemplate("car_row.html", nil))
	return tl
}

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	condition atomic.Value
}

func (c *Car) JawsClick(e *jaws.Element, name string) error {
	switch name {
	case "up":
		jaws.ListMove(globals.Cars, c, -1)
	case "down":
		jaws.ListMove(globals.Cars, c, 1)
	case "remove":
		globals.Cars = jaws.ListRemove(globals.Cars, c)
		e.Jaws.Remove(c)
		return nil
	case "+":
		oldVal := c.condition.Load().(int)
		if oldVal > 99 {
			return errors.New("condition too high")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal+1) {
			e.Jaws.Dirty(c.Condition())
		}
		return nil
	case "-":
		oldVal := c.condition.Load().(int)
		if oldVal < 1 {
			return errors.New("condition too low")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal-1) {
			e.Jaws.Dirty(c.Condition())
		}
		return nil
	}
	e.Jaws.Dirty(globals.CarsTable)
	return nil
}

func (c *Car) Condition() *atomic.Value {
	if c.condition.Load() == nil {
		c.condition.Store(rand.Intn(100))
	}
	return &c.condition
}
