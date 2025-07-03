package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"os"
	"path/filepath"
)

func list_dir() []string {
	pages := []string{}
	err := filepath.Walk("./journals", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			pages = append(pages, info.Name())
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return pages
}

func main() {
	fmt.Println("bunk")
	a := app.New()
	w := a.NewWindow("Encrypt UI")
	w.SetContent(widget.NewLabel("Test"))
	w.ShowAndRun()
}
