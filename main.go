package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"syscall"
)

type Settings struct {
	Recent   []string          `json:"Recent"`
	Settings map[string]string `json:"Settings"`
}

var settings Settings

func main() {
	load_settings()
	var items = settings.Recent
	var index widget.ListItemID
	a := app.New()
	w := a.NewWindow("Encrypt UI")
	message := widget.NewLabel("Recent Files")

	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(items[i])
		})
	button := widget.NewButton("Decrypt", func() {
		fmt.Printf("%s", func() []byte {
			data, err := os.ReadFile(items[index])
			if err != nil {
				if errors.Is(err.(*os.PathError).Err, syscall.ENOENT) {
					dialog.ShowConfirm("Failed To Find File", "Remove File", func(b bool) {
						if b {
							for i, item := range items {
								if item == items[index] {
									items = append(items[:i], items[i+1:]...)
									fmt.Println(items)
									settings.Recent = items
									set_settings(settings)
									list.Refresh()
								}
							}
						}
					}, w)
				}
			}
			return data
		}())
	})
	topContainer := container.NewVBox(button, message)

	content := container.NewBorder(topContainer, nil, nil, nil, list)
	list.OnSelected = func(id widget.ListItemID) {
		set_selected(&index, id)
	}
	w.SetContent(content)
	w.ShowAndRun()
}
func load_settings() {
	data, f_err := os.ReadFile("./config.json")
	if f_err != nil {
		fmt.Println(f_err)
	}
	err := json.Unmarshal(data, &settings)
	if err != nil {
		panic(fmt.Errorf("Failed to load settings: %v", err))
	}
}
func set_settings(settings Settings) {
	marshal, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("./config.json", marshal, os.ModePerm)
}

func set_selected(index *widget.ListItemID, id int) {
	*index = id

}
