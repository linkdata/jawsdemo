package main

import (
	"errors"
	"math/rand"
	"sync/atomic"

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

func (ct *CarsTable) JawsClick(e *jaws.Element, name string) error {
	switch name {
	case "add":
		AddRandomCar()
		e.Dirty(globals.CarsTable)
	}
	return nil
}

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	condition atomic.Value
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
		VIN:   string(vin),
		Make:  carMakes[rand.Intn(len(carMakes))],
		Model: carModels[rand.Intn(len(carModels))],
		Year:  1970 + rand.Intn(30),
	}
	globals.Cars = append(globals.Cars, car)
}

func (c *Car) JawsClick(e *jaws.Element, name string) error {
	switch name {
	case "up":
		jaws.ListMove(globals.Cars, c, -1)
	case "down":
		jaws.ListMove(globals.Cars, c, 1)
	case "remove":
		globals.Cars = jaws.ListRemove(globals.Cars, c)
	case "+":
		oldVal := c.condition.Load().(int)
		if oldVal > 99 {
			return errors.New("condition too high")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal+1) {
			e.Dirty(c.Condition())
		}
		return nil
	case "-":
		oldVal := c.condition.Load().(int)
		if oldVal < 1 {
			return errors.New("condition too low")
		}
		if c.condition.CompareAndSwap(oldVal, oldVal-1) {
			e.Dirty(c.Condition())
		}
		return nil
	}

	e.Dirty(globals.CarsTable)
	return nil
}

func (c *Car) Condition() *atomic.Value {
	if c.condition.Load() == nil {
		c.condition.Store(rand.Intn(100))
	}
	return &c.condition
}
