package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/linkdata/jaws"
)

type uiClock struct{}

func (uiClock) JawsGetHtml(rq *jaws.Request) template.HTML {
	now := time.Now()
	flash := "&nbsp;"
	if now.Second()%2 == 0 {
		flash = ":"
	}
	return template.HTML(fmt.Sprintf("%02d%s%02d", now.Hour(), flash, now.Minute()))
}
