package main

import (
	"github.com/linkdata/jaws"
)

type uiClient struct{ *Globals }

func (ui uiClient) getClient(rq *jaws.Request) (c *Client) {
	sess := rq.Session()
	ui.mu.Lock()
	c = ui.client[sess]
	if c == nil {
		c = &Client{}
		ui.client[sess] = c
	}
	ui.mu.Unlock()
	return
}

func (ui uiClient) JawsMakeJsVar(rq *jaws.Request) (v jaws.IsJsVar, err error) {
	return jaws.NewJsVar(&ui.mu, ui.getClient(rq)), nil
}

var _ jaws.JsVarMaker = uiClient{}

func (g *Globals) Client() uiClient {
	return uiClient{g}
}
