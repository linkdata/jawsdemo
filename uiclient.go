package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/ui"
)

type uiClient struct{ *Globals }

func (uic uiClient) getClient(rq *jaws.Request) (c *Client) {
	sess := rq.Session()
	if c, _ = sess.Get("client").(*Client); c == nil {
		c = &Client{
			X: -1,
			Y: -1,
			B: 0,
		}
		sess.Set("client", c)
	}
	return
}

func (uic uiClient) JawsMakeJsVar(rq *jaws.Request) (v ui.IsJsVar, err error) {
	return ui.NewJsVar(&uic.mu, uic.getClient(rq)), nil
}

var _ ui.JsVarMaker = uiClient{}

func (g *Globals) Client() uiClient {
	return uiClient{g}
}
