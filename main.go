package main

import (
	"encoding/json"
	"fmt"
	"github.com/drognisep/syspoll/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	flag "github.com/spf13/pflag"
	"os"
)

var (
	loadFile       string
	exportTemplate bool
	template       = []*System{
		{
			Name:          "Descriptive name",
			CheckInterval: "30s",
			Http: &CheckHttp{
				URL: "https://google.com",
			},
		},
	}
)

func main() {
	var systems []System

	flag.BoolVar(&exportTemplate, "template", false, "Export a template polling spec")
	flag.StringVar(&loadFile, "file", "", "Load a polling spec from file")
	flag.Parse()

	if exportTemplate {
		file, err := os.OpenFile("template.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to create/open file 'template.json': %v\n", err)
			os.Exit(1)
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(template); err != nil {
			fmt.Printf("Failed to write to template.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Wrote polling spec to template.json")
		os.Exit(0)
	}
	if len(loadFile) != 0 {
		file, err := os.Open(loadFile)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %v\n", loadFile, err)
			os.Exit(1)
		}
		if err := json.NewDecoder(file).Decode(&systems); err != nil {
			fmt.Printf("Failed to parse file '%s': %v\n", loadFile, err)
			os.Exit(1)
		}
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	table := DisplayTable(app, pages, systems...)
	if len(table.systems) == 0 {
		ShowSystemCreate(pages, func(newState *System, submitted bool) {
			if submitted && newState != nil {
				table.Add(*newState)
				pages.ShowPage(page.Systems)
				pages.SendToFront(page.Systems)
			}
		})
	}
	addBtn := tview.NewButton("Add System").SetSelectedFunc(func() {
		ShowSystemCreate(pages, func(newState *System, submitted bool) {
			if submitted && newState != nil {
				table.Add(*newState)
			}
		})
	})
	quitBtn := tview.NewButton("Quit").SetSelectedFunc(func() {
		app.Stop()
	})
	quitBtn.SetBackgroundColor(tcell.ColorRed)
	btnRow := Row(addBtn, saveButton(pages, table), quitBtn).SetItemPadding(1)
	grid := tview.NewGrid().
		SetColumns(0).
		SetRows(1, 0).
		AddItem(btnRow, 0, 0, 1, 1, 0, 0, true).
		AddItem(table, 1, 0, 1, 1, 0, 0, false)

	pages.AddPage(page.Systems, grid, true, len(table.systems) > 0)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func saveButton(pages *tview.Pages, table *systemTable) *tview.Button {
	const (
		savePage = "Save Page"
	)
	btn := tview.NewButton("Save")

	var filePath string
	saveForm := tview.NewForm().
		AddInputField("File path", "", 20, nil, func(text string) {
			filePath = text
		}).
		AddButton("Save", func() {
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				ShowErr(pages, err)
				return
			}
			defer file.Close()
			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			if err := enc.Encode(table.systems); err != nil {
				ShowErr(pages, err)
				return
			}
			pages.RemovePage(savePage)
		}).
		AddButton("Cancel", func() {
			pages.RemovePage(savePage)
		})

	btn.SetSelectedFunc(func() {
		pages.AddAndSwitchToPage(savePage, GridWrapper(saveForm, 50, 20), true)
	})
	return btn
}
