package main

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"slices"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	mu        deadlock.RWMutex
	condition int
}

var carMakes = []string{"Dodge", "Hyundai", "Acura", "Volvo", "Saab", "Lada", "Mazda"}
var carModels = []string{"Sedan", "Coupe", "SUV", "Truck", "Cabriolet"}

func intN(n int) int {
	x := rand.Intn(n) //#nosec G404
	return x
}

func AddRandomCar() {
	var vin []byte
	for i := 0; i < 17; i++ {
		n := byte(intN(26 + 10))
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
		condition: 30 + intN(70),
	}
	globals.mu.Lock()
	globals.Cars = append(globals.Cars, car)
	globals.mu.Unlock()
}

func (c *Car) JawsClick(e *jaws.Element, name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch name {
	case "up":
		if idx := slices.Index(globals.Cars, c); idx > 0 {
			globals.Cars[idx-1], globals.Cars[idx] = globals.Cars[idx], globals.Cars[idx-1]
		}
	case "down":
		if idx := slices.Index(globals.Cars, c); idx >= 0 && idx < len(globals.Cars)-1 {
			globals.Cars[idx+1], globals.Cars[idx] = globals.Cars[idx], globals.Cars[idx+1]
		}
	case "remove":
		globals.Cars = slices.DeleteFunc(globals.Cars, func(c2 *Car) bool { return c2 == c })
	case "+":
		if c.condition > 99 {
			return errors.New("condition too high")
		}
		c.condition++
		e.Dirty(c.Condition())
		return nil
	case "-":
		if c.condition < 1 {
			return errors.New("condition too low")
		}
		c.condition--
		e.Dirty(c.Condition())
		return nil
	default:
		return jaws.ErrEventUnhandled
	}

	e.Dirty(globals.CarsTable)
	return nil
}

type uiCondition struct{ *Car }

func (ui uiCondition) JawsGetFloat(e *jaws.Element) (v float64) {
	ui.mu.RLock()
	v = float64(ui.condition)
	ui.mu.RUnlock()
	return
}

func (ui uiCondition) JawsGetHtml(e *jaws.Element) (v template.HTML) {
	ui.mu.RLock()
	v = template.HTML(fmt.Sprint(ui.condition)) //#nosec G203
	ui.mu.RUnlock()
	return
}

func (ui uiCondition) JawsSetFloat(e *jaws.Element, v float64) error {
	ui.mu.Lock()
	ui.condition = int(v)
	ui.mu.Unlock()
	return nil
}

func (c *Car) Condition() jaws.FloatSetter {
	return uiCondition{c}
}
