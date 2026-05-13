package main

import (
	"html/template"
	"math"
	"strconv"
	"strings"

	"github.com/linkdata/jaws"
	"github.com/linkdata/jaws/lib/bind"
	"github.com/linkdata/jaws/lib/tag"
)

const (
	drawingWidth     = 640
	drawingHeight    = 320
	drawingMaxPoints = 240
)

var drawingColors = [...]string{
	"#0d6efd",
	"#d63384",
	"#198754",
	"#fd7e14",
	"#6f42c1",
	"#20c997",
	"#dc3545",
	"#0dcaf0",
}

type drawingPoint struct {
	X float64
	Y float64
}

type drawingTrace struct {
	SessionID uint64
	Color     string
	Points    []drawingPoint
}

type uiDrawing struct {
	globals *Globals
}

// Drawing returns the shared SVG drawing board.
func (g *Globals) Drawing() bind.HTMLGetter {
	return uiDrawing{globals: g}
}

func (ui uiDrawing) JawsGetTag(tag.Context) any {
	return &ui.globals.drawingTraces
}

func (ui uiDrawing) JawsGetHTML(elem *jaws.Element) (v template.HTML) {
	sessionID := drawingSessionID(elem)
	ui.globals.mu.RLock()
	traces := copyDrawingTraces(ui.globals.drawingTraces)
	ui.globals.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString(`<rect x="0" y="0" width="640" height="320" rx="8" fill="#f8f9fa"></rect>`)
	sb.WriteString(`<path d="M0 40H640M0 80H640M0 120H640M0 160H640M0 200H640M0 240H640M0 280H640M40 0V320M80 0V320M120 0V320M160 0V320M200 0V320M240 0V320M280 0V320M320 0V320M360 0V320M400 0V320M440 0V320M480 0V320M520 0V320M560 0V320M600 0V320" stroke="#dee2e6" stroke-width="1"></path>`)
	pointCount := 0
	currentColor := ""
	for _, trace := range traces {
		if trace.SessionID == sessionID {
			currentColor = trace.Color
		}
		if len(trace.Points) == 0 {
			continue
		}
		pointCount += len(trace.Points)
		sb.WriteString(`<polyline points="`)
		for i, pt := range trace.Points {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(formatDrawingCoord(pt.X))
			sb.WriteByte(',')
			sb.WriteString(formatDrawingCoord(pt.Y))
		}
		sb.WriteString(`" fill="none" stroke="`)
		sb.WriteString(trace.Color)
		sb.WriteString(`" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"></polyline>`)
		for _, pt := range trace.Points {
			sb.WriteString(`<circle cx="`)
			sb.WriteString(formatDrawingCoord(pt.X))
			sb.WriteString(`" cy="`)
			sb.WriteString(formatDrawingCoord(pt.Y))
			sb.WriteString(`" r="5" fill="`)
			sb.WriteString(trace.Color)
			sb.WriteString(`" stroke="#ffffff" stroke-width="2"></circle>`)
		}
	}
	if pointCount == 0 {
		sb.WriteString(`<text x="320" y="164" text-anchor="middle" fill="#6c757d" font-size="18">Drag to draw your trace. Right-click to clear it.</text>`)
	}
	if currentColor != "" {
		sb.WriteString(`<circle cx="24" cy="24" r="8" fill="`)
		sb.WriteString(currentColor)
		sb.WriteString(`" stroke="#ffffff" stroke-width="2"></circle>`)
		sb.WriteString(`<text x="40" y="29" fill="#495057" font-size="14">your trace</text>`)
	}
	return template.HTML(sb.String()) // #nosec G203
}

func (ui uiDrawing) JawsPointer(elem *jaws.Element, ptr jaws.Pointer) (err error) {
	if ptr.Name != "drawing" {
		return jaws.ErrEventUnhandled
	}
	switch {
	case ptr.Kind == jaws.PointerDown && ptr.Button == 0:
	case ptr.Kind == jaws.PointerMove && ptr.Buttons&jaws.PointerButtonPrimary != 0:
	default:
		return nil
	}
	ui.addDrawingPoint(elem, ptr.X, ptr.Y)
	return nil
}

func (ui uiDrawing) addDrawingPoint(elem *jaws.Element, x, y float64) {
	sessionID := drawingSessionID(elem)
	ui.globals.mu.Lock()
	trace := ui.globals.drawingTraceLocked(sessionID)
	trace.Points = append(trace.Points, drawingPoint{
		X: clampDrawingCoord(x, drawingWidth),
		Y: clampDrawingCoord(y, drawingHeight),
	})
	if len(trace.Points) > drawingMaxPoints {
		copy(trace.Points, trace.Points[len(trace.Points)-drawingMaxPoints:])
		trace.Points = trace.Points[:drawingMaxPoints]
	}
	ui.globals.mu.Unlock()

	elem.Dirty(ui)
}

func (ui uiDrawing) JawsContextMenu(elem *jaws.Element, click jaws.Click) (err error) {
	if click.Name != "drawing" {
		return jaws.ErrEventUnhandled
	}

	sessionID := drawingSessionID(elem)
	ui.globals.mu.Lock()
	if trace := ui.globals.findDrawingTraceLocked(sessionID); trace != nil {
		trace.Points = trace.Points[:0]
	}
	ui.globals.mu.Unlock()

	elem.Dirty(ui)
	return nil
}

func (g *Globals) drawingTraceLocked(sessionID uint64) *drawingTrace {
	if trace := g.findDrawingTraceLocked(sessionID); trace != nil {
		return trace
	}
	g.drawingTraces = append(g.drawingTraces, drawingTrace{
		SessionID: sessionID,
		Color:     drawingColors[len(g.drawingTraces)%len(drawingColors)],
	})
	return &g.drawingTraces[len(g.drawingTraces)-1]
}

func (g *Globals) findDrawingTraceLocked(sessionID uint64) *drawingTrace {
	for i := range g.drawingTraces {
		if g.drawingTraces[i].SessionID == sessionID {
			return &g.drawingTraces[i]
		}
	}
	return nil
}

func copyDrawingTraces(src []drawingTrace) (dst []drawingTrace) {
	dst = make([]drawingTrace, len(src))
	for i := range src {
		dst[i] = src[i]
		dst[i].Points = append([]drawingPoint(nil), src[i].Points...)
	}
	return
}

func drawingSessionID(elem *jaws.Element) uint64 {
	if elem == nil {
		return 0
	}
	return elem.Session().ID()
}

func clampDrawingCoord(v, max float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0
	}
	if v > max {
		return max
	}
	return v
}

func formatDrawingCoord(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}

var _ bind.HTMLGetter = uiDrawing{}
var _ jaws.ContextMenuHandler = uiDrawing{}
var _ jaws.PointerHandler = uiDrawing{}
var _ tag.TagGetter = uiDrawing{}
