package main

import (
	"github.com/linkdata/jaws"
)

type uiClient struct{ *Globals }

func (ui uiClient) JawsGetAny(e *jaws.Element) any {
	ui.mu.RLock()
	defer ui.mu.RUnlock()
	if v, ok := ui.client[e.Session()]; ok {
		return *v
	}
	return Client{}
}

func (ui uiClient) JawsSetAny(e *jaws.Element, v any) error {
	if v, ok := v.(map[string]any); ok {
		sess := e.Session()
		ui.mu.Lock()
		c := ui.client[sess]
		if c == nil {
			c = &Client{}
			ui.client[sess] = c
		}
		c.X = v["X"].(float64)
		c.Y = v["Y"].(float64)
		c.B = v["B"].(float64)
		ui.mu.Unlock()
	}
	return nil
}

func (g *Globals) Client() jaws.AnySetter {
	return uiClient{g}
}
