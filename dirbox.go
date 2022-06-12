package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type DirBox struct {
	*tview.Box

	horizontal bool
	fullScreen bool
	padding    int

	mux   sync.RWMutex
	items []tview.Primitive
}

func Row(items ...tview.Primitive) *DirBox {
	return &DirBox{
		Box:        tview.NewBox(),
		horizontal: true,
		items:      items,
	}
}

func Col(items ...tview.Primitive) *DirBox {
	return &DirBox{
		Box:        tview.NewBox(),
		horizontal: false,
		items:      items,
	}
}

func (b *DirBox) SetFullScreen(fullscreen bool) *DirBox {
	b.fullScreen = fullscreen
	return b
}

func (b *DirBox) SetItemPadding(padding int) *DirBox {
	if padding >= 0 {
		b.padding = padding
	}
	return b
}

func (b *DirBox) Draw(screen tcell.Screen) {
	b.Box.DrawForSubclass(screen, b)
	var x, y int
	if b.fullScreen {
		width, height := screen.Size()
		b.SetRect(0, 0, width, height)
	}

	x, y, _, _ = b.GetInnerRect()
	var dims [][4]int
	dims = b.calcDims(x, y)

	b.mux.Lock()
	defer b.mux.Unlock()
	for i, item := range b.items {
		item.SetRect(dims[i][0], dims[i][1], dims[i][2], dims[i][3])
		if item.HasFocus() {
			defer item.Draw(screen)
		} else {
			item.Draw(screen)
		}
	}
}

// Focus is called when this primitive receives focus.
func (b *DirBox) Focus(delegate func(p tview.Primitive)) {
	for _, item := range b.items {
		if item != nil && item.HasFocus() {
			delegate(item)
			return
		}
	}
	b.Box.Focus(delegate)
}

// HasFocus returns whether or not this primitive has focus.
func (b *DirBox) HasFocus() bool {
	for _, item := range b.items {
		if item != nil && item.HasFocus() {
			return true
		}
	}
	return b.Box.HasFocus()
}

// MouseHandler returns the mouse handler for this primitive.
func (b *DirBox) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return b.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !b.InRect(event.Position()) {
			return false, nil
		}

		// Pass mouse events along to the first child item that takes it.
		for _, item := range b.items {
			if item == nil {
				continue
			}
			consumed, capture = item.MouseHandler()(action, event, setFocus)
			if consumed {
				return
			}
		}

		return
	})
}

// InputHandler returns the handler for this primitive.
func (b *DirBox) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		for _, item := range b.items {
			if item != nil && item.HasFocus() {
				if handler := item.InputHandler(); handler != nil {
					handler(event, setFocus)
					return
				}
			}
		}
	})
}

func (b *DirBox) calcDims(x, y int) [][4]int {
	b.mux.RLock()
	defer b.mux.RUnlock()
	if b.horizontal {
		return b.calcHDims(x, y, b.maxItemOrBoxHeight())
	} else {
		return b.calcVDims(x, y, b.maxItemOrBoxWidth())
	}
}

func (b *DirBox) calcHDims(offX, offY, h int) [][4]int {
	b.mux.RLock()
	defer b.mux.RUnlock()

	var runningW int
	var padding int
	dims := make([][4]int, len(b.items))
	for i, item := range b.items {
		if i > 0 {
			padding = b.padding
		}
		_, _, w, _ := item.GetRect()
		dims[i] = [4]int{offX + runningW + padding, offY, w, h}
		runningW += w + padding
	}
	return dims
}

func (b *DirBox) calcVDims(offX, offY, w int) [][4]int {
	b.mux.RLock()
	defer b.mux.RUnlock()

	var runningH int
	var padding int
	dims := make([][4]int, len(b.items))
	for i, item := range b.items {
		if i > 0 {
			padding = b.padding
		}
		_, _, _, h := item.GetRect()
		dims[i] = [4]int{offX, offY + runningH + padding, w, h}
		runningH += h + padding
	}
	return dims
}

func (b *DirBox) maxItemOrBoxHeight() int {
	_, _, _, maxHeight := b.GetInnerRect()

	var runningMax int
	for _, item := range b.items {
		_, _, _, h := item.GetRect()
		if h > runningMax {
			runningMax = h
		}
	}
	if runningMax > maxHeight {
		return maxHeight
	}
	return runningMax
}

func (b *DirBox) maxItemOrBoxWidth() int {
	_, _, maxWidth, _ := b.GetInnerRect()

	var runningMax int
	for _, item := range b.items {
		_, _, w, _ := item.GetRect()
		if w > runningMax {
			runningMax = w
		}
	}
	if runningMax > maxWidth {
		return maxWidth
	}
	return runningMax
}
