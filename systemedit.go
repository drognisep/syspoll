package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/url"
	"regexp"
	"time"
)

type editCallback func(newState *System, submitted bool)

func ShowSystemCreate(pages *tview.Pages, callback editCallback) {
	ShowSystemEdit(pages, nil, callback)
}

func ShowSystemEdit(pages *tview.Pages, system *System, callback editCallback) {
	const pageName = "Edit System"

	system = system.Copy()
	if system == nil {
		system = &System{
			Http: &CheckHttp{},
		}
	}
	if system.Http == nil {
		system.Http = &CheckHttp{}
	}

	protoOpt := "https"

	form := tview.NewForm().SetButtonBackgroundColor(tcell.ColorBlue).SetButtonTextColor(tcell.ColorWhite)
	form.SetBorder(true)
	form.AddInputField("System Name", system.Name, 0, nil, func(text string) {
		system.Name = text
	})
	form.AddInputField("Check Interval", system.CheckInterval, 0, nil, func(text string) {
		system.CheckInterval = text
	})
	form.AddDropDown("Protocol", []string{"http", "https"}, 1, func(option string, optionIndex int) {
		protoOpt = option
	})
	form.AddInputField("URL", system.Http.URL, 0, nil, func(text string) {
		_url, err := url.Parse(text)
		if err != nil {
			return
		}
		_url.Scheme = protoOpt
		system.Http.URL = _url.String()
	})

	form.AddButton("Save", func() {
		valid := true
		valid = nameValidFunc(system.Name, 'a')
		valid = valid && durValidFunc(system.CheckInterval, 'a')
		valid = valid && urlValidFunc(system.Http.URL, 'a')
		if !valid {
			return
		}

		pages.RemovePage(pageName)
		callback(system, true)
	})
	form.AddButton("Cancel", func() {
		pages.RemovePage(pageName)
		callback(nil, false)
	})
	wrapper := GridWrapper(form, 50, 13)

	pages.AddAndSwitchToPage(pageName, wrapper, true)
}

func GridWrapper(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func urlValidFunc(textToCheck string, _ rune) bool {
	_, err := url.Parse(textToCheck)
	return err == nil
}

func durValidFunc(textToCheck string, _ rune) bool {
	_, err := time.ParseDuration(textToCheck)
	return err == nil
}

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9 \-_.,:/]+$`)

func nameValidFunc(textToCheck string, _ rune) bool {
	return nameRegex.MatchString(textToCheck)
}
