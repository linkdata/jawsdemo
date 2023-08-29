package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/linkdata/jaws"
)

type UiState struct {
	G *Globals
	*jaws.Request
}

func ClockString() string {
	now := time.Now()
	return fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
}

func (uis *UiState) ClockString() (v string) {
	return ClockString()
}

func RandomizeCarsLink() {
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
