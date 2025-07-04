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
	"io"
	"os"
	"strings"
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
				decryptedText, err := DecryptAES([]byte(password), strings.TrimSpace(string(encryptedData)))
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
	encrypt_button := widget.NewButton("Encrypt", func() {
		dialog.ShowCustomConfirm("Type Of File", "New File", "Encrypt Existing File", widget.NewLabel("What would you like to encrypt"), func(b bool) {
			if !b {

				dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err != nil {
						panic(err)
					}
					if reader == nil {
						dialog.ShowError(errors.New("No File Chosen"), w)
					}
					data, err := io.ReadAll(reader)
					if err != nil {
						panic(err)
					}
					create_password_box(w, func(password string, entered bool) {
						if entered == true {
							encryptedText := EncryptAES([]byte(password), string(data))
							dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
								if err != nil {
									panic(err)
								}
								_, err = writer.Write([]byte(encryptedText))
								settings.Recent = append([]string{writer.URI().String()}, settings.Recent...)
								list.Refresh()
								if err != nil {
									panic(err)
								}
							}, w)

						}
					})

				}, w)

			}
		}, w)
	})
	content := container.NewBorder(topContainer, encrypt_button, nil, nil, list)
	list.OnSelected = func(id widget.ListItemID) {
		set_selected(&index, id)
	}
	w.SetContent(content)
	w.ShowAndRun()
}
