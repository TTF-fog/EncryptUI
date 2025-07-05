package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	d "github.com/sqweek/dialog"
	"golang.design/x/clipboard"
	"os"
	"regexp"
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
	settings.Recent = append([]string{"Decrypt Existing File"}, settings.Recent...)

	a := app.New()
	w := a.NewWindow("Encrypt UI")

	message := container.NewCenter(widget.NewLabel("Recent Files"))

	list := widget.NewList(
		func() int { return len(settings.Recent) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(func() string {
				if i != 0 {
					return getFileName(settings.Recent[i])
				}
				return settings.Recent[i]
			}())
		},
	)
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search History")

	var selectedFile string

	decryptButton := widget.NewButton("Decrypt", func() {
		if selectedFile == "" {
			return
		}
		if selectedFile == "Decrypt Existing File" {
			decryptExistingFile(w, list)
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
		}, false)
	})
	list.OnSelected = func(id widget.ListItemID) {
		if id == 0 {
			decryptButton.SetText("Decrypt Existing File")
		} else {
			decryptButton.SetText("Decrypt")
		}

		selectedFile = settings.Recent[id]
	}

	encryptButton := widget.NewButton("Encrypt", func() {
		nf_button := widget.NewButton("New File", func() {
			encryptNewFile(w, list)
		})
		ef_button := widget.NewButton("Encrypt Existing file", func() {
			encryptExistingFile(w, list)
		})
		st_button := container.NewHBox(widget.NewButton("Create Standalone Encrypted File", func() {
			showStandaloneOverview(w)
		}), widget.NewButton("?", func() {
			dialog.ShowInformation("Standalone Help", "A standalone encrypted file is one that needs no extra software to be opened on a different PC (Same OS) \n"+
				"Requires Go Installed ", w)
		}))
		formItems := []*widget.FormItem{
			widget.NewFormItem("", nf_button),
			widget.NewFormItem("", ef_button),
			widget.NewFormItem("", st_button),
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
	topContainer := container.NewVBox(decryptButton, container.NewBorder(searchEntry, nil, message, nil))
	content := container.NewBorder(topContainer, container.NewVBox(clearHistoryButton, encryptButton), nil, nil, list)

	origin := settings.Recent
	searchEntry.OnChanged = func(s string) {
		var results []string
		fmt.Println("Search Entry:", s)
		if s == "" {
			settings.Recent = origin
			list.Refresh()
			return
		}
		list.Refresh()
		for _, item := range origin {
			match, err := regexp.Match(".*^"+regexp.QuoteMeta(s)+".*", []byte(getFileName(item)))
			fmt.Println(item, match)
			if err != nil {
				dialog.NewError(err, w)
			}
			if match {
				results = append(results, item)
			}
		}
		settings.Recent = results
		list.Refresh()
	}
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
		encryptedText, err := EncryptAES([]byte(password), []byte(string(data)))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if err := os.WriteFile(file, []byte(encryptedText), 777); err != nil {
			dialog.ShowError(err, w)
			return
		}
		addRecentFile(file, list)
		dialog.ShowConfirm("Delete Unencrypted files?", "", func(b bool) {
			if b {
				err := os.Remove(file)
				if err != nil {
					dialog.ShowError(err, w)
				}
			}
		}, w)
	}, true)
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
			encryptedText, err := EncryptAES([]byte(password), []byte(entry.Text))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			saveEncryptedFile(w, list, encryptedText)
		}, true)
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

func getFileName(filePath string) string {
	var name string
	for i := len(filePath) - 1; i > 0; i-- {
		if filePath[i] == '/' {
			name = filePath[i+1:]
			break
		}
	}
	return name
}
func decryptExistingFile(w fyne.Window, list *widget.List) {
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
		decryptedText, err := DecryptAES([]byte(password), string(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		showDecryptedContent(w, decryptedText)

		addRecentFile(file, list)
	}, false)
}
