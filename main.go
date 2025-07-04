package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
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
	clipboard.Init()
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

		filePath := items[index]
		encryptedData, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err.(*os.PathError).Err, syscall.ENOENT) {
				dialog.ShowConfirm("Failed To Find File", "Remove File", func(b bool) {
					if b {
						for i, item := range items {
							if item == items[index] {
								items = append(items[:i], items[i+1:]...)
								settings.Recent = items
								set_settings(settings)
								list.Refresh()
							}
						}
					}
				}, w)
			} else {
				dialog.ShowError(err, w)
			}
			return
		}

		create_password_box(w, func(password string, entered bool) {
			if entered == true {
				fmt.Println("Entered:")
				decryptedText, err := DecryptAES([]byte(password), string(encryptedData))
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				text := widget.NewMultiLineEntry()
				copy_button := widget.NewButton("copy", func() {
					clipboard.Write(clipboard.FmtText, []byte(decryptedText))
				})
				text.SetText(decryptedText)
				text.Disable()
				items := []*widget.FormItem{
					widget.NewFormItem("Name", text),
					widget.NewFormItem("Copy", copy_button),
				}

				d := dialog.NewForm("Decrypted Content", "Done", "Cancel", items, func(b bool) { /*it forces me to have 2 buttons :/*/ }, w)
				d.Resize(fyne.NewSize(400, 200))
				d.Show()
			}
		})

	})
	topContainer := container.NewVBox(button, message)

	content := container.NewBorder(topContainer, nil, nil, nil, list)
	list.OnSelected = func(id widget.ListItemID) {
		set_selected(&index, id)
	}
	w.SetContent(content)
	w.ShowAndRun()
}
