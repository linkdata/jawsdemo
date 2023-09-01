package main

import (
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/linkdata/jaws"
)

type CarList struct {
	Cars []*Car
}

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	condition atomic.Value
}

func (c *Car) Condition() *atomic.Value {
	if c.condition.Load() == nil {
		c.condition.Store(rand.Intn(100))
	}
	return &c.condition
}

func (c *Car) JawsClick(e *jaws.Element, name string) error {
	switch name {
	case "up":
		for i, oc := range globals.Cars {
			if i > 0 && oc == c {
				globals.Cars[i], globals.Cars[i-1] = globals.Cars[i-1], globals.Cars[i]
				break
			}
		}
	case "remove":
		var nl []*Car
		for _, oc := range globals.Cars {
			if oc != c {
				nl = append(nl, oc)
			}
		}
		globals.Cars = nl
		e.Request().Jaws.Remove(c)
		return nil
	case "down":
		for i, oc := range globals.Cars {
			if i < len(globals.Cars)-1 && oc == c {
				globals.Cars[i], globals.Cars[i+1] = globals.Cars[i+1], globals.Cars[i]
				break
			}
		}
	case "+":
		oldVal := c.condition.Load().(int)
		if oldVal > 99 {
			return errors.New("condition too high")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal+1) {
			e.Request().Jaws.Update([]interface{}{c.Condition()})
		}
		return nil
	case "-":
		oldVal := c.condition.Load().(int)
		if oldVal < 1 {
			return errors.New("condition too low")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal-1) {
			e.Request().Jaws.Update([]interface{}{c.Condition()})
		}
		return nil
	}
	var tags []interface{}
	tags = append(tags, "carlist")
	for _, oc := range globals.Cars {
		tags = append(tags, oc)
	}
	e.Request().Jaws.Order(tags...)
	return nil
}
