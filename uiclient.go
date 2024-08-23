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

func (ui uiClient) JawsVarMake(rq *jaws.Request) (v jaws.UI, err error) {
	return jaws.NewJsVar("client", jaws.Bind(&ui.mu, ui.getClient(rq))), nil
}

func (c *Client) JawsGetTag(*jaws.Request) any {
	return uiClientPos{}
}

var _ jaws.VarMaker = uiClient{}

func (g *Globals) Client() uiClient {
	return uiClient{g}
}
