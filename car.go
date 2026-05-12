package main

import (
	"errors"
	"math/rand"
	"slices"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
)

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	mu        deadlock.RWMutex
	condition float64
}

var carMakes = []string{"Dodge", "Hyundai", "Acura", "Volvo", "Saab", "Lada", "Mazda"}
var carModels = []string{"Sedan", "Coupe", "SUV", "Truck", "Cabriolet"}

func intN(n int) int {
	x := rand.Intn(n) //#nosec G404
	return x
}

func (g *Globals) AddRandomCar() *Car {
	var vin []byte
	for range 17 {
		n := byte(intN(26 + 10)) // #nosec G115
		if n < 10 {
			vin = append(vin, '0'+n)
		} else {
			vin = append(vin, 'A'+(n-10))
		}
	}
	car := &Car{
		VIN:       string(vin),
		Make:      carMakes[intN(len(carMakes))],
		Model:     carModels[intN(len(carModels))],
		Year:      1970 + intN(30),
		condition: 30 + float64(intN(70)),
	}
	g.mu.Lock()
	g.Cars = append(g.Cars, car)
	g.mu.Unlock()
	return car
}

func AddRandomCar() {
	globals.AddRandomCar()
}

func (c *Car) JawsClick(e *jaws.Element, data jaws.Click) error {
	switch data.Name {
	case "up":
		var changed bool
		globals.mu.Lock()
		if idx := slices.Index(globals.Cars, c); idx > 0 {
			globals.Cars[idx-1], globals.Cars[idx] = globals.Cars[idx], globals.Cars[idx-1]
			changed = true
		}
		globals.mu.Unlock()
		if changed {
			e.Dirty(globals.CarsTable)
		}
		return nil
	case "down":
		var changed bool
		globals.mu.Lock()
		if idx := slices.Index(globals.Cars, c); idx >= 0 && idx < len(globals.Cars)-1 {
			globals.Cars[idx+1], globals.Cars[idx] = globals.Cars[idx], globals.Cars[idx+1]
			changed = true
		}
		globals.mu.Unlock()
		if changed {
			e.Dirty(globals.CarsTable)
		}
		return nil
	case "remove":
		var changed bool
		globals.mu.Lock()
		before := len(globals.Cars)
		globals.Cars = slices.DeleteFunc(globals.Cars, func(c2 *Car) bool { return c2 == c })
		changed = len(globals.Cars) != before
		globals.mu.Unlock()
		if changed {
			e.Dirty(globals.CarsTable)
		}
		return nil
	case "+":
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.condition >= 100 {
			return errors.New("condition too high")
		}
		c.condition++
		e.Dirty(c.Condition())
		return nil
	case "-":
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.condition <= 0 {
			return errors.New("condition too low")
		}
		c.condition--
		e.Dirty(c.Condition())
		return nil
	default:
		return jaws.ErrEventUnhandled
	}
}

func (c *Car) Condition() any {
	return bind.New(&c.mu, &c.condition)
}
