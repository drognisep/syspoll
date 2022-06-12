package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

// DirBox is a collection Primitive that is a bit like tview.Flex, but does not try to fill the given space by default.
// DirBox will pack Primitives along the axis at their minimum reported size, and  will attempt to make them match their dimensions along the cross axis.
// DirBox will not allow a Primitive to exceed its own cross axis dimension.
type DirBox struct {
	*tview.Box

	horizontal  bool
	fullScreen  bool
	padding     int
	expandCross bool

	mux   sync.RWMutex
	items []tview.Primitive
}

// Row will return a horizontal DirBox.
func Row(items ...tview.Primitive) *DirBox {
	return &DirBox{
		Box:        tview.NewBox(),
		horizontal: true,
		items:      items,
	}
}

// Col will return a vertical DirBox.
func Col(items ...tview.Primitive) *DirBox {
	return &DirBox{
		Box:        tview.NewBox(),
		horizontal: false,
		items:      items,
	}
}

// SetFullScreen will allow the user to specify whether the DirBox should use the entire available screen space when drawing.
func (b *DirBox) SetFullScreen(fullscreen bool) *DirBox {
	b.fullScreen = fullscreen
	return b
}

// SetItemPadding will set the space between individual Primitives along the axis.
func (b *DirBox) SetItemPadding(padding int) *DirBox {
	if padding >= 0 {
		b.padding = padding
	}
	return b
}

// SetExpandCrossAxis will allow the user to specify that this DirBox should expand its elements along the cross axis as much as possible within its own allotted space.
func (b *DirBox) SetExpandCrossAxis(expand bool) *DirBox {
	b.expandCross = expand
	return b
}

// Draw will draw the DirBox and the Primitive items it contains.
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

	if b.expandCross && b.horizontal {
		return maxHeight
	}

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

	if b.expandCross && !b.horizontal {
		return maxWidth
	}

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
