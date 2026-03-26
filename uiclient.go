package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/ui"
)

type uiClient struct{ *Globals }

func (uic uiClient) getClient(rq *jaws.Request) (c *Client) {
	sess := rq.Session()
	uic.mu.Lock()
	c = uic.client[sess]
	if c == nil {
		c = &Client{}
		uic.client[sess] = c
	}
	uic.mu.Unlock()
	return
}

func (uic uiClient) JawsMakeJsVar(rq *jaws.Request) (v ui.IsJsVar, err error) {
	return ui.NewJsVar(&uic.mu, uic.getClient(rq)), nil
}

var _ ui.JsVarMaker = uiClient{}

func (g *Globals) Client() uiClient {
	return uiClient{g}
}
