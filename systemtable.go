package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/http"
	"strconv"
	"time"
)

const (
	statusUp   = "[green]UP"
	statusDown = "[red]DOWN"
	statusErr  = "[red]ERR"
	statusUnk  = "[gray]UNK"

	nameCol      = 0
	statusCol    = 1
	intervalCol  = 2
	failCountCol = 3
)

func DisplayTable(app *tview.Application) *tview.Table {
	table := tview.NewTable().SetBorders(true)
	table.SetBorderPadding(0, 0, 1, 1)
	table.SetCell(0, nameCol, tview.NewTableCell("System").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, statusCol, tview.NewTableCell("Status").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, intervalCol, tview.NewTableCell("Interval").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))
	table.SetCell(0, failCountCol, tview.NewTableCell("Failures").SetAttributes(tcell.AttrBold).SetAlign(tview.AlignCenter))

	for r, sys := range systems {
		r := r + 1
		sys := sys
		table.SetCell(r, nameCol, tview.NewTableCell(sys.Name))
		table.SetCell(r, statusCol, tview.NewTableCell(statusUnk))
		table.SetCell(r, intervalCol, tview.NewTableCell(sys.CheckInterval))
		table.SetCell(r, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))

		switch {
		case sys.Http != nil:
			go pollingLoop(app, table, sys, r)
		}
	}
	return table
}

func pollingLoop(app *tview.Application, table *tview.Table, sys System, r int) {
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
			sys.FailedChecks = append(sys.FailedChecks, time.Now())
			table.SetCell(r, statusCol, tview.NewTableCell(statusDown))
			table.SetCell(r, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))
			app.Draw()
			time.Sleep(dur)
			continue
		}
		code := resp.StatusCode
		if code > 399 {
			sys.FailedChecks = append(sys.FailedChecks, time.Now())
			table.SetCell(r, statusCol, tview.NewTableCell(fmt.Sprintf("%s - %d", statusErr, code)))
			table.SetCell(r, failCountCol, tview.NewTableCell(strconv.Itoa(len(sys.FailedChecks))).SetAlign(tview.AlignRight))
		} else {
			table.SetCell(r, statusCol, tview.NewTableCell(fmt.Sprintf("%s - %d", statusUp, code)))
		}
		app.Draw()
		time.Sleep(dur)
	}
}
