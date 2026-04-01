package main

import (
	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/ui"
)

type Client struct {
	X float64
	Y float64
	B float64 // button state
}

var _ ui.SetPather = (*Client)(nil)

func (c *Client) JawsPathSet(elem *jaws.Element, jspath string, value any) {
	elem.Dirty(uiClientPos{})
}
