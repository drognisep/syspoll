package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	statusUp = "[green]UP"

	nameCol      = 0
	statusCol    = 1
	intervalCol  = 2
	failCountCol = 3
)

type systemTable struct {
	*tview.Table

	app   *tview.Application
	pages *tview.Pages

	mux     sync.Mutex
	systems []*System
}

func (t *systemTable) Add(sys System) {
	rows := t.GetRowCount()
	t.mux.Lock()
	t.systems = append(t.systems, &sys)
	t.mux.Unlock()
	t.Table.SetCell(rows, nameCol, tview.NewTableCell(sys.Name))
	t.Table.SetCell(rows, statusCol, tview.NewTableCell(Unknown.UiString()))
	t.Table.SetCell(rows, intervalCol, tview.NewTableCell(sys.CheckInterval))
	t.Table.SetCell(rows, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))

	switch {
	case sys.Http != nil:
		go httpPollingLoop(t.app, t.Table, &sys, rows)
	}
}

func NewSystemTable(app *tview.Application, pages *tview.Pages, systems ...System) *systemTable {
	table := tview.NewTable().SetBorders(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetCell(0, nameCol, tview.NewTableCell("System").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, statusCol, tview.NewTableCell("Status").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, intervalCol, tview.NewTableCell("Interval").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, failCountCol, tview.NewTableCell("Failures").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))

	sysTable := &systemTable{
		Table: table,
		app:   app,
		pages: pages,
	}
	for _, sys := range systems {
		sys := sys
		sys.FailedChecks = nil
		sysTable.Add(sys)
	}
	return sysTable
}

func httpPollingLoop(app *tview.Application, table *tview.Table, sys *System, r int) {
	if sys.Http == nil {
		return
	}
	_url, err := sys.Http.ToURL()
	if err != nil {
		return
	}
	dur, err := sys.Interval()
	if err != nil {
		return
	}
	for {
		resp, err := http.Get(_url.String())
		if err != nil {
			sys.FailedChecks = append(sys.FailedChecks, DownFailure(time.Now()))
			table.SetCell(r, statusCol, tview.NewTableCell(Down.UiString()))
			table.SetCell(r, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))
			app.Draw()
			time.Sleep(dur)
			continue
		}
		code := resp.StatusCode
		if code > 399 {
			sys.FailedChecks = append(sys.FailedChecks, ErrorFailure(time.Now()))
			table.SetCell(r, statusCol, tview.NewTableCell(fmt.Sprintf("%s - %d", Error.UiString(), code)))
			table.SetCell(r, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))
		} else {
			table.SetCell(r, statusCol, tview.NewTableCell(fmt.Sprintf("%s - %d", statusUp, code)))
		}
		app.Draw()
		time.Sleep(dur)
	}
}
