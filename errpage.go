package main

import (
	"fmt"
	"github.com/drognisep/syspoll/page"
	"github.com/rivo/tview"
)

func ShowErr(pages *tview.Pages, err error) {
	errModel := tview.NewModal().
		SetText(fmt.Sprintf("Error occurred:\n%v", err)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage(page.Error)
		})
	pages.AddAndSwitchToPage(page.Error, GridWrapper(errModel, 50, 20), true)
}
