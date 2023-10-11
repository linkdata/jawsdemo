package main

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"

	"github.com/linkdata/deadlock"
	"github.com/linkdata/jaws"
)

type CarsTable struct{}

func (ct *CarsTable) JawsContains(rq *jaws.Request) (tl []jaws.UI) {
	for _, c := range globals.Cars {
		tl = append(tl, rq.MakeTemplate("car_row.html", c))
	}
	tl = append(tl, rq.MakeTemplate("car_row.html", nil))
	return tl
}

func (ct *CarsTable) JawsClick(e *jaws.Element, name string) (stop bool, err error) {
	switch name {
	case "add":
		AddRandomCar()
		e.Dirty(globals.CarsTable)
		stop = true
	}
	return
}

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

func AddRandomCar() {
	var vin []byte
	for i := 0; i < 17; i++ {
		n := byte(rand.Intn(26 + 10))
		if n < 10 {
			vin = append(vin, '0'+n)
		} else {
			vin = append(vin, 'A'+(n-10))
		}
	}
	car := &Car{
		VIN:       string(vin),
		Make:      carMakes[rand.Intn(len(carMakes))],
		Model:     carModels[rand.Intn(len(carModels))],
		Year:      1970 + rand.Intn(30),
		condition: 30 + rand.Intn(70),
	}
	globals.Cars = append(globals.Cars, car)
}

func (c *Car) JawsClick(e *jaws.Element, name string) (stop bool, err error) {
	stop = true
	switch name {
	case "up":
		jaws.ListMove(globals.Cars, c, -1)
	case "down":
		jaws.ListMove(globals.Cars, c, 1)
	case "remove":
		globals.Cars = jaws.ListRemove(globals.Cars, c)
	case "+":
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.condition > 99 {
			return true, errors.New("condition too high")
		}
		c.condition++
		e.Dirty(c.Condition())
		return
	case "-":
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.condition < 1 {
			return true, errors.New("condition too low")
		}
		c.condition--
		e.Dirty(c.Condition())
		return
	default:
		stop = false
		return
	}

	e.Dirty(globals.CarsTable)
	return
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
	v = template.HTML(fmt.Sprint(ui.condition))
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
