package main

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"strings"
	"time"

	"github.com/linkdata/jaws"
)

const uiClockID = "clock"

type UiState struct {
	G *Globals
	*jaws.Request
}

func (uis *UiState) ClockID() (v string) {
	return uiClockID
}

func ClockString() string {
	now := time.Now()
	return fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
}

func (uis *UiState) ClockString() (v string) {
	return ClockString()
}

const uiCarsLinkID = "cars"

func (*UiState) CarsLinkID() string {
	return uiCarsLinkID
}

func CarsLinkText() string {
	switch rand.Intn(5) {
	case 0:
		return "Check out these cars!"
	case 1:
		return "Did you know VIN numbers are encoded?"
	case 2:
		return "DO NOT CLICK HERE!"
	case 3:
		return "Cars"
	}
	return "This is a boring link to car info."
}

func (*UiState) CarsLinkText() string {
	return CarsLinkText()
}

func (uis *UiState) JawsVersion() (v string) {
	if bi, ok := debug.ReadBuildInfo(); ok {
		v = bi.Main.Version
		for _, dep := range bi.Deps {
			if strings.HasSuffix(dep.Path, "/jaws") {
				v += " - jaws@" + dep.Version
			}
		}
	}
	return
}

func (uis *UiState) OnMystical() jaws.ClickFn {
	return func(rq *jaws.Request, jid string) error {
		rq.Trigger("setinputbutton", "Wooooo....")
		return nil
	}
}

func NewUiState(rq *jaws.Request, g *Globals) *UiState {
	uis := &UiState{
		G:       g,
		Request: rq,
	}
	return uis
}
