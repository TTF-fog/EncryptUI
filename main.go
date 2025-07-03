package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"path/filepath"
	"time"
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
	var data = list_dir()
	fmt.Println("bunk")
	a := app.New()
	w := a.NewWindow("Encrypt UI")
	message := widget.NewLabel("Recent Files")
	button := widget.NewButton("Decrypt", func() {
		formatted := time.Now().Format("Time: 03:04:05")
		message.SetText(formatted)
	})
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

	topContainer := container.NewVBox(button, message)
	content := container.NewBorder(topContainer, nil, nil, nil, list)

	w.SetContent(content)
	w.ShowAndRun()
}
