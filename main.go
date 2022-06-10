package main

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	flag "github.com/spf13/pflag"
	"os"
)

var (
	systems []System
)

var (
	loadFile string
	export   bool
)

func main() {
	flag.BoolVar(&export, "template", false, "Export a template polling spec")
	flag.StringVar(&loadFile, "file", "", "Load a polling spec from file")
	flag.Parse()

	if export {
		file, err := os.OpenFile("template.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to create/open file 'template.json': %v\n", err)
			os.Exit(1)
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(systems); err != nil {
			fmt.Printf("Failed to write to ")
		}
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
		os.Exit(0)
	}

	if len(systems) == 0 {
		fmt.Println("No systems loaded for polling")
		flag.Usage()
		os.Exit(0)
	}

	app := tview.NewApplication()

	pages := tview.NewPages()

	table := DisplayTable(app)
	pages.AddPage("Systems", table, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
