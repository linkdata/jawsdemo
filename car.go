package main

import (
	"errors"
	"strconv"

	"github.com/linkdata/jaws"
)

type Car struct {
	VIN       string
	Make      string
	Model     string
	Year      int
	Condition int
}

func (c *Car) ConditionID() string {
	return c.VIN + ".cond"
}

func (c *Car) conditionDec(rq *jaws.Request, jid string) error {
	if c.Condition < 1 {
		return errors.New("condition too low")
	}
	c.Condition--
	rq.Jaws.SetInner(c.ConditionID(), strconv.Itoa(c.Condition))
	return nil
}

func (c *Car) RemoveButton() jaws.ClickFn {
	return func(rq *jaws.Request, jid string) error {
		rq.Jaws.Remove(c.VIN)
		return nil
	}
}

func (c *Car) UpButtonFn() jaws.ClickFn {
	return func(rq *jaws.Request, jid string) error {
		rq.Jaws.Insert("carlist", c.VIN, "<tr><td>Meh</td></tr>")
		return nil
	}
}

func (c *Car) AppendButton() jaws.ClickFn {
	return func(rq *jaws.Request, jid string) error {
		rq.Jaws.Append("carlist", "<tr><td>Foo</td></tr>")
		return nil
	}
}

func (c *Car) ConditionDec() jaws.ClickFn {
	return c.conditionDec
}

func (c *Car) conditionInc(rq *jaws.Request, jid string) error {
	if c.Condition > 99 {
		return errors.New("condition too high")
	}
	c.Condition++
	rq.Jaws.SetInner(c.ConditionID(), strconv.Itoa(c.Condition))
	return nil
}

func (c *Car) ConditionInc() jaws.ClickFn {
	return c.conditionInc
}
