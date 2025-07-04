package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	d "github.com/sqweek/dialog"
	"golang.design/x/clipboard"
	"os"
	"strings"
)

type Settings struct {
	Recent   []string          `json:"Recent"`
	Settings map[string]string `json:"Settings"`
}

var settings Settings

func main() {
	loadSettings()
	clipboard.Init()

	a := app.New()
	w := a.NewWindow("Encrypt UI")

	message := widget.NewLabel("Recent Files")
	list := widget.NewList(
		func() int { return len(settings.Recent) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(settings.Recent[i]) },
	)

	var selectedFile string
	list.OnSelected = func(id widget.ListItemID) {
		selectedFile = settings.Recent[id]
	}

	decryptButton := widget.NewButton("Decrypt", func() {
		if selectedFile == "" {
			dialog.ShowError(errors.New("no file selected"), w)
			return
		}

		encryptedData, err := os.ReadFile(selectedFile)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				dialog.ShowConfirm("File Not Found", "Remove from recent?", func(b bool) {
					if b {
						removeRecentFile(selectedFile)
						list.Refresh()
					}
				}, w)
			} else {
				dialog.ShowError(err, w)
			}
			return
		}

		createPasswordBox(w, func(password string, ok bool) {
			if !ok {
				return
			}
			decryptedText, err := DecryptAES([]byte(password), strings.TrimSpace(string(encryptedData)))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			showDecryptedContent(w, decryptedText)
		})
	})

	encryptButton := widget.NewButton("Encrypt", func() {
		nf_button := widget.NewButton("New File", func() {
			encryptNewFile(w, list)
		})
		ef_button := widget.NewButton("Encrypt Existing file", func() {
			encryptExistingFile(w, list)
		})
		formItems := []*widget.FormItem{
			widget.NewFormItem("", nf_button),
			widget.NewFormItem("", ef_button),
		}
		dialog.ShowForm("Encryption Type", "", "", formItems, func(encryptExisting bool) { /*mandatory function, has no effect since confirmation/dismissal does nothing
			:/ */
		}, w)
	})
	clearHistoryButton := widget.NewButton("Clear History", func() {
		settings.Recent = []string{}
		setSettings(settings)
		list.Refresh()
	})
	topContainer := container.NewVBox(decryptButton, message)
	content := container.NewBorder(topContainer, container.NewVBox(clearHistoryButton, encryptButton), nil, nil, list)

	w.SetContent(content)
	w.ShowAndRun()
}

func removeRecentFile(filePath string) {
	for i, p := range settings.Recent {
		if p == filePath {

			settings.Recent = append(settings.Recent[:i], settings.Recent[i+1:]...)
			setSettings(settings)
			return
		}
	}
}

func encryptExistingFile(w fyne.Window, list *widget.List) {
	file, err := d.File().Load()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	data, err := os.ReadFile(file)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	createPasswordBox(w, func(password string, ok bool) {
		if !ok {
			return
		}
		encryptedText, err := EncryptAES([]byte(password), string(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if err := os.WriteFile(file, []byte(encryptedText), 777); err != nil {
			dialog.ShowError(err, w)
			return
		}
		addRecentFile(file, list)
	})
}

func encryptNewFile(w fyne.Window, list *widget.List) {
	entry := widget.NewMultiLineEntry()
	dialog.ShowForm("Enter Text", "Encrypt", "Cancel", []*widget.FormItem{widget.NewFormItem("Text", entry)}, func(ok bool) {
		if !ok {
			return
		}
		createPasswordBox(w, func(password string, ok bool) {
			if !ok {
				return
			}
			encryptedText, err := EncryptAES([]byte(password), entry.Text)
			println(encryptedText)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			saveEncryptedFile(w, list, encryptedText)
		})
	}, w)
}

func saveEncryptedFile(w fyne.Window, list *widget.List, encryptedText string) {
	file, err := d.File().Save()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	_, err = os.Create(file)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if err := os.WriteFile(file, []byte(encryptedText), 0644); err != nil {
		dialog.ShowError(err, w)
		return
	}
	addRecentFile(file, list)
}

func addRecentFile(filePath string, list *widget.List) {
	settings.Recent = append([]string{filePath}, settings.Recent...)
	setSettings(settings)
	list.Refresh()
}
